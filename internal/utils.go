package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"runtime/debug"

	"github.com/dvasilas/proteus/internal/config"
	"github.com/dvasilas/proteus/internal/proto"
	"github.com/dvasilas/proteus/internal/proto/qpu"
	"github.com/dvasilas/proteus/internal/proto/qpu_api"
	cli "github.com/dvasilas/proteus/internal/qpu/client"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

// QPU ...
type QPU struct {
	Client               cli.Client
	Conns                []*QPU
	Dataset              *qpu_api.DataSet
	QueryingCapabilities []*qpu.AttributePredicate
	Config               *config.Config
}

//---------------- API Functions -------------------

// ConnectToQPUGraph ...
func ConnectToQPUGraph(q *QPU) error {
	conns := make([]*QPU, len(q.Config.Connections))
	for i, conn := range q.Config.Connections {
		c, err := cli.NewClient(conn.Address)
		if err != nil {
			return err
		}
		retries := 0
		connConf, err := c.GetConfig()
		for err != nil || retries > 10 {
			ReportError(err)
			connConf, err = c.GetConfig()
			if err == nil {
				break
			}
			retries++
		}
		if err != nil {
			ReportError(err)
			return err
		}
		conns[i] = &QPU{
			Client:               c,
			QueryingCapabilities: connConf.GetSupportedQueries(),
			Dataset:              connConf.GetDataset(),
			Config: &config.Config{
				QpuType: connConf.QpuType,
				Port:    conn.Address,
			},
		}
	}
	q.Conns = conns
	calcQueryingCapabilities(q, conns)
	calcDataset(q, conns)
	return nil
}

//----------------- ObjectState --------------------

// ObjectState ...
type ObjectState struct {
	ObjectID   string
	ObjectType qpu.LogOperation_ObjectType
	Bucket     string
	State      qpu.ObjectState
	Timestamp  qpu.Vectorclock
}

// Marshal ...
func (o *ObjectState) Marshal() ([]byte, error) {
	marshalledState, err := marshalState(&o.State)
	if err != nil {
		return nil, err
	}
	marshalledVC, err := MarshalVectorClock(&o.Timestamp)
	if err != nil {
		return nil, err
	}
	marshalledParts := make([][]byte, 0)
	marshalledParts = append(marshalledParts, []byte(o.ObjectID))
	marshalledParts = append(marshalledParts, marshalObjectType(o.ObjectType))
	marshalledParts = append(marshalledParts, []byte(o.Bucket))
	marshalledParts = append(marshalledParts, marshalledState)
	marshalledParts = append(marshalledParts, marshalledVC)
	marshalledObjetState := bytes.Join(marshalledParts, []byte{'|'})
	return marshalledObjetState, nil
}

func marshalObjectType(t qpu.LogOperation_ObjectType) []byte {
	return []byte(qpu.LogOperation_ObjectType_name[int32(t)])
}

func marshalState(objectState *qpu.ObjectState) ([]byte, error) {
	return proto.Marshal(objectState)
}

func unmarshalObjectType(encodedType []byte) qpu.LogOperation_ObjectType {
	return qpu.LogOperation_ObjectType(qpu.LogOperation_ObjectType_value[string(encodedType)])
}

// MarshalVectorClock ...
func MarshalVectorClock(vc *qpu.Vectorclock) ([]byte, error) {
	return proto.Marshal(vc)
}

func unmarshalState(encodedObjectState []byte) (qpu.ObjectState, error) {
	var objectState qpu.ObjectState
	err := proto.Unmarshal(encodedObjectState, &objectState)
	return objectState, err
}

// UnmarshalVectorClock ...
func UnmarshalVectorClock(encodedVC []byte) (qpu.Vectorclock, error) {
	var vc qpu.Vectorclock
	err := proto.Unmarshal(encodedVC, &vc)
	return vc, err
}

