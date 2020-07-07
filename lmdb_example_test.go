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
