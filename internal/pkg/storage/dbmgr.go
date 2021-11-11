package storage

import (
	"errors"
	"fmt"

	quorumpb "quorum/internal/pkg/pb"

	"github.com/gogo/protobuf/proto"
	logging "github.com/ipfs/go-log/v2"
)

var dbmgr_log = logging.Logger("dbmgr")


const TRX_PREFIX string = "trx" //trx
const BLK_PREFIX string = "blk" //block
const SEQ_PREFIX string = "seq" //sequence
const GRP_PREFIX string = "grp" //group
const CNT_PREFIX string = "cnt" //content
const ATH_PREFIX string = "ath" //auth
const PRD_PREFIX string = "prd" //producer
const ANN_PREFIX string = "ann" //announce
const SMA_PREFIX string = "sma" //schema
const CHD_PREFIX string = "chd" //cached


type DbMgr struct {
	GroupInfoDb	QuorumStorage
	Db	QuorumStorage
	Auth	QuorumStorage
	DataPath	string
}

func (dbMgr *DbMgr) CloseDb()  {
	dbMgr.GroupInfoDb.Close()
	dbMgr.Db.Close()
	dbmgr_log.Infof("ChainCtx Db closed")
}

func getPrefix(prefix ...string) string {
	nodeprefix := ""
	if len(prefix) == 1 {
		nodeprefix = prefix[0] + "_"
	}
	return nodeprefix
}

// save trx
func (dbMgr *DbMgr) AddTrx(trx *quorumpb.Trx, prefix ...string) error {
	nodeprefix := getPrefix(prefix...)

	key := nodeprefix + TRX_PREFIX + "_" + trx.TrxId
	value, err := proto.Marshal(trx)
	if err != nil {
		return err
	}
	return dbMgr.Db.Set([]byte(key), value)
}

// UNUSED
// rm Trx
func (dbMgr *DbMgr) RmTrx(trxId string, prefix ...string) error {
	nodeprefix := getPrefix(prefix...)
	key := nodeprefix + TRX_PREFIX + "_" + trxId
	return dbMgr.Db.Delete([]byte(key))
}

// get trx
func (dbMgr *DbMgr) GetTrx(trxId string, prefix ...string) (*quorumpb.Trx, error) {
	nodeprefix := getPrefix(prefix...)
	key := nodeprefix + TRX_PREFIX + "_" + trxId
	value, err := dbMgr.Db.Get([]byte(key))
	if err != nil {
		return nil, err
	}

	var trx quorumpb.Trx
	err = proto.Unmarshal(value, &trx)
	if err != nil {
		return nil, err
	}

	return &trx, err
}

func (dbMgr *DbMgr) UpdTrx(trx *quorumpb.Trx, prefix ...string) error {
	return dbMgr.AddTrx(trx, prefix...)
}

func (dbMgr *DbMgr) IsTrxExist(trxId string, prefix ...string) (bool, error) {
	nodeprefix := getPrefix(prefix...)
	key := nodeprefix + TRX_PREFIX + "_" + trxId
	return dbMgr.Db.IsExist([]byte(key))
}

func (dbMgr *DbMgr) AddGensisBlock(gensisBlock *quorumpb.Block, prefix ...string) error {
	nodeprefix := getPrefix(prefix...)
	key := nodeprefix + BLK_PREFIX + "_" + gensisBlock.BlockId
	
	chunk := quorumpb.BlockDbChunk{}
	chunk.BlockId = gensisBlock.BlockId
	chunk.BlockItem = gensisBlock
	chunk.ParentBlockId = ""
	chunk.Height = 0

	value, err := proto.Marshal(&chunk)
	if err != nil {
		return err
	}

	return dbMgr.Db.Set([]byte(key), value)
}

func (dbMgr *DbMgr) IsBlockExist(blockId string, cached bool, prefix ...string) (bool, error) {
	nodeprefix := getPrefix(prefix...)
	var key string

	if cached {
		key = nodeprefix + CHD_PREFIX + "_" + BLK_PREFIX + "_" + blockId
	} else {
		key = nodeprefix + BLK_PREFIX + "_" + blockId
	}

	return dbMgr.Db.IsExist([]byte(key))
}

