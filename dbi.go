package lmdb

import (
	"github.com/lestrrat-go/lmdb/internal/clib"
	"github.com/pkg/errors"
)

type DBI struct {
	ptr uintptr
	txn *Txn
}

func (dbi *DBI) Open(txn *Txn, name string, flags uint) error {
	if err := clib.DbiOpen(txn.ptr, name, flags, &dbi.ptr); err != nil {
		return errors.Wrap(err, `failed to open database`)
	}

	dbi.txn = txn
	return nil
}
