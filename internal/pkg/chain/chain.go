package chain

import (
	"encoding/hex"
	"errors"
	"sync"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"quorum/internal/pkg/nodectx"
	quorumpb "quorum/internal/pkg/pb"
	pubsubconn "quorum/internal/pkg/pubsubconn"
	"google.golang.org/protobuf/proto"

	localcrypto "quorum/internal/pkg/crypto"
)

const (
	USER_CHANNEL_PREFIX = "user_channel_"
	PRODUCER_CHANNEL_PREFIX = "prod_channel_"
)

var chain_log = logging.Logger("chain")

type GroupProducer struct {
	ProducerPubkey   string
	ProducerPriority int8
}

type Chain struct {
	nodename          string
	group             *Group
	userChannelId     string
	producerChannelId string
	trxMgrs           map[string]*TrxMgr
	ProducerPool      map[string]*quorumpb.ProducerItem

	Syncer    *Syncer
	Consensus Consensus
	statusmu  sync.RWMutex

	groupId string
}

func (chain *Chain) CustomInit(nodename string, group *Group, producerPubsubconn pubsubconn.P2pPubSubConn, userPubsubconn pubsubconn.PubSubConn)  {
	chain.group = group
	chain.trxMgrs = make(map[string]*TrxMgr)
	chain.nodename = nodename

	chain.producerChannelId = PRODUCER_CHANNEL_PREFIX + group.Item.GroupId
	producerTrxMgr := &TrxMgr{}
	producerTrxMgr.Init(chain.group.Item, &producerPubsubconn)
	producerTrxMgr.SetNodeName(nodename)
	chain.trxMgrs[chain.producerChannelId] = producerTrxMgr

	chain.Consensus = NewMolasses(&MolassesProducer{},MolassesUser{})


}