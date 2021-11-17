package main

import (
	logging "github.com/ipfs/go-log/v2"
	"os"
)

var (
	signalch chan os.Signal
	mainlog  = logging.Logger("main")
)