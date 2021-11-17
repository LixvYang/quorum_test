package api

import (
	"errors"
	"quorum/internal/pkg/options"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	logging "github.com/ipfs/go-log/v2"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	logger = logging.Logger("api")
)

type TokenItem struct {
	Token string `json:"token"`
}

func getJWTKey(h *Handler) (string, error) {
	// get JWTKey from node options config file
	nodeOpt := options.GetNodeOptions()
	if nodeOpt == nil {
		return "", errors.New("Call InitNodeOptions() before use it")
	}
	return nodeOpt.JWTKey, nil
}

func getToken(name string, jwtKey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = name
	//FIXME:hardcode
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	return token.SignedString([]byte(jwtKey))
}


func CustomJWTConfig(jwtKey string) middleware.JWTConfig {
	config := middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:	[]byte(jwtKey),
		AuthScheme: "Bearer",
		TokenLookup: "header:" + echo.HeaderAuthorization,
		Skipper: func(c echo.Context) bool {
			r := c.Request()
			if strings.HasPrefix(r.Host, "localhost:") || r.Host == "localhost" || strings.HasPrefix(r.Host,"127.0.0.1") {
				return true
			} else if strings.HasPrefix(r.URL.Path,"/app/api/v1/token/apply") {
				// FIXME:hardcode url path
				return true
			} 
			return false
		},
	}
	return config
}