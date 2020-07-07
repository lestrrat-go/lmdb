package clib

/*
// Note to self: lmdb does not provide a .pc file for pkg-config. So, um, yeah.
#cgo LDFLAGS: -llmdb
#include <lmdb.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Error is used to return a numeric return value from the
// C layer with its calling context
type Error struct {
	Message string
	Value   int
}

func (e Error) Error() string {
	v := C.mdb_strerror(C.int(e.Value))
	return fmt.Sprintf(`%s (error = %#v, raw = %d)`, e.Message, C.GoString(v), e.Value)
}

// TODO: these should be auto-generated
const (
	MDB_FIXEDMAP   = uint(C.MDB_FIXEDMAP)
	MDB_NOSUBDIR   = uint(C.MDB_NOSUBDIR)
	MDB_RDONLY     = uint(C.MDB_RDONLY)
	MDB_WRITEMAP   = uint(C.MDB_WRITEMAP)
	MDB_NOMETASYNC = uint(C.MDB_NOMETASYNC)
	MDB_NOSYNC     = uint(C.MDB_NOSYNC)
	MDB_MAPASYNC   = uint(C.MDB_MAPASYNC)
	MDB_NOTLS      = uint(C.MDB_NOTLS)
	MDB_NOLOCK     = uint(C.MDB_NOLOCK)
	MDB_NORDAHEAD  = uint(C.MDB_NORDAHEAD)
	MDB_NOMEMINIT  = uint(C.MDB_NOMEMINIT)
)

func EnvCreate(ptr *uintptr) error {
	var env *C.MDB_env
	if ret := C.mdb_env_create(&env); ret != 0 {
		return Error{Message: `mdb_env_create returned an error`, Value: int(ret)}
	}
	*ptr = uintptr(unsafe.Pointer(env))
	return nil
}

func EnvClose(ptr uintptr) error {
	env := (*C.MDB_env)(unsafe.Pointer(ptr))
	C.mdb_env_close(env)
	return nil
}

func EnvOpen(ptr uintptr, path string, flags uint, mode uint) error {
	env := (*C.MDB_env)(unsafe.Pointer(ptr))
	if ret := C.mdb_env_open(env, C.CString(path), C.uint(flags), C.mdb_mode_t(mode)); ret != 0 {
		return Error{Message: `mdb_env_open returned an error`, Value: int(ret)}
	}
	return nil
}

func TxnBegin(envptr uintptr, parentptr uintptr, flags uint, ptr *uintptr) error {
	env := (*C.MDB_env)(unsafe.Pointer(envptr))
	var parent *C.MDB_txn
	if parentptr != 0 {
		parent = (*C.MDB_txn)(unsafe.Pointer(envptr))
	}

	var txn *C.MDB_txn
	if ret := C.mdb_txn_begin(env, parent, C.uint(flags), &txn); ret != 0 {
		return Error{Message: `mdb_txn_begin returned an error`, Value: int(ret)}
	}
	*ptr = uintptr(unsafe.Pointer(txn))
	return nil
}

func TxnAbort(ptr uintptr) error {
	txn := (*C.MDB_txn)(unsafe.Pointer(ptr))
	C.mdb_txn_abort(txn)
	return nil
}

func TxnID(ptr uintptr) uint {
	txn := (*C.MDB_txn)(unsafe.Pointer(ptr))
	return uint(C.mdb_txn_id(txn))
}

func DbiOpen(txnptr uintptr, name string, flags uint, ptr *uintptr) error {
	txn := (*C.MDB_txn)(unsafe.Pointer(txnptr))
	var dbi *C.MDB_dbi
	if ret := C.mdb_dbi_open(txn, C.CString(name), C.uint(flags), dbi); ret != 0 {
		return Error{Message: `mdb_dbi_open returned an error`, Value: int(ret)}
	}
	*ptr = uintptr(unsafe.Pointer(dbi))
	return nil
}
