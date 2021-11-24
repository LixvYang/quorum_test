package chain

import (
	quorumpb "quorum/internal/pkg/pb"
)

type Producer interface {
	Init(item *quorumpb.GroupItem, nodename string, iface ChainMolassesIface)
	AddTrx(trx *quorumpb.Trx)
	AddBlockToPool(block *quorumpb.Block)
	GetBlockForward(trx *quorumpb.Trx) error
	GetBlockBackward(trx *quorumpb.Trx) error
	AddProducedBlock(trx *quorumpb.Trx) error
	AddBlock(block *quorumpb.Block) error
}
