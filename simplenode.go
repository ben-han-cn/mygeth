package mygeth

import (
	"encoding/hex"
	clog "log"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/release"
	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/params"
)

const (
	clientIdentifier = "geth"
)

var (
	gitCommit = "1.7.0-unstable"
	relOracle = common.HexToAddress("0xfa7b9770ca4cb04296cac84f37736d4041251cdf")
)

func initLog() {
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(4))
	glogger.Vmodule("")
	log.Root().SetHandler(glogger)
}

func defaultNodeConfig(dataDir string) node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = "1.7.0-unstable"
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "shh")
	cfg.WSModules = append(cfg.WSModules, "eth", "shh")
	cfg.IPCPath = "geth.ipc" //no ipc
	cfg.DataDir = dataDir
	cfg.NoUSB = true
	setP2PConfig(&cfg.P2P)
	return cfg
}

func MakeFullNode(dataDir string) *node.Node {
	initLog()
	nodeConf := defaultNodeConfig(dataDir)
	stack, err := node.New(&nodeConf)
	if err != nil {
		clog.Fatalf("create node failed: %v", err)
	}

	ethConf := eth.DefaultConfig
	ethConf.SyncMode = downloader.FastSync
	ethConf.MaxPeers = 25
	ethConf.LightPeers = 20
	ethConf.DatabaseHandles = 1024
	ethConf.EthashDatasetDir = filepath.Join(nodeConf.DataDir, "ethash")
	err = stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		return eth.New(ctx, &ethConf)
	})
	if err != nil {
		clog.Fatalf("Failed to register the Ethereum service: %v", err)
	}

	err = stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		config := release.Config{
			Oracle: relOracle,
			Major:  uint32(params.VersionMajor),
			Minor:  uint32(params.VersionMinor),
			Patch:  uint32(params.VersionPatch),
		}
		commit, _ := hex.DecodeString(gitCommit)
		copy(config.Commit[:], commit)
		return release.NewReleaseService(ctx, config)
	})
	if err != nil {
		clog.Fatalf("Failed to register oracle: %v", err)
	}

	return stack
}
