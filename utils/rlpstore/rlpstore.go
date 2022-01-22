package rlpstore

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/galaxy-foundation/icicb-base/kvdb"
	"github.com/goicicb/logger"
)

type Helper struct {
	logger.Instance
}

// Set RLP value
func (s *Helper) Set(table kvdb.Store, key []byte, val interface{}) {
	buf, err := rlp.EncodeToBytes(val)
	if err != nil {
		s.Log.Crit("Failed to encode rlp", "err", err)
	}

	if err := table.Put(key, buf); err != nil {
		s.Log.Crit("Failed to put key-value", "err", err)
	}
}

// Get RLP value
func (s *Helper) Get(table kvdb.Store, key []byte, to interface{}) interface{} {
	buf, err := table.Get(key)
	if err != nil {
		s.Log.Crit("Failed to get key-value", "err", err)
	}
	if buf == nil {
		return nil
	}

	err = rlp.DecodeBytes(buf, to)
	if err != nil {
		s.Log.Crit("Failed to decode rlp", "err", err, "size", len(buf))
	}
	return to
}
