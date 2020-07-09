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

		dbi.Zerocopy(true)
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
