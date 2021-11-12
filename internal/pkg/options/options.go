package options

import (
	"quorum/internal/pkg/utils"
	"sync"

	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/viper"
)

var optionslog = logging.Logger("options")

type NodeOptions struct {
	EnableNat        bool
	EnableDevNetwork bool
	NetworkName      string
	JWTToken         string
	JWTKey           string
	SignKeyMap       map[string]string
	mu               sync.Mutex
}

var nodeopts *NodeOptions
var nodeconfigdir string
var nodepeername string

const JWTKeyLength = 32
const defaultNetworkName = "nevis"

func GetNodeOptions() *NodeOptions {
	return nodeopts
}

func Load(dir string, keyname string) (*NodeOptions, error) {
	v, err := initConfigfile(dir, keyname)
	if err != nil {
		return nil, err
	}
	err = v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	options := &NodeOptions{}
	options.EnableNat = v.GetBool("EnableNat")
	options.EnableDevNetwork = v.GetBool("EnableDevNetwork")
	options.NetworkName = v.GetString("NetworkName")
	if options.NetworkName == "" {
		options.NetworkName = defaultNetworkName
	}

	options.SignKeyMap = v.GetStringMapString("SignKeyMap")
	options.JWTKey = v.GetString("JWTKey")
	options.JWTToken = v.GetString("JWTToken")
	return options, nil
}

func initConfigfile(dir, keyname string) (*viper.Viper, error) {
	if err := utils.EnsureDir(dir); err != nil {
		optionslog.Error("Check config directory failed: %s", err)
		return nil, err
	}

	v := viper.New()
	v.SetConfigFile(keyname + "_options.toml")
	v.SetConfigName(keyname + "_options")
	v.SetConfigType("toml")
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			optionslog.Infof("config file not found, generating...")
			writeDefaultToConfig(v)
		} else {
			return nil, err
		}
	}

	if v.GetString("JWTKey") == "" {
		v.Set("JWTKey", utils.GetRandomStr(JWTKeyLength))
		if err := v.WriteConfig(); err != nil {
			return nil ,err
		}
	}

	return v, nil
}

func writeDefaultToConfig(v *viper.Viper) error {
	v.Set("EnableNat", true)
	v.Set("EnableDevNetwork", false)
	v.Set("NetworkName", defaultNetworkName)
	v.Set("JWTKey", utils.GetRandomStr(JWTKeyLength))
	v.Set("JWTToken", "")
	v.Set("SignKeyMap", map[string]string{})
	return v.SafeWriteConfig()
}

// 初始化节点配置
func InitNodeOptions(configdir, peername string) (*NodeOptions, error) {
	var err error
	nodeopts, err := Load(configdir, peername)
	if err == nil {
		nodeconfigdir = configdir
		nodepeername = peername
	}
	return nodeopts, err
}


func (opt *NodeOptions) writeToconfig() error {
	v, err := initConfigfile(nodeconfigdir, nodepeername)
	if err != nil {
		return err
	}
	v.Set("EnableNat", opt.EnableNat)
	v.Set("EnableDevNetwork", opt.EnableDevNetwork)
	v.Set("SignKeyMap", opt.SignKeyMap)
	v.Set("JWTKey", opt.JWTKey)
	v.Set("JWTToken", opt.JWTToken)
	return v.WriteConfig()
}

func (opt *NodeOptions) SetJWTKey(jwtKey string) error {
	opt.mu.Lock()
	defer opt.mu.Unlock()
	opt.JWTKey = jwtKey
	return opt.writeToconfig()
}

func (opt *NodeOptions) SetJWTToken(jwtToken string) error {
	opt.mu.Lock()
	defer opt.mu.Unlock()
	opt.JWTToken = jwtToken
	return opt.writeToconfig()
}

func (opt *NodeOptions) SetSignKeyMap(keyname, addr string) error {
	opt.mu.Lock()
	defer opt.mu.Unlock()
	opt.SignKeyMap[keyname] = addr
	return opt.writeToconfig()
}