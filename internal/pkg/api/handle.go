package api

import (
	"context"
	"quorum/internal/pkg/appdata"
	"quorum/internal/pkg/nodectx"
	"quorum/internal/pkg/p2p"
)

type (
	Handler struct {
		Ctx       context.Context
		Node      *p2p.Node
		NodeCtx   *nodectx.NodeCtx
		GitCommit string
		Appdb     *appdata.AppDb
	}
)

type ErrorResponse struct {
	Error string `json:"error" validate:"required"`
}