// UnmarshalObject ...
func UnmarshalObject(data []byte) (ObjectState, error) {
	marshalledParts := bytes.Split(data, []byte{'|'})
	state, err := unmarshalState(marshalledParts[3])
	if err != nil {
		return ObjectState{}, err
	}
	vectorclock, err := UnmarshalVectorClock(marshalledParts[4])
	if err != nil {
		return ObjectState{}, err
	}
	objectState := ObjectState{
		ObjectID:   string(marshalledParts[0]),
		ObjectType: unmarshalObjectType(marshalledParts[1]),
		Bucket:     string(marshalledParts[2]),
		State:      state,
		Timestamp:  vectorclock,
	}
	return objectState, nil
}

// GetMessageSize ...
func GetMessageSize(streamRec *qpu_api.ResponseStreamRecord) (int, error) {
	buff, err := proto.Marshal(streamRec)
	if err != nil {
		return -1, err
	}
	bytesBuff := bytes.NewBuffer(buff)
	return bytesBuff.Len(), nil
}

// SubQuery ...
type SubQuery struct {
	SubQuery []*qpu.AttributePredicate
	Endpoint *QPU
}

// ObjectStateJSON ...
type ObjectStateJSON struct {
	ObjectID   string
	ObjectType string
	Bucket     string
	State      []struct {
		AttrKey   string
		AttrType  string
		AttrValue string
	}
	Timestamp map[string]uint64
}

// ValueToString converts an attribute value to a string
func ValueToString(val *qpu.Value) string {
	switch val.Val.(type) {
	case *qpu.Value_Int:
		return strconv.Itoa(int(val.GetInt()))
	case *qpu.Value_Flt:
		return fmt.Sprintf("%f", val.GetFlt())
	case *qpu.Value_Str:
		return val.GetStr()
	default:
		return ""
	}
}

// CanRespondToQuery ...
func CanRespondToQuery(predicate []*qpu.AttributePredicate, capabilities []*qpu.AttributePredicate) (bool, error) {
	if len(capabilities) == 0 {
		return true, nil
	}
	for _, p := range predicate {
		matchesCapabilities := false
		for _, c := range capabilities {
			if c.GetAttr().GetAttrKey() == p.GetAttr().GetAttrKey() {
				lb, err := Compare(p.GetLbound(), c.GetLbound())
				if err != nil {
					return false, err
				}
				ub, err := Compare(p.GetUbound(), c.GetUbound())
				if err != nil {
					return false, err
				}
				if lb >= 0 && ub <= 0 {
					matchesCapabilities = true
					break
				}
			}
		}
		if !matchesCapabilities {
			return false, nil
		}
	}
	return true, nil
}

// Compare ...
func Compare(a, b *qpu.Value) (int, error) {
	if valueType(a) != valueType(b) {
		return 0, errors.New("cannot compare different types of Value")
	}
	const TOLERANCE = 0.000001
	switch a.GetVal().(type) {
	case *qpu.Value_Flt:
		diff := a.GetFlt() - b.GetFlt()
		if diff := math.Abs(diff); diff < TOLERANCE {
			return 0, nil
		}
		if diff < 0 {
			return -1, nil
		}
		return 1, nil
	case *qpu.Value_Int:
		return int(a.GetInt() - b.GetInt()), nil
	case *qpu.Value_Str:
		return strings.Compare(a.GetStr(), b.GetStr()), nil
	}
	return 0, errors.New("unknown Value type")
}

// AttrMatchesPredicate checks if an object attribute matches a given predicate.
func AttrMatchesPredicate(predicate *qpu.AttributePredicate, attr *qpu.Attribute) (bool, error) {
	if keyMatch(predicate.GetAttr().GetAttrKey(), attr) {
		return rangeMatch(predicate, attr)
	}
	return false, nil
}

func keyMatch(objectName string, attr *qpu.Attribute) bool {
	if objectName == attr.GetAttrKey() {
		return true
	}
	return false
}

// within the range [greaterOrEqual, lessThan)
func rangeMatch(pred *qpu.AttributePredicate, attr *qpu.Attribute) (bool, error) {
	lb, err := Compare(attr.GetValue(), pred.GetLbound())
	if err != nil {
		return false, err
	}
	ub, err := Compare(attr.GetValue(), pred.GetUbound())
	if err != nil {
		return false, err
	}
	if lb >= 0 && ub < 0 {
		return true, nil
	}
	return false, nil
}

