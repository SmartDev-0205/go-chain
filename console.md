
### Make config file
```shell
./galaxy --genesis --rpc  -rpcaddr 0.0.0.0 --datadir=d:/icicb/galaxy --rpcapi "net,eth,txpool,web3" --maxpeers 128 --maxpendpeers 128 --txpool.globalqueue 4096  dumpconfig > d:/icicb/config.toml
```

### Running fakenet
```shell
./build/galaxy --fakenet 1/1 --http --http.addr="127.0.0.1" --http.port=5050 --http.corsdomain="*" --http.api="eth,debug,net,admin,web3,personal,txpool,icicb,dag" --datadir=d:/icicb/fake
```

### Running testnet
./build/galaxy --metrics  --cache 64000 --genesis
$HOME/icicb/mainnet.g --nousb --http --http.addr '0.0.0.0' --http.port 8545 --http.corsdomain "*" --http.vhosts "*" --ws --ws.addr '0.0.0.0' --ws.port 8546  --ws.origins '0.0.0.0' --graphql --graphql.corsdomain '*' --graphql.vhosts '*' --datadir "$HOME/icicb/node" --http.api "net,eth,web3" --ws.api "net,eth,web3"


### testing

generate genesis

```shell
./build/galaxy init d:/icicb/testgenesis.json --datadir d:/icicb/genesis.g
```

### generate validator key
```shell
./build/galaxy validator new
```

### run mainnet
```shell
./build/galaxy --genesis d:/icicb/genesis.g --datadir d:/icicb/mainnet --validator.id 0 --validator.pubkey 0xc0042ddbf4fda4bfefbb4368e3c0626faacfeb9319d08588cdb1a414eee7870ebc4c806aed6c5b01c2b2de5c60ba1cf98510dd4bcaecd56ad2dadc3c976279003797  --validator.password 123456789
```

```shell
./build/galaxy --genesis d:/icicb/genesis.g --datadir d:/icicb/mainnet --http --http.addr '0.0.0.0' --http.port 8545 --http.corsdomain "*" --http.vhosts "*" --http.api="eth,debug,net,admin,web3,personal,txpool,icicb,dag" --validator.id 1 --validator.pubkey 0xc004aa0713a0bc43226ab1f19e9547afc0d4b13fde317142436072d4d959621dc61a1d6ac185fb87ffe81d578c2a89a79e151a77ae71f31b8281e1c690c0f625dcf7 --validator.password D:\\icicb\\pass-1.txt

./build/galaxy --genesis d:/icicb/genesis.g --datadir d:/icicb/mainnet --nousb --validator.id 2 --validator.pubkey 0xc004398d3a0ee4514dfd111d22a4a54137b2ea79dc179f575f5b3dfe66ccc3348559a28af9ed7c075e3b93a5ff338dc657606104b3030ebcd95d7ad52fa62199553f --validator.password D:\\icicb\\pass-2.txt
```

`galaxy` will prompt you for a password to decrypt your validator private key. Optionally, you can
specify password with a file using `--validator.password` flag.

./build/galaxy --genesis d:/icicb/genesis.g --datadir d:/icicb/mainnet --bootnodes="enode://2cecf66045ee5f0defb2d0d88020a181504295882547b0442bd65246ab1a40e0164eb105558a22c4ce376bc01d25e3344521bc693915aaa5d19eb917c7acfc08@192.168.115.163:5060"
./build/galaxy init d:\\icicb\\genesis.json 

"init","d:\\icicb\\genesis.json"
"--genesis" 
"d:/icicb/genesis.g"
"--datadir"
"d:/icicb/mainnet"
"--http"
"--http.addr"
"'0.0.0.0'"
"--http.port"
"8545"
"--http.corsdomain"
"\"*\""
"--http.vhosts"
"\"*\""
"--http.api=\"eth,debug,net,admin,web3,personal,txpool,icicb,dag\""
"--nousb"
"--validator.id"
"1"
"--validator.pubkey"
"0xc004aa0713a0bc43226ab1f19e9547afc0d4b13fde317142436072d4d959621dc61a1d6ac185fb87ffe81d578c2a89a79e151a77ae71f31b8281e1c690c0f625dcf7"
"--validator.password"
"D:\\icicb\\pass-1.txt"
