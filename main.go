package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"path"

	"quorum/internal/pkg/cli"
	"quorum/internal/pkg/p2p"
	"quorum/internal/pkg/utils"
	"quorum/internal/pkg/storage"
	"quorum/internal/pkg/options"
	localcrypto "quorum/internal/pkg/crypto"
	ethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"

	dsbadger2 "github.com/ipfs/go-ds-badger2"
	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p"
	peerstore "github.com/libp2p/go-libp2p-core/peer"
)

const DEFAUT_KEY_NAME string = "default"


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

	signkeycount, err := localcrypto.InitKeystore(config.KeyStoreName, config.KeyStoreDir)
	ksi := localcrypto.GetKeystore()
	if err != nil {
		cancel()
		mainlog.Fatalf(err.Error())
	}
	ks, ok := ksi.(*localcrypto.DirKeyStore)
	if !ok {
		//TODO: test other keystore type?
		//if there are no other keystores, exit and show error info.
		cancel()
		mainlog.Fatalf(err.Error())
	}

	password := os.Getenv("RUM_KSPASSWD")
	if signkeycount > 0 {
		if password == "" {
			password, err := localcrypto.PassphrasePromptForUnlock()
		}
		err = ks.Unlock(nodeoptions.SignKeyMap,password)
		if err != nil {
			mainlog.Fatalf(err.Error())
			cancel()
			return 0
		}
	} else {
		if password == "" {
			password, err = localcrypto.PassphrasePromptForEncryption()
			if err != nil {
				mainlog.Fatalf(err.Error())
				cancel()
				return 0
			}
			fmt.Println("Please keeping your password safe, We can't recover or reset your password.")
			fmt.Println("Your password:", password)
			fmt.Println("After saving the password, press any key to continue.")
			os.Stdin.Read(make([]byte, 1))
		}

		signkeyhexstr, err := localcrypto.LoadEncodeKeyFrom(config.ConfigDir, peername, "txt")
		if err != nil {
			cancel()
			mainlog.Fatalf(err.Error())
		}
		var addr string
		if signkeyhexstr != "" {
			addr, err = ks.Import(DEFAUT_KEY_NAME, signkeyhexstr, localcrypto.Sign, password)
		} else {
			addr, err = ks.NewKey(DEFAUT_KEY_NAME,localcrypto.Sign, password)
			if err != nil {
				mainlog.Fatalf(err.Error())
				cancel()
				return 0
			}
		}

		if addr == "" {
			mainlog.Fatalf("Load or create new signkey failed")
			cancel()
			return 0
		}

		err = nodeoptions.SetSignKeyMap(DEFAUT_KEY_NAME, addr)
		if err != nil {
			mainlog.Fatalf(err.Error())
			cancel()
			return 0
		}
		err = ks.Unlock(nodeoptions.SignKeyMap, password)
		if err != nil {
			mainlog.Fatalf(err.Error())
			cancel()
			return 0
		}

		fmt.Printf("load signkey: %d press any key to continue...\n", signkeycount)
	}

	_, err = ks.GetKeyFromUnlocked(localcrypto.Sign.NameString(DEFAUT_KEY_NAME))
	signkeycount = ks.UnlockedKeyCount(localcrypto.Sign)
	signkeycount = ks.UnlockedKeyCount(localcrypto.Sign)
	if signkeycount == 0 {
		mainlog.Fatalf("load signkey error, exit... %s", err)
		cancel()
		return 0
	}

	//Load default sign keys
	key, err := ks.GetKeyFromUnlocked(localcrypto.Sign.NameString(DEFAUT_KEY_NAME))

	defaultkey, ok := key.(*ethkeystore.Key)
	if ok == false {
		fmt.Println("load default key error, exit...")
		mainlog.Fatalf(err.Error())
		cancel()
		return 0
	}
	keys, err := localcrypto.SignKeytoPeerKeys(defaultkey)

	if err != nil {
		mainlog.Fatalf(err.Error())
		cancel()
		return 0
	}

	peerid, ethaddr, err := ks.GetPeerInfo(DEFAUT_KEY_NAME)
	if err != nil {
		cancel()
		mainlog.Fatalf(err.Error())
	}

	mainlog.Infof("eth addresss: <%s>", ethaddr)
	ds, err := dsbadger2.NewDatastore(path.Join(config.DataDir, fmt.Sprintf("%s-%s", peername, "peerstore")), &dsbadger2.DefaultOptions)
	checkLockError(err)
	if err != nil {
		cancel()
		mainlog.Fatalf(err.Error())
	}

	if config.IsBootstrap == true {
		listenaddresses, _ := utils.StringsToAddrs([]string{config.ListenAddresses})
		// bootstrop node connections: low watermarks: 1000 hi watermarks 50000, grace 30s
		

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

