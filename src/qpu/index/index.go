package index

import (
	"context"
	"errors"
	"io"
	"math/rand"
	"time"

	"github.com/dvasilas/proteus/src"
	"github.com/dvasilas/proteus/src/config"
	"github.com/dvasilas/proteus/src/protos"
	pbQPU "github.com/dvasilas/proteus/src/protos/qpu"
	pbUtils "github.com/dvasilas/proteus/src/protos/utils"
	"github.com/dvasilas/proteus/src/qpu/filter"
	"github.com/dvasilas/proteus/src/qpu/index/antidote"
	"github.com/dvasilas/proteus/src/qpu/index/inMem"
	log "github.com/sirupsen/logrus"
)

// IQPU implements an index QPU
type IQPU struct {
	qpu         *utils.QPU
	index       indexStore
	cancelFuncs []context.CancelFunc
	forwardM    map[int]chan *pbQPU.ResponseStreamRecord
}

// indexStore describes the interface that any index implementation needs to expose
// to work with this module.
type indexStore interface {
	Put(*pbUtils.Attribute, utils.ObjectState) error
	Get(*pbUtils.AttributePredicate) (map[string]utils.ObjectState, error)
	RemoveOldEntry(*pbUtils.Attribute, string) error
	Print()
}

//---------------- API Functions -------------------

// QPU creates an index QPU
func QPU(conf *config.Config) (*IQPU, error) {
	rand.Seed(time.Now().UnixNano())
	q := &IQPU{
		qpu: &utils.QPU{
			Config:               conf,
			QueryingCapabilities: conf.IndexConfig.IndexingConfig,
		},
		forwardM: make(map[int]chan *pbQPU.ResponseStreamRecord),
	}

	if err := utils.ConnectToQPUGraph(q.qpu); err != nil {
		return nil, err
	}
	if len(q.qpu.Conns) > 1 {
		return nil, errors.New("index QPUs support a single connection")
	}
	var index indexStore
	var err error
	switch q.qpu.Config.IndexConfig.IndexStore.Store {
	case config.INMEM:
		index, err = btreeindex.New(
			q.qpu.Config.IndexConfig.IndexingConfig[0].GetAttr().GetAttrKey(),
			q.qpu.Config.IndexConfig.IndexingConfig[0].GetAttr().GetAttrType())
		if err != nil {
			return &IQPU{}, err
		}
	case config.ANT:
		index, err = mapindex.New(conf)
		if err != nil {
			return &IQPU{}, err
		}
	}
	q.index = index

	var sync bool
	switch conf.IndexConfig.ConsLevel {
	case "sync":
		sync = true
	case "async":
		sync = false
	default:
		return nil, errors.New("unknown index consistency level")
	}
	pred := []*pbUtils.AttributePredicate{}
	q.cancelFuncs = make([]context.CancelFunc, len(q.qpu.Conns))
	for i, conn := range q.qpu.Conns {
		streamIn, cancel, err := conn.Client.Query(
			pred,
			protoutils.SnapshotTimePredicate(
				protoutils.SnapshotTime(pbUtils.SnapshotTime_INF, nil),
				protoutils.SnapshotTime(pbUtils.SnapshotTime_INF, nil),
			),
			true,
			sync,
		)
		if err != nil {
			cancel()
			return nil, err
		}
		q.cancelFuncs[i] = cancel
		go q.opConsumer(streamIn, cancel, sync)
	}

	if err := q.catchUp(); err != nil {
		return nil, err
	}
	return q, nil

}

