package chain

import (
	"bytes"

	logging "github.com/ipfs/go-log/v2"
	localCrypto "quorum/internal/pkg/crypto"
	"quorum/internal/pkg/nodectx"
	quorumpb "quorum/internal/pkg/pb"
)

var molautil_log = logging.Logger("util")

func Hash(data []byte) []byte {
	return localCrypto.Hash(data)
}