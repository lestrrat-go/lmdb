package lmdb

import (
	"github.com/lestrrat-go/lmdb/internal/clib"
	"github.com/pkg/errors"
)

type Txn struct {
	ptr uintptr
	env *Env
}

func (txn *Txn) requireValidTxn() error {
	if txn.ptr == 0 {
		return errors.New(`invalid transaction: txn has not bee initialized via Begin(), or has already been freed via Commit()/Abort()`)
	}
	return nil
}

func NewTxn(env *Env, parent *Txn, flags uint) (*Txn, error) {
	var txn Txn
	var parentptr uintptr
	if parent != nil {
		parentptr = parent.ptr
	}
	if err := clib.TxnBegin(env.ptr, parentptr, flags, &txn.ptr); err != nil {
		return nil, errors.Wrap(err, `failed to begin transaction`)
	}

	txn.env = env
	return &txn, nil
}

// Begin creates a new transaction, using the receiver object as the
// parent transaction. The environment is also automatically shared
// by the parent.
func (txn *Txn) Begin(flags uint) (*Txn, error) {
	newTxn, err := NewTxn(txn.env, txn, flags)
	if err != nil {
		return nil, errors.Wrap(err, `failed to create sub transaction`)
	}
	return newTxn, nil
}

func (txn *Txn) Commit() error {
	if err := txn.requireValidTxn(); err != nil {
		return errors.Wrap(err, `failed to commit`)
	}

	if err := clib.TxnCommit(txn.ptr); err != nil {
		return errors.Wrap(err, `failed to commit`)
	}
	return nil
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

// Open opens a database in the current transaction.
func (txn *Txn) Open(name string, flags uint) (*DBI, error) {
	var dbi DBI
	if err := clib.DbiOpen(txn.ptr, name, flags, &dbi.handle); err != nil {
		return nil, errors.Wrap(err, `failed to open database`)
	}

	dbi.txn = txn
	return &dbi, nil
}
