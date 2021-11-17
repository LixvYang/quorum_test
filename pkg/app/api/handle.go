package api

import (
	"context"

	"quorum/internal/pkg/appdata"
	"quorum/internal/pkg/storage"
)

type Handler struct {
	Ctx       context.Context
	Appdb     *appdata.AppDb
	Chaindb   *storage.DbMgr
	Apiroot   string
	GitCommit string
	ConfigDir string
	PeerName  string
	NodeName  string
}
