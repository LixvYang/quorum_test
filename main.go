package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"quorum/internal/pkg/cli"
	"quorum/internal/pkg/p2p"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
)

var (
	ReleaseVersion string
	GitCommit      string
)

func main() {
	if ReleaseVersion == "" {
		ReleaseVersion = "v1.0.0"
	}
	if GitCommit == "" {
		GitCommit = "devel"
	}

	help := flag.Bool("h", false, "Display Help")
	version := flag.Bool("version", false, "Show the version")
	update := flag.Bool("update", false, "Update to the latest version")
	// updateFrom := flag.String("from", "github", "Update from: github/qingcloud, default to github")
	config, err := cli.ParseFlags()
	lvl, err := logging.LevelFromString("info")
	logging.SetAllLoggers(lvl)
	logging.SetLogLevel("appsync", "error")
	logging.SetLogLevel("appdata", "error")
	if err != nil {
		panic(err)
	}

	if config.IsDebug {
		logging.SetLogLevel("main", "debug")
		logging.SetLogLevel("crypto", "debug")
		logging.SetLogLevel("network", "debug")
		logging.SetLogLevel("pubsub", "debug")
		logging.SetLogLevel("autonat", "debug")
		logging.SetLogLevel("chain", "debug")
		logging.SetLogLevel("dbmgr", "debug")
		logging.SetLogLevel("chainctx", "debug")
		logging.SetLogLevel("group", "debug")
		logging.SetLogLevel("syncer", "debug")
		logging.SetLogLevel("producer", "debug")
		logging.SetLogLevel("user", "debug")
		logging.SetLogLevel("groupmgr", "debug")
	}

	if *help {
		fmt.Println("Output a help ")
		fmt.Println()
		fmt.Println("Usage:...")
		flag.PrintDefaults()
	}

	if *version {
		fmt.Printf("%s - %s\n", ReleaseVersion, GitCommit)
		return
	}

	if *update {
		return
	}

	if config.IsPing {
		if len(config.BootstrapPeers) == 0 {
			fmt.Println("Usage: ", os.Args[0], "-ping", "-peer <peer> [-peer <peer> ...]")
		}

		//FIXME: hardcode
		tcpAddr := "/ip4/127.0.0.1/tcp/0"
		wsAddr := "/ip4/127.0.0.1/tcp/0/ws"
		ctx := context.Background()
		node, err := libp2p.New(
			ctx,
			libp2p.ListenAddrStrings(tcpAddr, wsAddr),
			libp2p.Ping(false),
		)
		if err != nil {
			panic(err)
		}

		// configure our ping protrol
		pingService := &p2p.PingService{ Host:node }
		node.SetStreamHandler(p2p.PingID,pingService.PingHandler)
	}

	

}
