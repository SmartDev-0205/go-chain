package utils

import "math/big"

// ToIcicb number of ICICB to Wei
func ToIcicb(icicb uint64) *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(icicb), big.NewInt(1e18))
}
