package launcher

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
	"github.com/galaxy-foundation/icicb-base/abft"
	"github.com/galaxy-foundation/icicb-base/hash"
	"github.com/galaxy-foundation/icicb-base/utils/cachescale"
	"github.com/naoina/toml"
	"gopkg.in/urfave/cli.v1"

	"github.com/goicicb/evmcore"

	galaxy "github.com/goicicb/galaxy"
	"github.com/goicicb/galaxy/genesisstore"
	"github.com/goicicb/gossip"
	"github.com/goicicb/gossip/gasprice"
	"github.com/goicicb/integration"
	"github.com/goicicb/integration/makegenesis"
	futils "github.com/goicicb/utils"
	"github.com/goicicb/vecmt"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(nodeFlags, testFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}
	checkConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(checkConfig),
		Name:        "checkconfig",
		Usage:       "Checks configuration file",
		ArgsUsage:   "",
		Flags:       append(nodeFlags, testFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The checkconfig checks configuration file.`,
	}
	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}

	// DataDirFlag defines directory to store Lachesis state and user's wallets
	DataDirFlag = utils.DirectoryFlag{
		Name:  "datadir",
		Usage: "Data directory for the databases and keystore",
		Value: utils.DirectoryString(DefaultDataDir()),
	}

	CacheFlag = cli.IntFlag{
		Name:  "cache",
		Usage: "Megabytes of memory allocated to internal caching",
		Value: DefaultCacheSize,
	}

	// GenesisFlag specifies network genesis configuration
	GenesisFlag = cli.StringFlag{
		Name:  "genesis",
		Usage: "'path to genesis file' - sets the network genesis configuration.",
	}

	RPCGlobalGasCapFlag = cli.Uint64Flag{
		Name:  "rpc.gascap",
		Usage: "Sets a cap on gas that can be used in icicb_call/estimateGas (0=infinite)",
		Value: gossip.DefaultConfig(cachescale.Identity).RPCGasCap,
	}
	RPCGlobalTxFeeCapFlag = cli.Float64Flag{
		Name:  "rpc.txfeecap",
		Usage: "Sets a cap on transaction fee (in ICICB) that can be sent via the RPC APIs (0 = no cap)",
		Value: gossip.DefaultConfig(cachescale.Identity).RPCTxFeeCap,
	}

	AllowedGalaxyGenesisHashes = map[uint64]hash.Hash{
		galaxy.MainNetworkID: hash.HexToHash("0x9510812d06bc3ea1097ecded6500e795b8b376f6f22e738065f7b6c6b218a09e"), // real 0xa97be35f423207258c18624416d67950933456c9a549585b9517c2d81c42a0ce
		galaxy.TestNetworkID: hash.HexToHash("0x322f536c852ab9f8bc2cdc8fb15d5216d9cc56a0d9136f00b58292b159a5142a"), // real 0x0774ccb0a820e486c3f435b805240f0ab8ceb79bf19af9617b548a28ce6ff9d2
	}
)

const (

	// DefaultCacheSize is calculated as memory consumption in a worst case scenario with default configuration
	// Average memory consumption might be 3-5 times lower than the maximum
	DefaultCacheSize  = 3200
	ConstantCacheSize = 1024
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		return fmt.Errorf("field '%s' is not defined in %s", field, rt.String())
	},
}

type config struct {
	Node          node.Config
	Galaxy        gossip.Config
	GalaxyStore   gossip.StoreConfig
	Lachesis      abft.Config
	LachesisStore abft.StoreConfig
	VectorClock   vecmt.IndexConfig
}

func (c *config) AppConfigs() integration.Configs {
	return integration.Configs{
		Galaxy:         c.Galaxy,
		GalaxyStore:    c.GalaxyStore,
		Lachesis:       c.Lachesis,
		LachesisStore:  c.LachesisStore,
		VectorClock:    c.VectorClock,
		AllowedGenesis: AllowedGalaxyGenesisHashes,
	}
}

func loadAllConfigs(file string, cfg *config) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	if err != nil {
		return errors.New(fmt.Sprintf("TOML config file error: %v.\n"+
			"Use 'dumpconfig' command to get an example config file.\n"+
			"If node was recently upgraded and a previous network config file is used, then check updates for the config file.", err))
	}
	return err
}

func getGalaxyGenesis(ctx *cli.Context) integration.InputGenesis {

	var genesis integration.InputGenesis
	switch {
	case ctx.GlobalIsSet(FakeNetFlag.Name):
		_, num, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", FakeNetFlag.Name, "err", err)
		}
		fakeGenesisStore := makegenesis.FakeGenesisStore(num, futils.ToIcicb(10000000000), futils.ToIcicb(5000000))
		genesis = integration.InputGenesis{
			Hash: fakeGenesisStore.Hash(),
			Read: func(store *genesisstore.Store) error {
				buf := bytes.NewBuffer(nil)
				err = fakeGenesisStore.Export(buf)
				if err != nil {
					return err
				}
				return store.Import(buf)
			},
			Close: func() error {
				return nil
			},
		}
	case ctx.GlobalIsSet(GenesisFlag.Name):
		genesisPath := ctx.GlobalString(GenesisFlag.Name)

		genesisFile, err := os.Open(genesisPath)
		if err != nil {
			utils.Fatalf("Failed to open genesis file: %v", err)
		}
		inputGenesisHash, readGenesisStore, err := genesisstore.OpenGenesisStore(genesisFile)
		if err != nil {
			utils.Fatalf("Failed to read genesis file: %v", err)
		}

		genesis = integration.InputGenesis{
			Hash:  inputGenesisHash,
			Read:  readGenesisStore,
			Close: genesisFile.Close,
		}
	default:
		log.Crit("Network genesis is not specified")
	}
	return genesis
}

func setBootnodes(ctx *cli.Context, urls []string, cfg *node.Config) {
	cfg.P2P.BootstrapNodesV5 = []*enode.Node{}
	for _, url := range urls {
		if url != "" {
			node, err := enode.Parse(enode.ValidSchemes, url)
			if err != nil {
				log.Error("Bootstrap URL invalid", "enode", url, "err", err)
				continue
			}
			cfg.P2P.BootstrapNodesV5 = append(cfg.P2P.BootstrapNodesV5, node)
		}
	}
	cfg.P2P.BootstrapNodes = cfg.P2P.BootstrapNodesV5
}

func setDataDir(ctx *cli.Context, cfg *node.Config) {
	defaultDataDir := DefaultDataDir()

	switch {
	case ctx.GlobalIsSet(utils.DataDirFlag.Name):
		cfg.DataDir = ctx.GlobalString(utils.DataDirFlag.Name)
	case ctx.GlobalIsSet(FakeNetFlag.Name):
		_, num, err := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		if err != nil {
			log.Crit("Invalid flag", "flag", FakeNetFlag.Name, "err", err)
		}
		cfg.DataDir = filepath.Join(defaultDataDir, fmt.Sprintf("fakenet-%d", num))
	default:
		cfg.DataDir = defaultDataDir
	}
}

func setGPO(ctx *cli.Context, cfg *gasprice.Config) {
	if ctx.GlobalIsSet(utils.GpoMaxGasPriceFlag.Name) {
		cfg.MaxPrice = big.NewInt(ctx.GlobalInt64(utils.GpoMaxGasPriceFlag.Name))
	}
}

func setTxPool(ctx *cli.Context, cfg *evmcore.TxPoolConfig) {
	if ctx.GlobalIsSet(utils.TxPoolLocalsFlag.Name) {
		locals := strings.Split(ctx.GlobalString(utils.TxPoolLocalsFlag.Name), ",")
		for _, account := range locals {
			if trimmed := strings.TrimSpace(account); !common.IsHexAddress(trimmed) {
				utils.Fatalf("Invalid account in --txpool.locals: %s", trimmed)
			} else {
				cfg.Locals = append(cfg.Locals, common.HexToAddress(account))
			}
		}
	}
	if ctx.GlobalIsSet(utils.TxPoolNoLocalsFlag.Name) {
		cfg.NoLocals = ctx.GlobalBool(utils.TxPoolNoLocalsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolJournalFlag.Name) {
		cfg.Journal = ctx.GlobalString(utils.TxPoolJournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolRejournalFlag.Name) {
		cfg.Rejournal = ctx.GlobalDuration(utils.TxPoolRejournalFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceLimitFlag.Name) {
		cfg.PriceLimit = ctx.GlobalUint64(utils.TxPoolPriceLimitFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolPriceBumpFlag.Name) {
		cfg.PriceBump = ctx.GlobalUint64(utils.TxPoolPriceBumpFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountSlotsFlag.Name) {
		cfg.AccountSlots = ctx.GlobalUint64(utils.TxPoolAccountSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalSlotsFlag.Name) {
		cfg.GlobalSlots = ctx.GlobalUint64(utils.TxPoolGlobalSlotsFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolAccountQueueFlag.Name) {
		cfg.AccountQueue = ctx.GlobalUint64(utils.TxPoolAccountQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolGlobalQueueFlag.Name) {
		cfg.GlobalQueue = ctx.GlobalUint64(utils.TxPoolGlobalQueueFlag.Name)
	}
	if ctx.GlobalIsSet(utils.TxPoolLifetimeFlag.Name) {
		cfg.Lifetime = ctx.GlobalDuration(utils.TxPoolLifetimeFlag.Name)
	}
}

func gossipConfigWithFlags(ctx *cli.Context, src gossip.Config) (gossip.Config, error) {
	cfg := src

	setGPO(ctx, &cfg.GPO)
	setTxPool(ctx, &cfg.TxPool)

	if ctx.GlobalIsSet(RPCGlobalGasCapFlag.Name) {
		cfg.RPCGasCap = ctx.GlobalUint64(RPCGlobalGasCapFlag.Name)
	}
	if ctx.GlobalIsSet(RPCGlobalTxFeeCapFlag.Name) {
		cfg.RPCTxFeeCap = ctx.GlobalFloat64(RPCGlobalTxFeeCapFlag.Name)
	}

	err := setValidator(ctx, &cfg.Emitter)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func gossipStoreConfigWithFlags(ctx *cli.Context, src gossip.StoreConfig) (gossip.StoreConfig, error) {
	cfg := src
	if !ctx.GlobalBool(utils.SnapshotFlag.Name) {
		cfg.EVM.EnableSnapshots = false
	}
	return cfg, nil
}

func nodeConfigWithFlags(ctx *cli.Context, cfg node.Config) node.Config {
	utils.SetNodeConfig(ctx, &cfg)

	if !ctx.GlobalIsSet(FakeNetFlag.Name) {
		setBootnodes(ctx, Bootnodes, &cfg)
	}
	setDataDir(ctx, &cfg)
	return cfg
}

func cacheScaler(ctx *cli.Context) cachescale.Func {
	if !ctx.GlobalIsSet(CacheFlag.Name) {
		return cachescale.Identity
	}
	totalCache := ctx.GlobalInt(CacheFlag.Name)
	if totalCache < DefaultCacheSize {
		log.Crit("Invalid flag", "flag", CacheFlag.Name, "err", fmt.Sprintf("minimum cache size is %d MB", DefaultCacheSize))
	}
	return cachescale.Ratio{
		Base:   DefaultCacheSize - ConstantCacheSize,
		Target: uint64(totalCache - ConstantCacheSize),
	}
}

func mayMakeAllConfigs(ctx *cli.Context) (*config, error) {
	// Defaults (low priority)
	cacheRatio := cacheScaler(ctx)
	cfg := config{
		Node:          defaultNodeConfig(),
		Galaxy:        gossip.DefaultConfig(cacheRatio),
		GalaxyStore:   gossip.DefaultStoreConfig(cacheRatio),
		Lachesis:      abft.DefaultConfig(),
		LachesisStore: abft.DefaultStoreConfig(cacheRatio),
		VectorClock:   vecmt.DefaultConfig(cacheRatio),
	}

	if ctx.GlobalIsSet(FakeNetFlag.Name) {
		_, num, _ := parseFakeGen(ctx.GlobalString(FakeNetFlag.Name))
		cfg.Galaxy = gossip.FakeConfig(num, cacheRatio)
	}

	// Load config file (medium priority)
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadAllConfigs(file, &cfg); err != nil {
			return &cfg, err
		}
	}

	// Apply flags (high priority)
	var err error
	cfg.Galaxy, err = gossipConfigWithFlags(ctx, cfg.Galaxy)
	if err != nil {
		return nil, err
	}
	cfg.GalaxyStore, err = gossipStoreConfigWithFlags(ctx, cfg.GalaxyStore)
	if err != nil {
		return nil, err
	}
	cfg.Node = nodeConfigWithFlags(ctx, cfg.Node)
	if cfg.Galaxy.Emitter.Validator.ID != 0 && len(cfg.Galaxy.Emitter.PrevEmittedEventFile.Path) == 0 {
		cfg.Galaxy.Emitter.PrevEmittedEventFile.Path = cfg.Node.ResolvePath(path.Join("emitter", fmt.Sprintf("last-%d", cfg.Galaxy.Emitter.Validator.ID)))
	}

	if err := cfg.Galaxy.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func makeAllConfigs(ctx *cli.Context) *config {
	cfg, err := mayMakeAllConfigs(ctx)
	if err != nil {
		utils.Fatalf("%v", err)
	}
	return cfg
}

func defaultNodeConfig() node.Config {
	cfg := NodeDefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "icicb", "dag", "sfc", "abft", "web3")
	cfg.WSModules = append(cfg.WSModules, "eth", "icicb", "dag", "sfc", "abft", "web3")
	cfg.IPCPath = "galaxy.ipc"
	cfg.DataDir = DefaultDataDir()
	return cfg
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	cfg := makeAllConfigs(ctx)
	comment := ""

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}

	dump := os.Stdout
	if ctx.NArg() > 0 {
		dump, err = os.OpenFile(ctx.Args().Get(0), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer dump.Close()
	}
	dump.WriteString(comment)
	dump.Write(out)

	return nil
}

func checkConfig(ctx *cli.Context) error {
	_, err := mayMakeAllConfigs(ctx)
	return err
}
