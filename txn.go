package lmdb

import (
	"github.com/lestrrat-go/lmdb/clib"
	"github.com/pkg/errors"
)

type Txn struct {
	ptr uintptr
}

func (txn *Txn) requireValidTxn() error {
	if txn.ptr == 0 {
		return errors.New(`invalid transaction: txn has not bee initialized via Begin(), or has already been freed via Commit()/Abort()`)
	}
	return nil
}

func (txn *Txn) Begin(env *Env, parent *Txn, flags uint) error {
	var parentptr uintptr
	if parent != nil {
		parentptr = parent.ptr
	}
	return clib.TxnBegin(env.ptr, parentptr, flags, &txn.ptr)
}

func (txn *Txn) Abort() error {
	if err := txn.requireValidTxn(); err != nil {
		return errors.Wrap(err, `failed to abort`)
	}

	if err := clib.TxnAbort(txn.ptr); err != nil {
		return errors.Wrap(err, `failed to abort`)
	}
	return nil
}

func (txn *Txn) ID() int {
	if err := txn.requireValidTxn(); err != nil {
		return -1
	}

	return int(clib.TxnID(txn.ptr))
}
