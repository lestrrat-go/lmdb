package lmdb_test

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/lestrrat-go/lmdb"
)

func ExampleEnv() {
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

	if err := lmdb.Run(env, lmdb.MDB_RDONLY, lmdb.TxnBodyFunc(func(txn *lmdb.Txn) error {
		fmt.Printf("txn id = %d\n", txn.ID())
		return nil
	})); err != nil {
		fmt.Printf("failed to execute lmdb.Run: %s\n", err)
		return
	}

	// OUTPUT:
	// txn id = 0
}
