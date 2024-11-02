package models

import "github.com/f24-cse535/pbft/pkg/rpc/pbft"

type Log struct {
	Request      *pbft.RequestMsg
	PrePrepareds []*pbft.PrePreparedMsg
	Prepareds    []*pbft.PreparedMsg
}
