package makegenesis

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/galaxy-foundation/icicb-base/hash"
	"github.com/galaxy-foundation/icicb-base/inter/idx"

	galaxy "github.com/goicicb/galaxy"
	"github.com/goicicb/galaxy/genesis"
	"github.com/goicicb/galaxy/genesis/driver"
	"github.com/goicicb/galaxy/genesis/driverauth"
	"github.com/goicicb/galaxy/genesis/evmwriter"
	"github.com/goicicb/galaxy/genesis/gpos"
	"github.com/goicicb/galaxy/genesis/netinit"
	"github.com/goicicb/galaxy/genesis/sfc"
	"github.com/goicicb/galaxy/genesisstore"
	"github.com/goicicb/inter"
	"github.com/goicicb/inter/validatorpk"
	futils "github.com/goicicb/utils"
)

var (
	FakeGenesisTime = inter.Timestamp(1608600000 * time.Second)
)

// FakeKey gets n-th fake private key.
func FakeKey(n int) *ecdsa.PrivateKey {
	reader := rand.New(rand.NewSource(int64(n)))

	key, err := ecdsa.GenerateKey(crypto.S256(), reader)

	fmt.Printf("\nYour new privatekey was generated %x\n", key.D)

	if err != nil {
		panic(err)
	}

	return key
}

type ValidatorAccount struct {
	address   string
	validator string
}