// Query implements the Query API for the index QPU
func (q *IQPU) Query(streamOut pbQPU.QPU_QueryServer) error {
	request, err := streamOut.Recv()
	if err == io.EOF {
		return errors.New("Query received EOF")
	}
	if err != nil {
		return err
	}
	req := request.GetRequest()
	log.WithFields(log.Fields{"req": req}).Debug("Query request")

	if req.GetClock().GetUbound().GetType() < req.GetClock().GetUbound().GetType() {
		return errors.New("lower bound of timestamp attribute cannot be greater than the upper bound")
	}
	if req.GetClock().GetLbound().GetType() != pbUtils.SnapshotTime_LATEST &&
		req.GetClock().GetLbound().GetType() != pbUtils.SnapshotTime_INF &&
		req.GetClock().GetUbound().GetType() != pbUtils.SnapshotTime_LATEST &&
		req.GetClock().GetUbound().GetType() != pbUtils.SnapshotTime_INF {
		return errors.New("not supported")
	}
	if req.GetClock().GetLbound().GetType() == pbUtils.SnapshotTime_LATEST || req.GetClock().GetUbound().GetType() == pbUtils.SnapshotTime_LATEST {
		indexResult, err := q.lookup(req.GetPredicate()[0])
		if err != nil {
			return err
		}
		var seqID int64
		for _, item := range indexResult {
			logOp := protoutils.LogOperation(
				item.ObjectID,
				item.Bucket,
				item.ObjectType,
				&item.Timestamp,
				protoutils.PayloadState(&item.State),
			)
			if err := streamOut.Send(protoutils.ResponseStreamRecord(
				seqID,
				pbQPU.ResponseStreamRecord_STATE,
				logOp,
			)); err != nil {
				return err
			}
			seqID++
		}
	}
	if req.GetClock().GetLbound().GetType() == pbUtils.SnapshotTime_INF || req.GetClock().GetUbound().GetType() == pbUtils.SnapshotTime_INF {
		forwardCh := make(chan *pbQPU.ResponseStreamRecord, 0)
		errCh := make(chan error)
		chID := rand.Int()
		q.forwardM[chID] = forwardCh
		go q.forward(req.GetPredicate(), forwardCh, streamOut, errCh, chID)
		err := <-errCh
		if err != nil {
			return err
		}
	}
	return nil
}

// GetConfig implements the GetConfig API for the index QPU
func (q *IQPU) GetConfig() (*pbQPU.ConfigResponse, error) {
	resp := protoutils.ConfigRespοnse(
		q.qpu.Config.QpuType,
		q.qpu.QueryingCapabilities,
		q.qpu.Dataset)
	return resp, nil
}

// Cleanup ...
func (q *IQPU) Cleanup() {
	log.Info("index QPU cleanup")
	for _, cFunc := range q.cancelFuncs {
		cFunc()
	}
}

//----------- Stream Consumer Functions ------------

func (q *IQPU) forward(predicate []*pbUtils.AttributePredicate, forwardCh chan *pbQPU.ResponseStreamRecord, streamOut pbQPU.QPU_QueryServer, errCh chan error, chanID int) {
	for streamRec := range forwardCh {
		match, err := filter.Filter(predicate, streamRec)
		if err != nil {
			errCh <- err
			break
		}
		if match {
			if err := streamOut.Send(streamRec); err != nil {
				errCh <- err
				break
			}
		}
	}
	errCh <- nil
	close(forwardCh)
	delete(q.forwardM, chanID)
}

// Receives an stream of update operations
// Updates the index for each operation
// TODO: Query a way to handle an error here
func (q *IQPU) opConsumer(stream pbQPU.QPU_QueryClient, cancel context.CancelFunc, sync bool) {
	for {
		streamRec, err := stream.Recv()
		if err == io.EOF {
			// TODO: see datastoredriver to fix this
			log.Fatal("indexQPU:opConsumer received EOF, which is not expected")
			return
		} else if err != nil {
			log.Fatal("opConsumer err", err)
			return
		} else {
			if streamRec.GetType() == pbQPU.ResponseStreamRecord_UPDATEDELTA {
				log.WithFields(log.Fields{
					"operation": streamRec,
				}).Debug("index QPU received operation")

				if err := q.updateIndex(streamRec.GetLogOp()); err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"op":    streamRec,
					}).Fatal("opConsumer: index Update failed")
					return
				}
				if sync {
					log.Debug("QPUServer:index updated, sending ACK")
					if err := stream.Send(protoutils.RequestStreamAck(streamRec.GetSequenceId())); err != nil {
						log.Fatal("opConsumer stream.Send failed")
						return
					}
				}
				for _, ch := range q.forwardM {
					ch <- streamRec
				}
			}
		}
	}
}

