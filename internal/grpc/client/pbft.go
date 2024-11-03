package client

import (
	"context"
	"io"
	"time"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.uber.org/zap"
)

// Commit calls the Commit RPC on the target machine (nodes to nodes).
func (c *Client) Commit(target string, msg *pbft.CommitMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call commit RPC
	if _, err := pbft.NewPBFTClient(conn).Commit(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call Commit RPC", zap.String("address", address), zap.Error(err))
	}
}

// PrePrepare calls the PrePrepare RPC on the target machine (nodes to nodes).
func (c *Client) PrePrepare(target string, msg *pbft.PrePrepareMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call preprepare RPC
	if _, err := pbft.NewPBFTClient(conn).PrePrepare(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call Preprepare RPC", zap.String("address", address), zap.Error(err))
	}
}

// PrePrepared calls the PrePrepared RPC on the target machine (nodes to nodes).
func (c *Client) PrePrepared(target string, msg *pbft.PrePreparedMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call preprepare RPC
	if _, err := pbft.NewPBFTClient(conn).PrePrepared(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call Preprepared RPC", zap.String("address", address), zap.Error(err))
	}
}

// Prepare calls the Prepare RPC on the target machine (nodes to nodes).
func (c *Client) Prepare(target string, msg *pbft.PrepareMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call prepare RPC
	if _, err := pbft.NewPBFTClient(conn).Prepare(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call Prepare RPC", zap.String("address", address), zap.Error(err))
	}
}

// Prepared calls the Prepared RPC on the target machine (nodes to nodes).
func (c *Client) Prepared(target string, msg *pbft.PreparedMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call prepare RPC
	if _, err := pbft.NewPBFTClient(conn).Prepared(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call Prepared RPC", zap.String("address", address), zap.Error(err))
	}
}

// Reply calls the Reply RPC on the target machine (nodes to clients).
func (c *Client) Reply(target string, msg *pbft.ReplyMsg) {
	address := c.clients[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call reply RPC
	if _, err := pbft.NewPBFTClient(conn).Reply(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call Reply RPC", zap.String("address", address), zap.Error(err))
	}
}

// Request calls the Request RPC on the target machine (clients to nodes).
func (c *Client) Request(target string, msg *pbft.RequestMsg) {
	address := c.nodes[target]
	msg.NodeId = c.nodeId
	msg.ClientId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))
		return
	}
	defer conn.Close()

	// call request RPC
	if _, err := pbft.NewPBFTClient(conn).Request(context.Background(), msg); err != nil {
		c.logger.Debug("failed to call Request RPC", zap.String("address", address), zap.Error(err))
	}
}

// PrintDB gets a target datastore.
func (c *Client) PrintDB(target string) []*pbft.RequestMsg {
	address := c.nodes[target]
	list := make([]*pbft.RequestMsg, 0)

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))

		return list
	}
	defer conn.Close()

	// open a stream on print db rpc
	stream, err := pbft.NewPBFTClient(conn).PrintDB(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.logger.Debug("failed to call PrintDB RPC", zap.String("address", address), zap.Error(err))

		return list
	}

	for {
		// get requests one by one
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				stream.CloseSend()
				break
			}
		}

		// append to the list of requests
		list = append(list, in)
	}

	return list
}

// PrintLog gets a target logs.
func (c *Client) PrintLog(target string) []string {
	address := c.nodes[target]
	list := make([]string, 0)

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))

		return list
	}
	defer conn.Close()

	// open a stream on print log rpc
	stream, err := pbft.NewPBFTClient(conn).PrintLog(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.logger.Debug("failed to call PrintLog RPC", zap.String("address", address), zap.Error(err))

		return list
	}

	for {
		// get requests one by one
		in, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				stream.CloseSend()
				break
			}
		}

		// append to the list of requests
		list = append(list, in.GetText())
	}

	return list
}

// PrintStatus gets a target and a sequence number to get an specific request status.
func (c *Client) PrintStatus(target string, sequenceNumber int) string {
	address := c.nodes[target]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))

		return err.Error()
	}
	defer conn.Close()

	// call print status RPC
	resp, err := pbft.NewPBFTClient(conn).PrintStatus(context.Background(), &pbft.StatusMsg{
		SequenceNumber: int64(sequenceNumber),
	})
	if err != nil {
		c.logger.Debug("failed to call PrintStatus RPC", zap.String("address", address), zap.Error(err))

		return err.Error()
	}

	// convert status enum to string
	switch resp.GetStatus() {
	case pbft.RequestStatus_REQUEST_STATUS_PP:
		return "preprepared"
	case pbft.RequestStatus_REQUEST_STATUS_P:
		return "prepared"
	case pbft.RequestStatus_REQUEST_STATUS_C:
		return "committed"
	case pbft.RequestStatus_REQUEST_STATUS_E:
		return "executed"
	default:
		return "no status"
	}
}

// TODO: print view
func (c *Client) PrintView(target string) bool {
	return false
}

// Transaction sends a transaction to one client.
func (c *Client) Transaction(sender, reseiver string, amount int) string {
	address := c.clients[sender]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		c.logger.Debug("failed to connect", zap.String("address", address), zap.Error(err))

		return err.Error()
	}
	defer conn.Close()

	// call transaction RPC
	resp, err := pbft.NewPBFTClient(conn).Transaction(context.Background(), &pbft.TransactionMsg{
		Sender:    sender,
		Reciever:  reseiver,
		Amount:    int64(amount),
		Timestamp: time.Now().Unix(),
	})
	if err != nil {
		c.logger.Debug("failed to call Transaction RPC", zap.String("address", address), zap.Error(err))

		return err.Error()
	}

	// extract and return message
	return resp.GetText()
}
