package lmdb

import (
	"github.com/lestrrat-go/lmdb/internal/clib"
)

func (dbi *DBI) Txn() *Txn {
	return dbi.txn
}

func (dbi *DBI) Get(key []byte) ([]byte, error) {
	var ret []byte
	if err := clib.GetBytes(dbi.Txn().ptr, dbi.handle, key, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (dbi *DBI) Put(key, value []byte, flags uint) error {
	return clib.PutBytes(dbi.Txn().ptr, dbi.handle, key, value, flags)
}
