package lmdb

import (
	"github.com/lestrrat-go/lmdb/internal/clib"
	"github.com/pkg/errors"
)

// NewEnv calls mdb_env_create and allocates memory for the MDB_env structure.
// You must call `Close()` on it when you are done to release the resources.
// This method is not thread-safe: callers must protect its call
func NewEnv() (*Env, error) {
	var env Env
	if err := clib.EnvCreate(&env.ptr); err != nil {
		return nil, errors.Wrap(err, `failed to create environment`)
	}
	return &env, nil
}

func (e *Env) Close() error {
	if err := clib.EnvClose(e.ptr); err != nil {
		return errors.Wrap(err, `failed to close environment`)
	}
	return nil
}

func (e *Env) Open(path string, flags uint, mode uint) error {
	if err := clib.EnvOpen(e.ptr, path, flags, mode); err != nil {
		return errors.Wrap(err, `failed to open environment`)
	}
	return nil
}

func (e *Env) Begin(flags uint) (*Txn, error) {
	txn, err := NewTxn(e, nil, flags)
	if err != nil {
		return nil, errors.Wrap(err, `failed to create a new transaction`)
	}

	return txn, nil
}
