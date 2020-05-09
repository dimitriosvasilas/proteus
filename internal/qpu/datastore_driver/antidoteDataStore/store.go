package store

import (
	"context"
	"errors"
	"io"
	"strconv"
	"strings"

	antidote "github.com/AntidoteDB/antidote-go-client"
	"github.com/dvasilas/proteus/internal/proto"
	pbAntidote "github.com/dvasilas/proteus/internal/proto/antidote"
	"github.com/dvasilas/proteus/internal/proto/qpu"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

//AntidoteDataStore ...
type AntidoteDataStore struct {
	logStreamEndpoint string
	endpoint          string
	antidoteCli       *antidote.Client
}

//---------------- API Functions -------------------

//New creates and initializes an instance of AntidoteDataStore
func New(endp, logPropEndp string) AntidoteDataStore {
	s := AntidoteDataStore{
		endpoint:          endp,
		logStreamEndpoint: logPropEndp,
	}
	return s
}

//SubscribeOps subscribes to updates from AntidoteDB (supports only async mode)
// each time an update is received, it is formated as a qpu.LogOperation
// and sent to datastoredriver via a channel
func (ds AntidoteDataStore) SubscribeOps(msg chan *qpu.LogOperation, ack chan bool, sync bool) (*grpc.ClientConn, <-chan error) {
	errs := make(chan error, 1)

	if sync {
		errs <- errors.New("SubscribeOps sync not supported for AntidoteDataStore")
		return nil, errs
	}

	conn, err := grpc.Dial(ds.logStreamEndpoint, grpc.WithInsecure())
	if err != nil {
		errs <- err
		return nil, errs
	}
	client := pbAntidote.NewServiceClient(conn)
	ctx := context.Background()
	stream, err := client.WatchAsync(ctx, &pbAntidote.SubRequest{Timestamp: 0})
	if err != nil {
		errs <- err
		return conn, errs
	}
	ds.opConsumer(stream, msg, errs)
	return conn, errs
}

//GetSnapshot reads a snapshot of all objects stored in an Antidotedb bucket,
// not yet implemented
func (ds AntidoteDataStore) GetSnapshot(bucket string, msg chan *qpu.LogOperation) <-chan error {
	return ds.readSnapshot(bucket, msg, ds.processAndForwardObject)
}

// Op ...
func (ds AntidoteDataStore) Op(op *qpu.LogOperation) {
}

//----------- Stream Consumer Functions ------------

//opConsumer creates a goroutine that receives a stream of updates from AntidoteDB,
// each time an update is received, it is parsed to a qpu.LogOperation object
// which is then sent to the datastoredriver via a channel
func (ds AntidoteDataStore) opConsumer(stream pbAntidote.Service_WatchAsyncClient, msg chan *qpu.LogOperation, errs chan error) {
	go func() {
		for {
			op, err := stream.Recv()
			log.WithFields(log.Fields{"op": op}).Debug("antidoteDataStore:opConsumer received op")
			if err == io.EOF {
				errs <- errors.New("antidoteDataStore:opConsumer received EOF")
				break
			} else if err != nil {
				errs <- err
				break
			} else {
			}
			msg <- ds.formatOperation(op)
		}
		close(msg)
		close(errs)
	}()
}

//---------------- Internal Functions --------------

//readSnapshot retrieves all objects stored in the given bucket
// and for each object calls the processObj function
func (ds AntidoteDataStore) readSnapshot(bucketName string, msg chan *qpu.LogOperation, processObj func(string, string, *antidote.MapReadResult, []antidote.MapEntryKey, chan *qpu.LogOperation, chan error)) <-chan error {
	errs := make(chan error, 1)
	endpoint := strings.Split(ds.endpoint, ":")
	port, err := strconv.ParseInt(endpoint[1], 10, 64)
	if err != nil {
		errs <- err
		return errs
	}
	c, err := antidote.NewClient(antidote.Host{Name: endpoint[0], Port: int(port)})
	if err != nil {
		errs <- err
		return errs
	}
	ds.antidoteCli = c
	bucket := antidote.Bucket{Bucket: []byte(bucketName)}
	tx := ds.antidoteCli.CreateStaticTransaction()
	pIndex, err := bucket.ReadSet(tx, antidote.Key([]byte("."+bucketName)))
	if err != nil {
		errs <- err
		return errs
	}
	go func() {
		for _, obj := range pIndex {
			objVal, err := bucket.ReadMap(tx, antidote.Key([]byte(obj)))
			if err != nil {
				errs <- err
			}
			entries := objVal.ListMapKeys()
			if err != nil {
				errs <- err
			}
			processObj(string(obj), bucketName, objVal, entries, msg, errs)
		}
		close(msg)
	}()
	return errs
}

//processAndForwardObject reads the content of an object (map crdt)
// creates a *qpu.LogOperation and sends it to the datastoredriver to datastoredriver via a channel
func (ds AntidoteDataStore) processAndForwardObject(bucket string, key string, val *antidote.MapReadResult, entries []antidote.MapEntryKey, msg chan *qpu.LogOperation, errs chan error) {
	attrs := make([]*qpu.Attribute, len(entries))
	for i, e := range entries {
		switch e.CrdtType {
		case antidote.CRDTType_LWWREG:
			r, err := val.Reg(e.Key)
			if err != nil {
				errs <- err
				close(errs)
			}
			attrs[i] = protoutils.Attribute(string(e.Key), protoutils.ValueStr(string(r)))
		case antidote.CRDTType_COUNTER:
			c, err := val.Counter(e.Key)
			if err != nil {
				errs <- err
				close(errs)
			}
			attrs[i] = protoutils.Attribute(string(e.Key), protoutils.ValueInt(int64(c)))
		case antidote.CRDTType_ORSET:
		case antidote.CRDTType_RRMAP:
		case antidote.CRDTType_MVREG:
		}
	}
	state := protoutils.ObjectState(attrs)
	obj := protoutils.LogOperation(
		string(key),
		bucket,
		qpu.LogOperation_MAPCRDT,
		protoutils.Vectorclock(map[string]uint64{"antidote": uint64(0)}),
		protoutils.PayloadState(state),
	)
	log.WithFields(log.Fields{
		"object": obj,
	}).Debug("antidoteDataStore: snapshot")
	msg <- obj
}

func (ds AntidoteDataStore) formatOperation(logOp *pbAntidote.LogOperation) *qpu.LogOperation {
	var payload *qpu.Payload
	switch logOp.GetPayload().GetVal().(type) {
	case *pbAntidote.LogOperation_Payload_Op:
		ops := logOp.GetPayload().GetOp().GetOp()
		attrs := make([]*qpu.Attribute, len(ops))
		updates := make([]*qpu.Operation_Update, len(ops))
		for i := range ops {
			// var typ qpu.Attribute_AttributeType
			// if ops[i].GetObject().GetType() == "antidote_crdt_counter_pn" {
			// 	typ = qpu.Attribute_CRDTCOUNTER
			// } else if ops[i].GetObject().GetType() == "antidote_crdt_register_lww" {
			// 	typ = qpu.Attribute_CRDTLWWREG
			// }
			attrs[i] = protoutils.Attribute(ops[i].GetObject().GetKey(), nil)
			updates[i] = protoutils.Update(ops[i].GetUpdate().GetOpType(), crdtValToValue(ops[i].GetUpdate().GetValue()))
		}
		payload = protoutils.PayloadOp(attrs, updates)
	case *pbAntidote.LogOperation_Payload_Delta:
		oldState := mapCrdtStateToObjectState(logOp.GetPayload().GetDelta().GetOld().GetState())
		newState := mapCrdtStateToObjectState(logOp.GetPayload().GetDelta().GetNew().GetState())
		payload = protoutils.PayloadDelta(oldState, newState)
	}
	op := protoutils.LogOperation(
		logOp.GetKey(),
		logOp.GetBucket(),
		qpu.LogOperation_MAPCRDT,
		//TODO antidote should return vector clock
		protoutils.Vectorclock(map[string]uint64{"antidote": uint64(logOp.GetCommitTime())}),
		payload,
	)
	return op
}

func mapCrdtStateToObjectState(crdtSt []*pbAntidote.CrdtMapState_MapState) *qpu.ObjectState {
	attrs := make([]*qpu.Attribute, len(crdtSt))
	for i := range crdtSt {
		// var typ qpu.Attribute_AttributeType
		// if crdtSt[i].GetObject().GetType() == "antidote_crdt_counter_pn" {
		// 	typ = qpu.Attribute_CRDTCOUNTER
		// } else if crdtSt[i].GetObject().GetType() == "antidote_crdt_register_lww" {
		// 	typ = qpu.Attribute_CRDTLWWREG
		// }
		attrs[i] = protoutils.Attribute(crdtSt[i].GetObject().GetKey(), crdtValToValue(crdtSt[i].GetValue()))
	}
	return protoutils.ObjectState(attrs)
}

func crdtValToValue(val *pbAntidote.CrdtValue) *qpu.Value {
	value := &qpu.Value{}
	switch val.GetVal().(type) {
	case *pbAntidote.CrdtValue_Str:
		value.Val = &qpu.Value_Str{
			Str: val.GetStr(),
		}
	case *pbAntidote.CrdtValue_Int:
		value.Val = &qpu.Value_Int{
			Int: val.GetInt(),
		}
	}
	return value
}