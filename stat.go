package lmdb

import "github.com/lestrrat-go/lmdb/internal/clib"

func (stat *Stat) PSize() (uint, error) {
	return clib.StatPSize(stat.ptr)
}

func (stat *Stat) Depth() (uint, error) {
	return clib.StatDepth(stat.ptr)
}

func (stat *Stat) BranchPages() (uint, error) {
	return clib.StatBranchPages(stat.ptr)
}
