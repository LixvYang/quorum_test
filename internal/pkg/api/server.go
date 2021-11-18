package api

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"

	"quorum/internal/pkg/cli"
	localcrypto "quorum/internal/pkg/crypto"
	"quorum/internal/pkg/options"
	"quorum/internal/pkg/p2p"
	"quorum/internal/pkg/utils"
	appapi "quorum/pkg/app/api"

	quorumpb "quorum/internal/pkg/pb"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/protobuf/encoding/protojson"
)

var quitch chan os.Signal

func StartAPIServer(config cli.Config,signalch chan os.Signal, h *Handler,node *p2p.Node,nodeopt *options.NodeOptions,ks  localcrypto.Keystore,ethaddr string,isbootstrapnode bool)  {
	quitch = signalch
	e := echo.New()
	e.Binder = new(CustomBinder)
	e.Use(middleware.JWTWithConfig(appapi.CustomJWTConfig(nodeopt.JWTKey)))
	r := e.Group("/api")
	a := e.Group("/app/api")
	r.GET("/quit",quitapp)
	 if isbootstrapnode == false {
		r.POST("/v1/group", h.CreateGroup())
	 }

	 certPath, keyPath, err := utils.GetTLSCerts()
	 if err != nil {
		panic(err)
	 }
	 e.Logger.Fatal(e.StartTLS(config.APIListenAddresses, certPath, keyPath))
}


type CustomBinder struct{}

func (cb *CustomBinder) Bind(i interface{}, c echo.Context) (err error) {
	db := new(echo.DefaultBinder)
	switch i.(type) {
	case *quorumpb.Activity:
		bodyBytes, err := ioutil.ReadAll(c.Request().Body)
		err = protojson.Unmarshal(bodyBytes,i.(*quorumpb.Activity))
		return err
	default:
		if err = db.Bind(i,c); err != echo.ErrUnsupportedMediaType {
			return
		}
		return err
	}
}

func quitapp(c echo.Context) (err error) {
	fmt.Println("/api/quit has been called, send Signal SIGNERM...")
	quitch <- syscall.SIGTERM
	return nil
}