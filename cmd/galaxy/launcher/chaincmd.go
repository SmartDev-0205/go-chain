package launcher

import (
	"os"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/log"
	"gopkg.in/urfave/cli.v1"

	"github.com/goicicb/galaxy/genesisstore"
	"github.com/goicicb/integration/makegenesis"
)

var (
	EventsCheckFlag = cli.BoolTFlag{
		Name:  "check",
		Usage: "true if events should be fully checked before importing",
	}
	importCommand = cli.Command{
		Name:      "import",
		Usage:     "Import a blockchain file",
		ArgsUsage: "<filename> (<filename 2> ... <filename N>) [check=false]",
		Category:  "MISCELLANEOUS COMMANDS",
		Description: `
    galaxy import events

The import command imports events from an RLP-encoded files.
Events are fully verified by default, unless overridden by check=false flag.`,

		Subcommands: []cli.Command{
			{
				Action:    utils.MigrateFlags(importEvents),
				Name:      "events",
				Usage:     "Import blockchain events",
				ArgsUsage: "<filename> (<filename 2> ... <filename N>) [--check=false]",
				Flags: []cli.Flag{
					DataDirFlag,
					EventsCheckFlag,
				},
				Description: `
The import command imports events from RLP-encoded files.
Events are fully verified by default, unless overridden by --check=false flag.`,
			},
			{
				Action:    utils.MigrateFlags(importEvm),
				Name:      "evm",
				Usage:     "Import EVM storage",
				ArgsUsage: "<filename> (<filename 2> ... <filename N>)",
				Flags: []cli.Flag{
					DataDirFlag,
					EventsCheckFlag,
				},
				Description: `
    galaxy import evm

The import command imports EVM storage (trie nodes, code, preimages) from files.`,
			},
		},
	}
	exportCommand = cli.Command{
		Name:     "export",
		Usage:    "Export blockchain",
		Category: "MISCELLANEOUS COMMANDS",

		Subcommands: []cli.Command{
			{
				Name:      "events",
				Usage:     "Export blockchain events",
				ArgsUsage: "<filename> [<epochFrom> <epochTo>]",
				Action:    utils.MigrateFlags(exportEvents),
				Flags: []cli.Flag{
					DataDirFlag,
				},
				Description: `
    galaxy export events

Requires a first argument of the file to write to.
Optional second and third arguments control the first and
last epoch to write. If the file ends with .gz, the output will
be gzipped
`,
			},
		},
	}
	checkCommand = cli.Command{
		Name:     "check",
		Usage:    "Check blockchain",
		Category: "MISCELLANEOUS COMMANDS",

		Subcommands: []cli.Command{
			{
				Name:   "evm",
				Usage:  "Check EVM storage",
				Action: utils.MigrateFlags(checkEvm),
				Flags: []cli.Flag{
					DataDirFlag,
				},
				Description: `
    galaxy check evm

Checks EVM storage roots and code hashes
`,
			},
		},
	}

	initGenesisCommand = cli.Command{
		Action:    utils.MigrateFlags(initGenesis),
		Name:      "init",
		Usage:     "Checks configuration file",
		ArgsUsage: "",
		Flags: []cli.Flag{
			DataDirFlag,
		},
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The initGenesis make init genesis file .`,
	}

	initTestnetGenesisCommand = cli.Command{
		Action:    utils.MigrateFlags(initTestnetGenesis),
		Name:      "inittest",
		Usage:     "Checks configuration file",
		ArgsUsage: "",
		Flags: []cli.Flag{
			DataDirFlag,
		},
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The initTestnetGenesis make init genesis file for testnet.`,
	}
)

func initGenesis(ctx *cli.Context) error {
	/* var savePath string */
	genesisPath := ctx.Args().First()
	/* if len(genesisPath) == 0 {
		utils.Fatalf("Must supply path to genesis JSON file")
	}
	p1 := s.LastIndex(genesisPath, "/")
	p2 := s.LastIndex(genesisPath, ".")
	if p1 == -1 || p2 == -1 {
		utils.Fatalf("invalid genesis JSON file: %v", genesisPath)
	}
	savePath = genesisPath[0:p1] + genesisPath[p1:p2] + ".g"
	log.Info(savePath) */
	/* file, err := os.Open(genesisPath)
	if err != nil {
		utils.Fatalf("Failed to read genesis file: %v", err)
	}
	defer file.Close() */
	savefile, err := os.OpenFile(genesisPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		utils.Fatalf("Failed to read genesis file: %v", err)
	}
	defer savefile.Close()
	/* genesis := new(core.Genesis) */
	// if err := json.NewDecoder(file).Decode(genesis); err != nil {
	// 	utils.Fatalf("invalid genesis file: %v", err)
	// }
	/* var balance *big.Int = futils.ToIcicb(1000000) */
	/* keys := make([]common.Address, 0, len(genesis.Alloc))
	for a := range genesis.Alloc {
		keys = append(keys, a)
	} */

	store := makegenesis.MakeGenesisStore()
	h := store.Hash()
	log.Info(h.String())
	genesisstore.WriteGenesisStore(savefile, store)
	/* log.Info(balance.String(), savefile) */
	return nil
}

func initTestnetGenesis(ctx *cli.Context) error {
	genesisPath := ctx.Args().First()
	savefile, err := os.OpenFile(genesisPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		utils.Fatalf("Failed to read genesis file: %v", err)
	}
	defer savefile.Close()
	store := makegenesis.MakeTestnetGenesisStore()
	h := store.Hash()
	log.Info(h.String())
	genesisstore.WriteGenesisStore(savefile, store)
	return nil
}
