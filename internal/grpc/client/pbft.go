package client

import (
	"context"
	"io"

	"github.com/f24-cse535/pbft/pkg/rpc/pbft"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Commit calls the Commit RPC on the target machine (nodes to nodes).
func (c *Client) Commit(target string, msg *pbft.AckMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call commit RPC
	if _, err := pbft.NewPBFTClient(conn).Commit(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// PrePrepare calls the PrePrepare RPC on the target machine (nodes to nodes).
func (c *Client) PrePrepare(target string, msg *pbft.PrePrepareMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call preprepare RPC
	if _, err := pbft.NewPBFTClient(conn).PrePrepare(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// PrePrepared calls the PrePrepared RPC on the target machine (nodes to nodes).
func (c *Client) PrePrepared(target string, msg *pbft.AckMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call preprepare RPC
	if _, err := pbft.NewPBFTClient(conn).PrePrepared(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// Prepare calls the Prepare RPC on the target machine (nodes to nodes).
func (c *Client) Prepare(target string, msg *pbft.AckMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call prepare RPC
	if _, err := pbft.NewPBFTClient(conn).Prepare(context.Background(), msg); err != nil {
		return err
	}

	return err
}

// Prepared calls the Prepared RPC on the target machine (nodes to nodes).
func (c *Client) Prepared(target string, msg *pbft.AckMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call prepare RPC
	if _, err := pbft.NewPBFTClient(conn).Prepared(context.Background(), msg); err != nil {
		return err
	}

	return err
}

// Request calls the Request RPC on the target machine (clients to nodes).
func (c *Client) Request(target string, msg *pbft.RequestMsg) error {
	address := c.nodes[target]

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call request RPC
	if _, err := pbft.NewPBFTClient(conn).Request(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// ViewChange calls the view change RPC on the target machine (nodes to nodes).
func (c *Client) ViewChange(target string, msg *pbft.ViewChangeMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call view change RPC
	if _, err := pbft.NewPBFTClient(conn).ViewChange(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// NewView calls the new view RPC on the target machine (nodes to nodes).
func (c *Client) NewView(target string, msg *pbft.NewViewMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call new view RPC
	if _, err := pbft.NewPBFTClient(conn).NewView(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// Checkpoint calls the new checkpoint RPC on the target machine.
func (c *Client) Checkpoint(target string, msg *pbft.CheckpointMsg) error {
	address := c.nodes[target]
	msg.NodeId = c.nodeId

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	// call checkpoint RPC
	if _, err := pbft.NewPBFTClient(conn).Checkpoint(context.Background(), msg); err != nil {
		return err
	}

	return nil
}

// PrintDB gets a target datastore.
func (c *Client) PrintDB(target string) []*pbft.RequestRsp {
	address := c.nodes[target]
	list := make([]*pbft.RequestRsp, 0)

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return list
	}
	defer conn.Close()

	// open a stream on print db rpc
	stream, err := pbft.NewPBFTClient(conn).PrintDB(context.Background(), &emptypb.Empty{})
	if err != nil {
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
		return list
	}
	defer conn.Close()

	// open a stream on print log rpc
	stream, err := pbft.NewPBFTClient(conn).PrintLog(context.Background(), &emptypb.Empty{})
	if err != nil {
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
		return err.Error()
	}
	defer conn.Close()

	// call print status RPC
	resp, err := pbft.NewPBFTClient(conn).PrintStatus(context.Background(), &pbft.StatusMsg{
		SequenceNumber: int64(sequenceNumber),
	})
	if err != nil {
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

// PrintView gets a target to print the view change messages.
func (c *Client) PrintView(target string) []*pbft.ViewRsp {
	address := c.nodes[target]
	list := make([]*pbft.ViewRsp, 0)

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return list
	}
	defer conn.Close()

	// open a stream on print view rpc
	stream, err := pbft.NewPBFTClient(conn).PrintView(context.Background(), &emptypb.Empty{})
	if err != nil {
		return list
	}

	for {
		// get messages one by one
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

// PrintCheckpoints gets a target to print checkpoint messages.
func (c *Client) PrintCheckpoints(target string) []*pbft.CheckpointRsp {
	address := c.nodes[target]
	list := make([]*pbft.CheckpointRsp, 0)

	// base connection
	conn, err := c.connect(address)
	if err != nil {
		return list
	}
	defer conn.Close()

	// open a stream on print checkpoints rpc
	stream, err := pbft.NewPBFTClient(conn).PrintCheckpoints(context.Background(), &emptypb.Empty{})
	if err != nil {
		return list
	}

	for {
		// get messages one by one
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
