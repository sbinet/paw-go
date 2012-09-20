package pawmgr

import (
	"fmt"
)

type NtupleMgr struct {
}

func (mgr *NtupleMgr) Create(id int, title string, vars []string) error {
	fmt.Printf("==> /ntuple/create id=%v title=%q vars=%v\n",
		id, title, vars)
	return nil
}

// EOF
