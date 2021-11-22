package chain

import (
	"bytes"
	"encoding/hex"
	"errors"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"quorum/internal/pkg/nodectx"
	quorumpb "quorum/internal/pkg/pb"
	"google.golang.org/protobuf/proto"
)

const (
	USER_CHANNEL_PREFIX = "user_channel_"
	PRODUCER_CHANNEL_PREFIX = "prod_channel_"
)

type Group struct {
	// Group Item
	Item *quorumpb.GroupItem
	ChainCtx *Chain

}