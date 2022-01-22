# ICICB galaxy 

EVM-compatible chain secured by the Lachesis consensus algorithm.

## Building the source

Building `galaxy` requires both a Go (version 1.14 or later) and a C compiler. You can install
them using your favourite package manager. Once the dependencies are installed, run

```shell
make galaxy
```
The build output is ```build/galaxy``` executable.

## Running `galaxy`

Going through all the possible command line flags is out of scope here,
but we've enumerated a few common parameter combos to get you up to speed quickly
on how you can run your own `galaxy` instance.

### Launching a network

Launching `galaxy` for a network:

```shell
$ ./build/galaxy --genesis /path/to/genesis.g
```

### Configuration

As an alternative to passing the numerous flags to the `galaxy` binary, you can also pass a
configuration file via:

```shell
$ ./build/galaxy --config /path/to/your_config_file.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to
export your existing configuration:

```shell
$ ./build/galaxy --your-favourite-flags dumpconfig
```

#### Validator

New validator private key may be created with `galaxy validator new` command.
```shell
$ ./build/galaxy validator new
```


To launch a validator, you have to use `--validator.id` and `--validator.pubkey` flags to enable events emitter.

```shell
$ ./build/galaxy --nousb --validator.id YOUR_ID --validator.pubkey 0xYOUR_PUBKEY
```

`galaxy` will prompt you for a password to decrypt your validator private key. Optionally, you can
specify password with a file using `--validator.password` flag.

#### Participation in discovery

Optionally you can specify your public IP to straighten connectivity of the network.
Ensure your TCP/UDP p2p port (5050 by default) isn't blocked by your firewall.

```shell
$ ./build/galaxy --nat extip:1.2.3.4
```

## Dev

### Running testnet

The network is specified only by its genesis file, so running a testnet node is equivalent to
using a testnet genesis file instead of a mainnet genesis file:
```shell
$ ./build/galaxy --genesis /path/to/testnet.g # launch node
```

It may be convenient to use a separate datadir for your testnet node to avoid collisions with other networks:
```shell
$ ./build/galaxy --genesis /path/to/testnet.g --datadir /path/to/datadir # launch node
$ ./build/galaxy --datadir /path/to/datadir account new # create new account
$ ./build/galaxy --datadir /path/to/datadir attach # attach to IPC
```

### Testing

Lachesis has extensive unit-testing. Use the Go tool to run tests:
```shell
go test ./...
```

If everything goes well, it should output something along these lines:
```
ok  	github.com/galaxy-foundation/go-galaxy/app	0.033s
?   	github.com/galaxy-foundation/go-galaxy/cmd/cmdtest	[no test files]
ok  	github.com/galaxy-foundation/go-galaxy/cmd/galaxy	13.890s
?   	github.com/galaxy-foundation/go-galaxy/cmd/galaxy/metrics	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/cmd/galaxy/tracing	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/crypto	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/debug	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/ethapi	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/eventcheck	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/eventcheck/basiccheck	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/eventcheck/gaspowercheck	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/eventcheck/heavycheck	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/eventcheck/parentscheck	[no test files]
ok  	github.com/galaxy-foundation/go-galaxy/evmcore	6.322s
?   	github.com/galaxy-foundation/go-galaxy/gossip	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/gossip/emitter	[no test files]
ok  	github.com/galaxy-foundation/go-galaxy/gossip/filters	1.250s
?   	github.com/galaxy-foundation/go-galaxy/gossip/gasprice	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/gossip/occuredtxs	[no test files]
?   	github.com/galaxy-foundation/go-galaxy/gossip/piecefunc	[no test files]
ok  	github.com/galaxy-foundation/go-galaxy/integration	21.640s
```

Also it is tested with [fuzzing](./FUZZING.md).


### operating a private network (fakenet)

Fakenet is a private network optimized for your private testing.
It'll generate a genesis containing N validators with equal stakes.
To launch a validator in this network, all you need to do is specify a validator ID you're willing to launch.

Pay attention that validator's private keys are deterministically generated in this network, so you must use it only for private testing.

Maintaining your own private network is more involved as a lot of configurations taken for
granted in the official networks need to be manually set up.

To run the fakenet with just one validator (which will work practically as a PoA blockchain), use:
```shell
$ ./build/galaxy --fakenet 1/1
```

To run the fakenet with 5 validators, run the command for each validator:
```shell
$ ./build/galaxy --fakenet 1/5 # first node, use 2/5 for second node
```

If you have to launch a non-validator node in fakenet, use 0 as ID:
```shell
$ ./build/galaxy --fakenet 0/5
```

After that, you have to connect your nodes. Either connect them statically or specify a bootnode:
```shell
$ ./build/galaxy --fakenet 1/5 --bootnodes "enode://galaxy@1.2.3.4:5050"
```

### Running the demo

For the testing purposes, the full demo may be launched using:
```shell
cd demo/
./start.sh # start the galaxy processes
./stop.sh # stop the demo
./clean.sh # erase the chain data
```