// Ping ...
func Ping(stream qpu_api.QPU_QueryServer, msg *qpu_api.PingMsg) error {
	seqID := msg.GetSeqId()
	for {
		seqID++
		if err := stream.Send(protoutils.ResponseStreamRecord(seqID, qpu_api.ResponseStreamRecord_HEARTBEAT, nil)); err != nil {
			return err
		}
		p, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		seqID = p.GetPing().GetSeqId()
		fmt.Println(seqID)
	}
	return nil
}

//----------- Stream Consumer Functions ------------

//QueryResponseConsumer receives a QueryResponseStream, iteratively reads from the stream, and processes each input element based on a given function
func QueryResponseConsumer(pred []*qpu.AttributePredicate, streamIn qpu_api.QPU_QueryClient, streamOut qpu_api.QPU_QueryServer, process func([]*qpu.AttributePredicate, *qpu_api.ResponseStreamRecord, qpu_api.QPU_QueryServer, *int64) error, errChan chan error) {
	seqID := int64(0)
	go func() {
		for {
			streamRec, err := streamIn.Recv()
			if err == io.EOF {
				if streamOut != nil {
					errChan <- streamOut.Send(
						protoutils.ResponseStreamRecord(
							seqID,
							qpu_api.ResponseStreamRecord_END_OF_STREAM,
							&qpu.LogOperation{},
						),
					)
				}
				errChan <- err
				return
			} else if err != nil {
				errChan <- err
				return
			}
			if err = process(pred, streamRec, streamOut, &seqID); err != nil {
				errChan <- err
				return
			}
		}
	}()
}

//---------------- Internal Functions --------------

func mergeDatasets(a, b *qpu_api.DataSet) {
	for databaseID := range b.GetDatabases() {
		if db, ok := a.GetDatabases()[databaseID]; ok {
			for datacenterID := range b.GetDatabases()[databaseID].GetDatacenters() {
				if dc, ok := db.GetDatacenters()[datacenterID]; ok {
					dc.Shards = append(dc.GetShards(), b.GetDatabases()[databaseID].GetDatacenters()[datacenterID].GetShards()...)
				} else {
					db.GetDatacenters()[datacenterID] = b.GetDatabases()[databaseID].GetDatacenters()[datacenterID]
				}
			}
		} else {
			a.GetDatabases()[databaseID] = b.GetDatabases()[databaseID]
		}
	}
}

func calcQueryingCapabilities(q *QPU, conns []*QPU) {
	for _, c := range conns {
		q.QueryingCapabilities = append(q.QueryingCapabilities, c.QueryingCapabilities...)
	}
}

func calcDataset(q *QPU, conns []*QPU) {
	q.Dataset = conns[0].Dataset
	for _, c := range conns[:1] {
		mergeDatasets(q.Dataset, c.Dataset)
	}
}

func valueType(v *qpu.Value) int {
	switch v.GetVal().(type) {
	case *qpu.Value_Flt:
		return 0
	case *qpu.Value_Int:
		return 1
	case *qpu.Value_Str:
		return 2
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return rand.Int()
}

//----------------- responseTime -------------------

// ResponseTime ...
type ResponseTime []time.Duration

func (t ResponseTime) Len() int {
	return len(t)
}
func (t ResponseTime) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}
func (t ResponseTime) Less(i, j int) bool {
	return t[i].Nanoseconds() < t[j].Nanoseconds()
}

//---------------- error handling ------------------

// ReportError prints the given error and the stack trace returned by runtime.Stack.
func ReportError(e error) {
	log.WithFields(log.Fields{"error": e}).Warn("error")
	debug.PrintStack()
}

// Warn prints the given error
func Warn(e error) {
	log.WithFields(log.Fields{"error": e}).Warn("warning")
	debug.PrintStack()
}

//----------- query metadata parameters ------------

// MaxResponseCount ..
func MaxResponseCount(metadata map[string]string) (int64, error) {
	maxResponseCount := int64(-1)
	if metadata != nil {
		if val, ok := metadata["maxResponseCount"]; ok {
			mdMaxResponseCountVal, err := strconv.ParseInt(val, 10, 0)
			if err != nil {
				return maxResponseCount, err
			}
			maxResponseCount = mdMaxResponseCountVal
		}
	}
	return maxResponseCount, nil
}