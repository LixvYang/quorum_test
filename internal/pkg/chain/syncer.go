package chain

import (
	"errors"
	"fmt"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"quorum/internal/pkg/nodectx"
	quorumpb "quorum/internal/pkg/pb"
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
