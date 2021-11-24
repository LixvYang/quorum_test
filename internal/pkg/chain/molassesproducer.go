package chain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	guuid "github.com/google/uuid"
	logging "github.com/ipfs/go-log/v2"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	localcrypto "quorum/internal/pkg/crypto"
	"quorum/internal/pkg/nodectx"
	quorumpb "quorum/internal/pkg/pb"
	"google.golang.org/protobuf/proto"
)

var molaproducer_log = logging.Logger("producer")

const PRODUCE_TIMER time.Duration = 5
const MERGE_TIMER time.Duration = 5

const TRXS_TOTAL_SIZE int = 900 * 1024

type ProducerStatus int

const (
	StatusIdle ProducerStatus = iota
	StatusMerging
	StatusProducing
)

type MolassesProducer struct {
	grpItem *quorumpb.GroupItem
	blockPool map[string]*quorumpb.Block
	trxPool map[string]*quorumpb.Trx
	trxMgr map[string]*TrxMgr
	status ProducerStatus
	ProduceTimer *time.Timer
	ProduceDone chan bool
	statusmu sync.RWMutex
	nodename string
	cIface ChainMolassesIface
	groupId string
}