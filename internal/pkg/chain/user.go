package chain

import (
	quorumpb "quorum/internal/pkg/pb"
	"google.golang.org/protobuf/proto"
)

type User interface {
	Init(item *quorumpb.GroupItem, nodename string, iface ChainMolassesIface)
	UpdAnnounce(item *quorumpb.AnnounceItem) (string, error)
	UpdBlkList(item *quorumpb.DenyUserItem) (string, error)
	UpdSchema(item *quorumpb.SchemaItem) (string, error)
	UpdProducer(item *quorumpb.ProducerItem) (string, error)
	PostToGroup(content proto.Message) (string, error)
	AddBlock(block *quorumpb.Block) error
}
