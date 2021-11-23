package pubsubconn

import (
	quorumpb "quorum/internal/pkg/pb"
)

type Chain interface {
	HandleTrx(trx *quorumpb.Trx) error
	HandleBlock(block *quorumpb.Block) error
}

type PubSubConn interface {
	JoinChannel(cId string, chain Chain) error
	Publish(data []byte) error
}