func MakeGenesisStore() *genesisstore.Store {
	genStore := genesisstore.NewMemStore()
	genStore.SetRules(galaxy.MainNetRules())

	var validatorAccounts = []ValidatorAccount{
		// for mainnet
		{
			address:   "0x111118fa5f43538A4Fafc956920C575c29fCED43",
			validator: "04e96025f5c550b66b1c87b8ebeb0cbcb6b306b5b2d1b45186e04e9d1b8087752896b0a95983fd025a867f87848326ece5cfd75983d9b37706240438e30d5eecbb",
		},
		{
			address:   "0x2222a13de6d0ab4bfb017e1b29d6a5b430e07574",
			validator: "0484a101d22ac8b715683219bcea75c5b192821e6a9a5a6080eff983bfe0a4b8e5fbd1452300f9c9fae116987e030a80cd1fafd2e9bff2b2955af623dacf2fb6f2",
		},
		{
			address:   "0x3333bfdaf4a1b2f0b748f3a88e3c853d1e2b8f5d",
			validator: "047301efa2653128b50490d37bd7a58a1af88ba397b1f2a362487e1d3a02113d6450ce444ca39071867ce874172e8c465bf41bdf2f68ff5fc667eeae366df4377d",
		},
		{
			address:   "0x44441858f827b0ba5a72ede2f8615c6f3deaeab4",
			validator: "04c33e9f40aa4c4db2a117a48581b17d3e31acebf87ea54ce6d19d95da134ea0cf207654a5e2d312d965a4a56a7e24c39a11c69d9995b053b617047c4e1f53e55e",
		},
		{
			address:   "0x5555e6e7f5d013ed35e76c117a18ac196bd0eee1",
			validator: "046cf5490c7702a6678cfaaf6ad0f5099b38a0db4cce030b5f48844ccbe94d05299cedbf6d12e5909b508330d8132c17b8a066ada99235593364529faa86f39bba",
		},
		{
			address:   "0x66664259885aa5a378edc418587f85943d899eb6",
			validator: "045024d1fded08f26d87261794e3b4a7c75a4b56206894596df2ed67114bccf4f09647ba962f929543781c7de7ede4fdc119bbf04bcd815268ff8c440354def788",
		},
		{
			address:   "0x77773049d7754024711c4615f0714752acc8676b",
			validator: "04b8bd4d5deeb04a89788b9fcb937beecd0755ead73bb91d65c2a1ffc702fc5f8c02b6586a471389aeaf2ac76bcb910ee998bf978fb710221c225a7c4e15b4fd73",
		},
		{
			address:   "0x8888a8d58edeb333d35ef9d93112e79383890a8f",
			validator: "0481d3d50a85f48b791ced92002db6f61a0e7209445e31b3a23537332bec6b02e99ae691ab446ba20ef09e1c167cb8a07a31823eeb94828d5e34074af52f63719e",
		},
	}

	var initialAccounts = []string{
		"0xb8B5BE7122f317F86b47778422e277cD91C0B031",
	}
	num := len(validatorAccounts)

	_total := 5000
	_validator := 1
	_staker := 5
	_initial := (5000 - (_validator+_staker)*num) / len(initialAccounts)

	totalSupply := futils.ToIcicb(uint64(_total) * 1e6)
	balance := futils.ToIcicb(uint64(_validator) * 1e6)
	stake := futils.ToIcicb(uint64(_staker) * 1e6)
	initialBalance := futils.ToIcicb(uint64(_initial) * 1e6)

	validators := make(gpos.Validators, 0, num)

	now := time.Now() // current local time
	// sec := now.Unix()      // number of seconds since January 1, 1970 UTC
	nsec := now.UnixNano()
	time := inter.Timestamp(nsec)
	for i := 1; i <= num; i++ {
		addr := common.HexToAddress(validatorAccounts[i-1].address)
		pubkeyraw := common.Hex2Bytes(validatorAccounts[i-1].validator)
		// fmt.Printf("\n# addr %x pubkeyraw %s len %d\n", addr, hex.EncodeToString(pubkeyraw), len(pubkeyraw))
		validatorID := idx.ValidatorID(i)
		pubKey := validatorpk.PubKey{
			Raw:  pubkeyraw,
			Type: validatorpk.Types.Secp256k1,
		}

		validators = append(validators, gpos.Validator{
			ID:               validatorID,
			Address:          addr,
			PubKey:           pubKey,
			CreationTime:     time,
			CreationEpoch:    0,
			DeactivatedTime:  0,
			DeactivatedEpoch: 0,
			Status:           0,
		})
	}
	for _, val := range initialAccounts {
		genStore.SetEvmAccount(common.HexToAddress(val), genesis.Account{
			Code:    []byte{},
			Balance: initialBalance,
			Nonce:   0,
		})
	}
	for _, val := range validators {
		genStore.SetEvmAccount(val.Address, genesis.Account{
			Code:    []byte{},
			Balance: balance,
			Nonce:   0,
		})
		genStore.SetDelegation(val.Address, val.ID, genesis.Delegation{
			Stake:              stake,
			Rewards:            new(big.Int),
			LockedStake:        new(big.Int),
			LockupFromEpoch:    0,
			LockupEndTime:      0,
			LockupDuration:     0,
			EarlyUnlockPenalty: new(big.Int),
		})
	}

	var owner common.Address
	if num != 0 {
		owner = validators[0].Address
	}

	genStore.SetMetadata(genesisstore.Metadata{
		Validators:    validators,
		FirstEpoch:    2,
		Time:          time,
		PrevEpochTime: time - inter.Timestamp(time.Time().Hour()),
		ExtraData:     []byte("galaxy"),
		DriverOwner:   owner,
		TotalSupply:   totalSupply,
	})
	genStore.SetBlock(0, genesis.Block{
		Time:        time - inter.Timestamp(time.Time().Minute()),
		Atropos:     hash.Event{},
		Txs:         types.Transactions{},
		InternalTxs: types.Transactions{},
		Root:        hash.Hash{},
		Receipts:    []*types.ReceiptForStorage{},
	})
	// pre deploy NetworkInitializer
	genStore.SetEvmAccount(netinit.ContractAddress, genesis.Account{
		Code:    netinit.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy NodeDriver
	genStore.SetEvmAccount(driver.ContractAddress, genesis.Account{
		Code:    driver.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy NodeDriverAuth
	genStore.SetEvmAccount(driverauth.ContractAddress, genesis.Account{
		Code:    driverauth.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy SFC
	genStore.SetEvmAccount(sfc.ContractAddress, genesis.Account{
		Code:    sfc.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// set non-zero code for pre-compiled contracts
	genStore.SetEvmAccount(evmwriter.ContractAddress, genesis.Account{
		Code:    []byte{0},
		Balance: new(big.Int),
		Nonce:   0,
	})

	return genStore
}
func MakeTestnetGenesisStore() *genesisstore.Store {
	genStore := genesisstore.NewMemStore()
	genStore.SetRules(galaxy.TestNetRules())
	var validatorAccounts = []ValidatorAccount{
		{
			address:   "0x034c37E7850A0DB0298664f723b357ff25FF31E2",
			validator: "0439b9f3f5a56c6aa8c79e01094d496d4e5b0b2116f6e26790177fb7639ffdf473ed428b71eec45e9789e3210cd46e663b9852d2f58ce7070bf1c928ace37d904a",
		},
		{
			address:   "0xA3F3571734840d9A01279D2696bDD20342eBF302",
			validator: "04be3ddfc6d48ad5d0ab793968f30e412cfaf0a1e1bdf3af63f542c4082191349ef8d2d13a1a1ce2b1512526f1ff53bbb8365d7f2e953b64fcc8cc93ce6ab60d9d",
		},
	}
	/* var validatorAccounts = []ValidatorAccount{
		{
			address:   "0x034c37E7850A0DB0298664f723b357ff25FF31E2",
			validator: "0439b9f3f5a56c6aa8c79e01094d496d4e5b0b2116f6e26790177fb7639ffdf473ed428b71eec45e9789e3210cd46e663b9852d2f58ce7070bf1c928ace37d904a",
		},
		{
			address:   "0xA3F3571734840d9A01279D2696bDD20342eBF302",
			validator: "04be3ddfc6d48ad5d0ab793968f30e412cfaf0a1e1bdf3af63f542c4082191349ef8d2d13a1a1ce2b1512526f1ff53bbb8365d7f2e953b64fcc8cc93ce6ab60d9d",
		},
		{
			address:   "0x4d16A5DCA915C2f2CE039f4204548A94f520EF2e",
			validator: "0454c530e6781b0c7bb378199d903651745c68239c150c793751d4fbb2bf923eb7e3fd155b5235140e8534fdafd726f0dd213ab5bc07d51bcc598aa67bb260901a",
		},
		{
			address:   "0x7ECB5240FB7237bE35ddd5E6B08994A8FC43E52D",
			validator: "04847343604e986ba2b2fdd64c905503b85918fa206cac3df141e8bb61afee7c07e77196b4372d0b9ff9331b1eeff0a989ca22ab5c95cbfec4f77dc584a96f7278",
		},
		{
			address:   "0xdD7792225BD36410F9deFD98878890eb6c8135ad",
			validator: "04303bbac1433b0d46feae674db33fd2b62794803f53f0335d8743ae7f6005608bee87c7f73a25fe296fd184536bd83a79e60e6c35d53ca208d2aa00f0b18d36d2",
		},
		{
			address:   "0xaA492E71d793C99D824efe945aE2091eb6e41977",
			validator: "043ecdc02c855c322be643aee9b8f735e82bc664746b09304a9883e553cf64ba2e91c1f3ab39ae24ae27ed3d41c93d947b0e64cb78c06dd408532f99ef2d207895",
		},
		{
			address:   "0xa73B9365479fBB5008E5222078F639b6039c2Aed",
			validator: "0486fb9204c56ce4fc2e2de29db5c7df9917ea22cdd39b21234dd35a191d3e9677e199142799001101d79b5b8c4a2966c072b3a9f7c06c55151cd2ad27a3d3cd8b",
		},
		{
			address:   "0x0716C6Bb0573e3FD902Eb4A7311863f8a1E411b9",
			validator: "04be2adee5c7b3d15cdb8ae0a099a759117cc6d5bbe45018fe6e0b05d645ca43069051819a3173c7e87b8947d1d6d3ae85c9dbf725f6778ef417caf186bd8fcac9",
		},
	}
	*/
	var initialAccounts = []string{
		"0x9cD60D0D9e4404Be3a1C890cAF477A08903Aca2b",
		"0x8E1c7C2960B5298Bc2580619224E56023a27996B",
		"0xEe8E84116F1903c1F0d723E9d1a92D20613a50d2",
		"0x5f4632ceD4D32B02c9d2217B19888b8eC9749114",
		"0x294cD6A64d63e9cbd92358C74cE751c43DE9F3dC",
		"0x6752bDd135D92025611c01ab0f16b532a046E863",
		"0xFaca4DAe41dcDD2618FfD083cf03Ef4C05078B79",
		"0xd68ccE056fe53c6C349AdE5De472597B8D2b576c",
		"0x4Af5d38b634C36d29F28Fe948383fA3be9fccda2",
		"0x9509eb170B5007e5Ac607944F800b8A475cc9bC7",
	}

	num := len(validatorAccounts)

	_total := 5000
	_validator := 10
	_staker := 100
	_initial := (5000 - (_validator+_staker)*num) / 10

	totalSupply := futils.ToIcicb(uint64(_total) * 1e6)
	balance := futils.ToIcicb(uint64(_validator) * 1e6)
	stake := futils.ToIcicb(uint64(_staker) * 1e6)
	initialBalance := futils.ToIcicb(uint64(_initial) * 1e6)

	validators := make(gpos.Validators, 0, num)

	now := time.Now() // current local time
	// sec := now.Unix()      // number of seconds since January 1, 1970 UTC
	nsec := now.UnixNano()
	time := inter.Timestamp(nsec)
	for i := 1; i <= num; i++ {
		addr := common.HexToAddress(validatorAccounts[i-1].address)
		pubkeyraw := common.Hex2Bytes(validatorAccounts[i-1].validator)
		fmt.Printf("\n# addr %x pubkeyraw %s len %d\n", addr, hex.EncodeToString(pubkeyraw), len(pubkeyraw))
		validatorID := idx.ValidatorID(i)
		pubKey := validatorpk.PubKey{
			Raw:  pubkeyraw,
			Type: validatorpk.Types.Secp256k1,
		}

		validators = append(validators, gpos.Validator{
			ID:               validatorID,
			Address:          addr,
			PubKey:           pubKey,
			CreationTime:     time,
			CreationEpoch:    0,
			DeactivatedTime:  0,
			DeactivatedEpoch: 0,
			Status:           0,
		})
	}

	for _, val := range initialAccounts {
		genStore.SetEvmAccount(common.HexToAddress(val), genesis.Account{
			Code:    []byte{},
			Balance: initialBalance,
			Nonce:   0,
		})
	}

	for _, val := range validators {
		genStore.SetEvmAccount(val.Address, genesis.Account{
			Code:    []byte{},
			Balance: balance,
			Nonce:   0,
		})
		genStore.SetDelegation(val.Address, val.ID, genesis.Delegation{
			Stake:              stake,
			Rewards:            new(big.Int),
			LockedStake:        new(big.Int),
			LockupFromEpoch:    0,
			LockupEndTime:      0,
			LockupDuration:     0,
			EarlyUnlockPenalty: new(big.Int),
		})
	}

	var owner common.Address
	if num != 0 {
		owner = validators[0].Address
	}

	genStore.SetMetadata(genesisstore.Metadata{
		Validators:    validators,
		FirstEpoch:    2,
		Time:          time,
		PrevEpochTime: time - inter.Timestamp(time.Time().Hour()),
		ExtraData:     []byte("fake"),
		DriverOwner:   owner,
		TotalSupply:   totalSupply,
	})
	genStore.SetBlock(0, genesis.Block{
		Time:        time - inter.Timestamp(time.Time().Minute()),
		Atropos:     hash.Event{},
		Txs:         types.Transactions{},
		InternalTxs: types.Transactions{},
		Root:        hash.Hash{},
		Receipts:    []*types.ReceiptForStorage{},
	})
	// pre deploy NetworkInitializer
	genStore.SetEvmAccount(netinit.ContractAddress, genesis.Account{
		Code:    netinit.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy NodeDriver
	genStore.SetEvmAccount(driver.ContractAddress, genesis.Account{
		Code:    driver.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy NodeDriverAuth
	genStore.SetEvmAccount(driverauth.ContractAddress, genesis.Account{
		Code:    driverauth.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy SFC
	genStore.SetEvmAccount(sfc.ContractAddress, genesis.Account{
		Code:    sfc.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// set non-zero code for pre-compiled contracts
	genStore.SetEvmAccount(evmwriter.ContractAddress, genesis.Account{
		Code:    []byte{0},
		Balance: new(big.Int),
		Nonce:   0,
	})

	return genStore
}
func FakeGenesisStore(num int, balance, stake *big.Int) *genesisstore.Store {
	genStore := genesisstore.NewMemStore()
	genStore.SetRules(galaxy.FakeNetRules())

	validators := GetFakeValidators(num)

	totalSupply := new(big.Int)
	for _, val := range validators {
		genStore.SetEvmAccount(val.Address, genesis.Account{
			Code:    []byte{},
			Balance: balance,
			Nonce:   0,
		})
		genStore.SetDelegation(val.Address, val.ID, genesis.Delegation{
			Stake:              stake,
			Rewards:            new(big.Int),
			LockedStake:        new(big.Int),
			LockupFromEpoch:    0,
			LockupEndTime:      0,
			LockupDuration:     0,
			EarlyUnlockPenalty: new(big.Int),
		})
		totalSupply.Add(totalSupply, balance)
	}

	var owner common.Address
	if num != 0 {
		owner = validators[0].Address
	}

	genStore.SetMetadata(genesisstore.Metadata{
		Validators:    validators,
		FirstEpoch:    2,
		Time:          FakeGenesisTime,
		PrevEpochTime: FakeGenesisTime - inter.Timestamp(time.Hour),
		ExtraData:     []byte("fake"),
		DriverOwner:   owner,
		TotalSupply:   totalSupply,
	})
	genStore.SetBlock(0, genesis.Block{
		Time:        FakeGenesisTime - inter.Timestamp(time.Minute),
		Atropos:     hash.Event{},
		Txs:         types.Transactions{},
		InternalTxs: types.Transactions{},
		Root:        hash.Hash{},
		Receipts:    []*types.ReceiptForStorage{},
	})
	// pre deploy NetworkInitializer
	genStore.SetEvmAccount(netinit.ContractAddress, genesis.Account{
		Code:    netinit.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy NodeDriver
	genStore.SetEvmAccount(driver.ContractAddress, genesis.Account{
		Code:    driver.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy NodeDriverAuth
	genStore.SetEvmAccount(driverauth.ContractAddress, genesis.Account{
		Code:    driverauth.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// pre deploy SFC
	genStore.SetEvmAccount(sfc.ContractAddress, genesis.Account{
		Code:    sfc.GetContractBin(),
		Balance: new(big.Int),
		Nonce:   0,
	})
	// set non-zero code for pre-compiled contracts
	genStore.SetEvmAccount(evmwriter.ContractAddress, genesis.Account{
		Code:    []byte{0},
		Balance: new(big.Int),
		Nonce:   0,
	})

	return genStore
}

func GetFakeValidators(num int) gpos.Validators {
	validators := make(gpos.Validators, 0, num)

	for i := 1; i <= num; i++ {
		key := FakeKey(i)
		addr := crypto.PubkeyToAddress(key.PublicKey)
		pubkeyraw := crypto.FromECDSAPub(&key.PublicKey)

		validatorID := idx.ValidatorID(i)
		validators = append(validators, gpos.Validator{
			ID:      validatorID,
			Address: addr,
			PubKey: validatorpk.PubKey{
				Raw:  pubkeyraw,
				Type: validatorpk.Types.Secp256k1,
			},
			CreationTime:     FakeGenesisTime,
			CreationEpoch:    0,
			DeactivatedTime:  0,
			DeactivatedEpoch: 0,
			Status:           0,
		})
	}

	return validators
}

type Genesis struct {
	Nonce      uint64         `json:"nonce"`
	Timestamp  uint64         `json:"timestamp"`
	ExtraData  []byte         `json:"extraData"`
	GasLimit   uint64         `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int       `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash    `json:"mixHash"`
	Coinbase   common.Address `json:"coinbase"`
	Alloc      GenesisAlloc   `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
	BaseFee    *big.Int    `json:"baseFeePerGas"`
}

type GenesisAlloc map[common.Address]GenesisAccount

type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}
