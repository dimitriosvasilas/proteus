package client

import (
	"context"

	pb "github.com/dimitriosvasilas/modqp/protos/datastore"
	"google.golang.org/grpc"
)

type activeStreams struct {
	opSubStreams map[int64]context.CancelFunc
}

//Client ...
type Client struct {
	dsClient      pb.DataStoreClient
	activeStreams activeStreams
}

//SubscribeStates ...
func (c *Client) SubscribeStates(ts int64) (pb.DataStore_SubscribeStatesClient, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.dsClient.SubscribeStates(ctx, &pb.SubRequest{Timestamp: ts})
	return stream, cancel, err
}

// SubscribeOps ...
func (c *Client) SubscribeOps(ts int64) (pb.DataStore_SubscribeOpsClient, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	stream, err := c.dsClient.SubscribeOps(ctx, &pb.SubRequest{Timestamp: ts})
	c.activeStreams.opSubStreams[ts] = cancel
	return stream, cancel, err
}

//StopOpsSubscription ...
func (c *Client) StopOpsSubscription(subID int64) {
	cancel := c.activeStreams.opSubStreams[subID]
	cancel()
	return
}

//GetSnapshot ...
func (c *Client) GetSnapshot(ts int64) (pb.DataStore_GetSnapshotClient, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())

	stream, err := c.dsClient.GetSnapshot(ctx, &pb.SubRequest{Timestamp: ts})
	return stream, cancel, err
}

//NewClient ...
func NewClient(address string) (Client, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return Client{}, nil, err
	}
	activeStrMap := make(map[int64]context.CancelFunc)
	return Client{pb.NewDataStoreClient(conn), activeStreams{activeStrMap}}, conn, nil
}
