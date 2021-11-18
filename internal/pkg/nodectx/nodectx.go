package nodectx

import (
	"context"
	"quorum/internal/pkg/p2p"
	"quorum/internal/pkg/storage"

	localcrypto "quorum/internal/pkg/crypto"

	logging "github.com/ipfs/go-log/v2"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/peer"
)

var chainctx_log = logging.Logger("chainctx")

type NodeStatus int8

const (
	USER_CHANNEL_PREFIX     = "user_channel_"
	PRODUCER_CHANNEL_PREFIX = "prod_channel_"
)

const (
	NODE_ONLINE  = 0
	NODE_OFFLINE = 1
)

type NodeCtx struct {
	Node      *p2p.Node
	PeerId    peer.ID
	Keystore  localcrypto.Keystore
	PublicKey p2pcrypto.PubKey
	Name      string
	Ctx       context.Context
	Version   string
	Status    NodeStatus
}

var nodeCtx *NodeCtx

var dbMgr *storage.DbMgr

// singlaton
func GetNodeCtx() *NodeCtx {
	return nodeCtx
}

// singlaton
func GetDbMgr() *storage.DbMgr {
	return dbMgr
}

func InitCtx(ctx context.Context, name string, node *p2p.Node, db *storage.DbMgr, channeltype string, gitcommit string) {
	nodeCtx := &NodeCtx{}
	nodeCtx.Name = name
	nodeCtx.Node = node

	dbMgr = db

	nodeCtx.Status = NODE_OFFLINE
	nodeCtx.Ctx = ctx
	nodeCtx.Version = "1.0.0"
}

func (nodeCtx *NodeCtx) PeersProtocol() *map[string][]string {
	return nodeCtx.Node.PeersProtocol()
}

func (nodeCtx *NodeCtx) ProtocolPrefix() string {
	return p2p.ProtocolPrefix
}

func (nodeCtx *NodeCtx) UpdateOnlineStatus(status NodeStatus) {
	nodeCtx.Status = status
}

func (nodeCtx *NodeCtx) GetNodePubKey() (string, error) {
	var pubkey string
	pubkeybytes, err := p2pcrypto.MarshalPublicKey(nodeCtx.PublicKey)
	if err != nil {
		return pubkey, err
	}

	pubkey = p2pcrypto.ConfigEncodeKey(pubkeybytes)
	return pubkey, nil
}

func (nodeCtx *NodeCtx) ListGroupPeers(groupid string) []peer.ID {
	userChannelId := USER_CHANNEL_PREFIX + groupid
	return nodeCtx.Node.Pubsub.ListPeers(userChannelId)
}

func (nodeCtx *NodeCtx) AddPeers(peers []peer.AddrInfo) int {
	return nodeCtx.Node.AddPeers(nodeCtx.Ctx, peers)
}
