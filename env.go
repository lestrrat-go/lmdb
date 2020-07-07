package lmdb

import (
	"github.com/lestrrat-go/lmdb/clib"
	"github.com/pkg/errors"
)

// Create calls mdb_env_create and allocates memory for the MDB_env structure.
// You must call `Close()` on it when you are done to release the resources.
// This method is not thread-safe: callers must protect its call
func (e *Env) Create() error {
	// XXX check if e.ptr is already populated? 

	if err := clib.EnvCreate(&e.ptr); err != nil {
		return errors.Wrap(err, `failed to create environment`)
	}
	return nil
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
