package lmdb

import (
	"runtime"

	"github.com/lestrrat-go/lmdb/internal/clib"
	"github.com/pkg/errors"
)

// TODO: these should be auto-generated
const (
	MDB_FIXEDMAP   = clib.MDB_FIXEDMAP
	MDB_NOSUBDIR   = clib.MDB_NOSUBDIR
	MDB_RDONLY     = clib.MDB_RDONLY
	MDB_WRITEMAP   = clib.MDB_WRITEMAP
	MDB_NOMETASYNC = clib.MDB_NOMETASYNC
	MDB_NOSYNC     = clib.MDB_NOSYNC
	MDB_MAPASYNC   = clib.MDB_MAPASYNC
	MDB_NOTLS      = clib.MDB_NOTLS
	MDB_NOLOCK     = clib.MDB_NOLOCK
	MDB_NORDAHEAD  = clib.MDB_NORDAHEAD
	MDB_NOMEMINIT  = clib.MDB_NOMEMINIT
)

func (fn TxnBodyFunc) Run(txn *Txn) error {
	return fn(txn)
}

// Run is a utility function that allows the user to execute a
// piece of code under a transaction.
//
// `parent` is the parent context in which the transaction is to be started.
// For root transactions, specify the *Env object. For nested transactions,
// specify the *Txn object to be used as the parent. If a nil parent is passed,
// an error is returned.
//
// If the transaction body returns an error, the transaction is always aborted.
//
// If the transaction
func Run(parent TxnBeginner, flags uint, body TxnBody) error {
	if parent == nil {
		return errors.New(`you must specify a parent transaction or the environment`)
	}

	rdonly := (flags & MDB_RDONLY) == 1

	if !rdonly {
		// transactions that requires writing, we need to lock the OS thread
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}

	txn, err := parent.Begin(flags)
	if err != nil {
		return errors.Wrap(err, `failed to begin a new transaction`)
	}

	if err := body.Run(txn); err != nil {
		// we're going to ignore this
		//nolint:errcheck
		txn.Abort()
		return errors.Wrap(err, `failed to run code under transaction`)
	}

	if rdonly {
		// we're going to ignore this
		//nolint:errcheck
		txn.Abort()
	} else {
		if err := txn.Commit(); err != nil {
			return errors.Wrap(err, `failed to commit transaction`)
		}
	}

	return nil
}
