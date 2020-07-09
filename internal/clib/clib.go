/* Package clib is the only place where we interact witht he world of C directly.
 * There are two reasons for this.
 *
 * First, This allows us to minimize the locations where we need to import the "C" and
 * "unsafe" packages directly.
 *
 * Second, and most importanly, encpasulating the C world in a single package
 * makes it easier for consumers of this package -- the user-facing Go packages --
 * to build a more idomatic Go-ish interface. This is because C is inherently "flat":
 * there are no namespaces and modules, so it's far easier to deal with when the
 * code is laid out in the same way. In order to translate this to Go, it's much
 * easier to create single "flat" set of APIs in Go-land, and create a public,
 * namespaced API on top of it, as opposed to creating a namespaced API directly
 * on top of a flat C API.
 */
package clib

/*
// Note to self: lmdb does not provide a .pc file for pkg-config. So, um, yeah.
#cgo LDFLAGS: -llmdb
#include <lmdb.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/pkg/errors"
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

func newError(name string, val int) error {
	return Error{
		Message: name + " returned an error",
		Value:   val,
	}
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
	cstrpath := C.CString(path)
	defer C.free(unsafe.Pointer(cstrpath))
	if ret := C.mdb_env_open(env, cstrpath, C.uint(flags), C.mdb_mode_t(mode)); ret != 0 {
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

func TxnCommit(ptr uintptr) error {
	txn := (*C.MDB_txn)(unsafe.Pointer(ptr))
	if ret := C.mdb_txn_commit(txn); ret != 0 {
		return Error{Message: `mdb_txn_commit returned an error`, Value: int(ret)}
	}
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

func DbiOpen(txnptr uintptr, name string, flags uint, handle *uint) error {
	txn := (*C.MDB_txn)(unsafe.Pointer(txnptr))
	var dbi C.MDB_dbi

	var cstrname *C.char
	if name != "" {
		cstrname = C.CString(name)
		defer C.free(unsafe.Pointer(cstrname))
	}
	if ret := C.mdb_dbi_open(txn, cstrname, C.uint(flags), &dbi); ret != 0 {
		return Error{Message: `mdb_dbi_open returned an error`, Value: int(ret)}
	}
	*handle = uint(dbi)
	return nil
}

var nullTerminated = []byte{0}

func makeByteVal(b []byte) *C.MDB_val {
	l := len(b)
	if l == 0 {
		b = nullTerminated
	}

	return &C.MDB_val{
		mv_size: C.size_t(l),
		mv_data: unsafe.Pointer(&b[0]),
	}
}

func PutBytes(txnptr uintptr, dbi uint, key, val []byte, flags uint) error {
	var keyval, valval C.MDB_val
	keyval.mv_size = C.size_t(len(key))
	keyval.mv_data = unsafe.Pointer((*C.char)(unsafe.Pointer(&key[0])))
	valval.mv_size = C.size_t(len(val))
	valval.mv_data = unsafe.Pointer((*C.char)(unsafe.Pointer(&val[0])))
	return Put(txnptr, dbi, uintptr(unsafe.Pointer(&keyval)), uintptr(unsafe.Pointer(&valval)), flags)
}

func Put(txnptr uintptr, dbi uint, keyptr, valptr uintptr, flags uint) error {
	txn := (*C.MDB_txn)(unsafe.Pointer(txnptr))
	key := (*C.MDB_val)(unsafe.Pointer(keyptr))
	val := (*C.MDB_val)(unsafe.Pointer(valptr))
	if ret := C.mdb_put(txn, C.MDB_dbi(dbi), key, val, C.uint(flags)); ret != 0 {
		return newError(`mdb_put`, int(ret))
	}
	return nil
}

// very thin wrapper
func mdbGet(txnptr uintptr, dbi uint, keyptr, valptr uintptr) error {
	txn := (*C.MDB_txn)(unsafe.Pointer(txnptr))
	key := (*C.MDB_val)(unsafe.Pointer(keyptr))
	val := (*C.MDB_val)(unsafe.Pointer(valptr))
	if ret := C.mdb_get(txn, C.MDB_dbi(dbi), key, val); ret != 0 {
		return newError(`mdb_get`, int(ret))
	}
	return nil
}

func GetBytes(txnptr uintptr, zerocopy bool, dbi uint, key []byte, val *[]byte) error {
	var keyval, valval C.MDB_val
	keyval.mv_size = C.size_t(len(key))
	keyval.mv_data = unsafe.Pointer((*C.char)(unsafe.Pointer(&key[0])))

	if err := mdbGet(txnptr, dbi, uintptr(unsafe.Pointer(&keyval)), uintptr(unsafe.Pointer(&valval))); err != nil {
		return errors.Wrap(err, `failed to execute GetBytes`)
	}

	if zerocopy {
		*val = (*[0xffffffff]byte)(unsafe.Pointer(valval.mv_data))[:valval.mv_size:valval.mv_size]
	} else {
		*val = C.GoBytes(valval.mv_data, C.int(valval.mv_size))
	}
	return nil
}

func Get(txnptr uintptr, dbi uint, keyptr, valptr uintptr) error {
	return mdbGet(txnptr, dbi, keyptr, valptr)
}
