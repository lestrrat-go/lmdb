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
* [x] Actually issue get/put
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
  env, err := lmdb.NewEnv()
  if err != nil {
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

  if err := lmdb.Run(env, lmdb.EmptyFlags, lmdb.TxnBodyFunc(func(txn *lmdb.Txn) error {
    dbi, err := txn.Open("", 0)
    if err != nil {
      fmt.Printf("failed to open database: %s\n", err)
      return err
    }

    if err := dbi.Put([]byte("Hello"), []byte("World"), 0); err != nil {
      fmt.Printf("failed to put value: %s\n", err)
      return err
    }

    val, err := dbi.Get([]byte("Hello"))
    if err != nil {
      fmt.Printf("failed to get value: %s\n", err)
      return err
    }
    fmt.Printf("val = %s\n", val)

    return nil
  })); err != nil {
    fmt.Printf("failed to execute lmdb.Run: %s\n", err)
    return
  }

  // OUTPUT:
  // val = World
}
```
