package api

import (
	"context"
	"quorum/internal/pkg/nodectx"
	"quorum/internal/pkg/p2p"
)

type (
	Handler struct {
		Ctx context.Context
		Node	*p2p.Node
		NodeCtx *nodectx.NodeCtx
		GitCommit string
	}
)