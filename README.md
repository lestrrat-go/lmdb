# lmdb

Go Binding to LMDB

# DISCLAIMERS

Yes, I know there are other libraries that provide bindings to LMDB. I'm for the moment
just fooling around, seeing what LMDB's API is all about, and also thinking out loud
(writing actual code), seeing how _I_ would implement such a thing.

If you just want to casually contribute API for a Go/LMDB binding, feel free to send us a
PR. Who know, we *may* actually finish implementing it.

# CURRENT STATE OF THE PROJECT

* [x] Create/Close MDB_env
* [x] Create/Close MDB_txn
* [x] Create MDB_dbi
* [ ] Actually issue get/put
* [ ] Cursors
* [ ] Nice-to-have: elegantly integrate with context.Context
* [ ] Nice-to-have: semi-automatic resource cleanup
* [ ] Nice-to-have: detect-and-error or protect from accidental concurrent accesses

# SYNOPSIS

```go
package lmdb_test

import (
	"io/ioutil"
	"fmt"
	"os"

	"github.com/lestrrat-go/lmdb"
)

func ExampleEnv() {
	var env lmdb.Env
	if err := env.Create(); err != nil {
		fmt.Printf("failed to create the environment: %s\n", err)
		return
	}
	defer env.Close()

	dir, err := ioutil.TempDir("", "lmdb-example-")
	if err != nil {
		fmt.Printf("failed to create temporary directory: %s\n", err)
		return
	}
	defer os.RemoveAll(dir)

	if err := env.Open(dir, 0, 0644); err != nil {
		fmt.Printf("failed to open the environment: %s\n", err)
		return
	}

	// OUTPUT:
}

func ExampleTxn() {
	var env lmdb.Env
	if err := env.Create(); err != nil {
		fmt.Printf("failed to create the environment: %s\n", err)
		return
	}
	defer env.Close()

	dir, err := ioutil.TempDir("", "lmdb-example-")
	if err != nil {
		fmt.Printf("failed to create temporary directory: %s\n", err)
		return
	}
	defer os.RemoveAll(dir)

	if err := env.Open(dir, 0, 0644); err != nil {
		fmt.Printf("failed to open the environment: %s\n", err)
		return
	}

	var txn lmdb.Txn
	if err := txn.Begin(&env, nil, 0); err != nil {
		fmt.Printf("failed to begin transaction: %s\n", err)
		return
	}
	defer txn.Abort()

	fmt.Printf("txn id = %d\n", txn.ID())

	// OUTPUT:
	// txn id = 1
}
```
