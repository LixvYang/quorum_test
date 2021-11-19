package api

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rumsystem/quorum/internal/pkg/options"
	quorumpb "quorum/internal/pkg/pb"

	guuid "github.com/google/uuid"
	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
	chain "quorum/internal/pkg/chain"
	localcrypto "quorum/internal/pkg/crypto"
	"quorum/internal/pkg/nodectx"
)

type CreateGroupParam struct {
	GroupName      string `from:"group_name" json:"group_name" validate:"required,max=20,min=5"`
	ConsensusType  string `from:"consensus_type"  json:"consensus_type"  validate:"required,oneof=pos poa"`
	EncryptionType string `from:"encryption_type" json:"encryption_type" validate:"required,oneof=public private"`
	AppKey         string `from:"app_key"         json:"app_key"         validate:"required,max=20,min=5"`
}

type CreateGroupResult struct {
	GenesisBlock       *quorumpb.Block `json:"genesis_block"`
	GroupId            string          `json:"group_id"`
	GroupName          string          `json:"group_name"`
	OwnerPubkey        string          `json:"owner_pubkey"`
	OwnerEncryptPubkey string          `json:"owner_encryptpubkey"`
	ConsensusType      string          `json:"consensus_type"`
	EncryptionType     string          `json:"encryption_type"`
	CipherKey          string          `json:"cipher_key"`
	AppKey             string          `json:"app_key"`
	Signature          string          `json:"signature"`
}

// @Tag Groups
// @Summary CreateGroup
// @Description Create a new group
// @Accept json
// @Produce json
// @Param data body CreateGroupParam true "GroupInfo"
// Success 200 {object} CreateGroupResult
// @Router /api/v1/group [post]
func (h *Handler) CreateGroup() echo.HandlerFunc {
	return func(c echo.Context) error {
		var err error
		output := make(map[string]string)
		
		validate := validate.New()
		params := new(CreateGroupParam)
		if err = c.Bind(params); err != nil {
			output[ERROR_INFO] = err.Error()
			return c.JSON(http.StatusBadRequest, output)
		}

		if err = validate.Struct(params); err != nil {
			output[ERROR_INFO] = err.Error()
			return c.JSON(http.StatusBadRequest, output)
		}

		if params.ConsensusType != "poa" {
			output := "Other types of groups are not supported yet"
			return c.JSON(http.StatusBadRequest, output)
		}

		groupid := guuid.New()

		nodeoptions := options.GetNodeOptions()

		var groupSignPubkey []byte
		var p2ppubkey p2pcrypto.PubKey
		ks := nodectx.GetNodeCtx().Keystore
		dirks, ok := ks.(*localcrypto.DirKeyStore)
		if ok == true {
			hexkey, err := dirks.GetEncodedPubkey(groupid.String(),localcrypto.Sign)
			if err != nil && strings.HasPrefix(err.Error(), "key not exist ") {
				newsignaddr, err := dirks.NewKeyWithDefaultPassword(groupid.String(),localcrypto.Sign)
				if err == nil && newsignaddr != "" {
					err = nodeoptions.SetSignKeyMap(groupid.String(),newsignaddr)
					if err != nil {
						output[ERROR_INFO] = fmt.Sprintf("save key map %s err : %s",newsignaddr,err.Error())
						return c.JSON(http.StatusBadRequest, output)
					}
				}
				hexkey, err = dirks.GetEncodedPubkey(guuidid.String(),localcrypto.Sign)
			} else {
				output[ERROR_INFO] = "Create new group key err:" + err.Error()
				return c.JSON(http.StatusBadRequest,output)
			}
		}
	}
}

