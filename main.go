package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"quorum/internal/pkg/cli"
	"quorum/internal/pkg/p2p"
	"quorum/internal/pkg/utils"
	"quorum/internal/pkg/storage"
	"quorum/internal/pkg/options"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

const DEFAULT_KEY_NAME = "default"

var (
	ReleaseVersion string
	GitCommit      string
	mainlog = logging.Logger("main")
)

// return EBUSY if LOCK is exist
func checkLockError(err error)  {
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "Another process is using this Badger database.") {
			mainlog.Errorf(errStr)
			os.Exit(16)
		}
	}
}

func createDb(path string) (*storage.DbMgr, error) {
	var err error
	groupDb := storage.QSBadger{}
	dataDb := storage.QSBadger{}
	err = groupDb.Init(path + "_groups")
	if err != nil {
		return nil, err
	}

	err = dataDb.Init(path + "_db")
	if err != nil {
		return nil, err
	}

	manager := storage.DbMgr{&groupDb, &dataDb, nil, path}
	return &manager, nil
}

func mainRet(config cli.Config)  {
	signalch = make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mainlog.Infof("Version: %s", GitCommit)
	peername := config.PeerName

	if config.IsBootstrap == true {
		peername = "bootstrap"
	}

	//Load node options
	nodeoptions, err := options.InitNodeOptions(config.ConfigDir, peername)
	if err != nil {
		cancel()
		mainlog.Fatalf(err.Error())
	}

	

}


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

		for _, addr := range config.BootstrapPeers {
			peer, err := peerstore.AddrInfoFromP2pAddr(addr)
			if err != nil {
				panic(err)
			}

			if err := node.Connect(ctx,*peer); err != nil {
				panic(err)
			}

			ch := pingService.Ping(ctx, peer.ID)
			fmt.Println()
			fmt.Println("pinging remote peer at",addr)
			for i := 0; i < 4; i ++ {
				res := <- ch
				fmt.Println("PING",addr,"in",res.RTT)
			}
		}

		return
	}

	// check dir
	if err := utils.EnsureDir(config.ConfigDir); err != nil {
		panic(err)
	}
	
	_, _, err = utils.NewTLSCert()
	if err != nil {
		panic(err)
	}

	os.Exit(1)
}

