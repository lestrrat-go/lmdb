package lmdb

type DBI struct {
	handle uint
	txn *Txn
}

type Env struct {
	ptr uintptr
}

// TxnBeginner is used to abstract between Env / Txn when
// we want to create a new transaction.
type TxnBeginner interface {
	Begin(uint) (*Txn, error)
}

type TxnBody interface {
	Run(*Txn) error
}

type TxnBodyFunc func(*Txn) error
