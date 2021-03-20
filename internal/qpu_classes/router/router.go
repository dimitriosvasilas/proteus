package routerqpu

import (
	"errors"
	"time"

	"github.com/dvasilas/proteus/internal/libqpu"
	"github.com/dvasilas/proteus/internal/libqpu/utils"
	"github.com/dvasilas/proteus/internal/proto/qpuextapi"
	qpugraph "github.com/dvasilas/proteus/internal/qpuGraph"
	responsestream "github.com/dvasilas/proteus/internal/responseStream"
	"github.com/opentracing/opentracing-go"

	"github.com/dvasilas/proteus/internal/proto/qpuapi"
)

// RouterQPU ...
type RouterQPU struct {
	adjacentQPUs []*libqpu.AdjacentQPU
	conf         *libqpu.QPUConfig
}

// ---------------- API Functions -------------------

// InitClass ...
func InitClass(qpu *libqpu.QPU, catchUpDoneCh chan int) (*RouterQPU, error) {

	rqpu := &RouterQPU{
		adjacentQPUs: qpu.AdjacentQPUs,
		conf:         qpu.Config,
	}

	go func() {
		time.Sleep(2)
		catchUpDoneCh <- 0
	}()

	return rqpu, nil
}

// ProcessQuerySnapshot ...
func (q *RouterQPU) ProcessQuerySnapshot(query libqpu.ASTQuery, md map[string]string, sync bool, parentSpan opentracing.Span) (<-chan libqpu.LogOperation, <-chan error) {
	return nil, nil
}

// ClientQuery ...
func (q *RouterQPU) ClientQuery(query libqpu.ASTQuery, parentSpan opentracing.Span) (*qpuextapi.QueryResp, error) {
	var forwardTo *libqpu.AdjacentQPU
	found := false
	for _, adjQPU := range q.adjacentQPUs {
		for _, table := range adjQPU.OutputSchema {
			if table == query.GetTable() {
				forwardTo = adjQPU
				found = true
				break
			}
		}
	}

	if !found {
		return nil, utils.Error(errors.New("unknown table"))
	}

	subQueryResponseStream, err := qpugraph.SendQuery(libqpu.NewQuery(nil, query.Q), forwardTo)
	if err != nil {
		return nil, err
	}

	respCh := make(chan libqpu.ResponseRecord)
	go func() {
		if err = responsestream.StreamConsumer(subQueryResponseStream, q.conf.ProcessingConfig.Input.MaxWorkers, q.conf.ProcessingConfig.Input.MaxJobQueue, q.processRespRecord, nil, respCh); err != nil {
			panic(err)
		}
	}()

	respRecords := make([]*qpuextapi.QueryRespRecord, 0)

	for record := range respCh {
		attributes := make(map[string]string)
		for k, v := range record.GetAttributes() {
			valStr, err := utils.ValueToStr(v)
			if err != nil {
				return nil, err
			}
			attributes[k] = valStr
		}

		respRecords = append(respRecords, &qpuextapi.QueryRespRecord{
			RecordId:   record.GetRecordID(),
			Attributes: attributes,
			Timestamp:  record.GetLogOp().GetTimestamp().GetVc(),
		})
	}

	return &qpuextapi.QueryResp{
		RespRecord: respRecords,
	}, nil
}

// ProcessQuerySubscribe ...
func (q *RouterQPU) ProcessQuerySubscribe(query libqpu.ASTQuery, md map[string]string, sync bool) (int, <-chan libqpu.LogOperation, <-chan error) {
	return -1, nil, nil
}

// RemovePersistentQuery ...
func (q *RouterQPU) RemovePersistentQuery(table string, queryID int) {
}

// GetMetrics ...
func (q *RouterQPU) GetMetrics(*qpuextapi.MetricsRequest) (*qpuextapi.MetricsResponse, error) {
	return nil, nil
}

// ---------------- Internal Functions --------------

func (q *RouterQPU) processRespRecord(respRecord libqpu.ResponseRecord, data interface{}, recordCh chan libqpu.ResponseRecord) error {
	respRecordType, err := respRecord.GetType()
	if err != nil {
		return err
	}

	if respRecordType == libqpu.EndOfStream {
		close(recordCh)
	} else {
		recordCh <- respRecord
	}

	return nil
}

// GetConfig ...
func (q RouterQPU) GetConfig() *qpuapi.ConfigResponse {
	return &qpuapi.ConfigResponse{}
}