//---------------- Internal Functions --------------

func (q *IQPU) update(attrStateO, attrStateN *pbUtils.Attribute, objectID string, object utils.ObjectState) error {
	if attrStateO != nil {
		q.index.RemoveOldEntry(attrStateO, objectID)
	}
	if err := q.index.Put(attrStateN, object); err != nil {
		return err
	}
	return nil
}

func (q *IQPU) lookup(predicate *pbUtils.AttributePredicate) (map[string]utils.ObjectState, error) {
	return q.index.Get(predicate)
}

// Given an operation sent from the data store, updates the index
func (q *IQPU) updateIndex(op *pbUtils.LogOperation) error {
	log.WithFields(log.Fields{
		"operation":       op,
		"querying config": q.qpu.QueryingCapabilities,
	}).Debug("updateIndex")

	state := utils.ObjectState{
		ObjectID:   op.GetObjectId(),
		ObjectType: op.GetObjectType(),
		Bucket:     op.GetBucket(),
		State:      *op.GetPayload().GetDelta().GetNew(),
	}
	for _, attr := range state.State.GetAttrs() {
		toIndex := true
		for _, pred := range q.qpu.QueryingCapabilities {
			match, err := filter.AttrMatchesPredicate(pred, attr)
			if err != nil {
				return err
			}
			if !match {
				toIndex = false
				break
			}
		}
		if toIndex {
			for _, attrOld := range op.GetPayload().GetDelta().GetOld().GetAttrs() {
				if attr.GetAttrKey() == attrOld.GetAttrKey() && attr.GetAttrType() == attrOld.GetAttrType() {
					return q.update(attrOld, attr, state.ObjectID, state)
				}
			}
			return q.update(nil, attr, state.ObjectID, state)
		}
	}
	return nil
}

// catchUp performs an index catch-up operation.
// It reads the latest snapshot for the underlying data store, and builds an index.
func (q *IQPU) catchUp() error {
	errChan := make(chan error)
	for _, conn := range q.qpu.Conns {
		pred := make([]*pbUtils.AttributePredicate, 0)
		stream, cancel, err := conn.Client.Query(
			pred,
			protoutils.SnapshotTimePredicate(
				protoutils.SnapshotTime(pbUtils.SnapshotTime_LATEST, nil),
				protoutils.SnapshotTime(pbUtils.SnapshotTime_LATEST, nil),
			),
			false,
			false,
		)
		defer cancel()
		if err != nil {
			return err
		}
		utils.QueryResponseConsumer(pred, stream, nil, q.updateIndexCatchUp, errChan)
	}
	streamCnt := len(q.qpu.Conns)
	for streamCnt > 0 {
		select {
		case err := <-errChan:
			if err == io.EOF {
				streamCnt--
			} else if err != nil {
				return err
			}
		}
	}
	return nil
}

func (q *IQPU) updateIndexCatchUp(pred []*pbUtils.AttributePredicate, streamRec *pbQPU.ResponseStreamRecord, streamOut pbQPU.QPU_QueryServer) error {
	delta := protoutils.PayloadDelta(nil, streamRec.GetLogOp().GetPayload().GetState())
	op := protoutils.LogOperation(
		streamRec.GetLogOp().GetObjectId(),
		streamRec.GetLogOp().GetBucket(),
		streamRec.GetLogOp().GetObjectType(),
		streamRec.GetLogOp().GetTimestamp(),
		delta,
	)
	return q.updateIndex(op)
}