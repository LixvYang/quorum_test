package chain

import quorumpb "quorum/internal/pkg/pb"

type ChainMolassesIface interface {
	GetUserTrxMgr() *TrxMgr
	GetProducerTrxMgr() *TrxMgr
	UpdChainInfo(height int64, blockId string) error
	UpdProducerList()
	CreateConsensus()
	IsSyncerReady() bool
	SyncBackward(block *quorumpb.Block) error
}
