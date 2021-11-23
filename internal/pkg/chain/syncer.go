package chain

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"quorum/internal/pkg/nodectx"
	quorumpb "quorum/internal/pkg/pb"

	logging "github.com/ipfs/go-log/v2"
)

var syncer_log = logging.Logger("syncer")

var WAIT_BLOCK_TIME_S = 10
var RETRY_LIMIT = 30

const (
	SYNCING_FORWARD  = 0
	SYNCING_BACKWARD = 1
	SYNC_FAILED      = 2
	IDLE             = 3
)

type Syncer struct {
	nodeName         string
	group            *Group
	trxMgr           *TrxMgr
	AskNextTimer     *time.Timer
	AskNextTimerDone chan bool
	Status           int8
	retryCount       int8
	statusBeforeFail int8
	responses        map[string]*quorumpb.ReqBlockResp
	groupId          string
}

func (syncer *Syncer) Init(grp *Group, trxMgr *TrxMgr)  {
	syncer_log.Debug("Init called")
	syncer.Status = IDLE
	syncer.group = grp
	syncer.trxMgr = trxMgr
	syncer.retryCount = 0
	syncer.responses = make(map[string]*quorumpb.ReqBlockResp)
	syncer.groupId = grp.Item.GroupId
	syncer_log.Infof("<%s> syncer initialed", syncer.groupId)
}


