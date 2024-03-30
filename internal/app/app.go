package app

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
	"github.com/alfonzm/lazydb/internal/ui"
)

func Start() int {
	db, err := db.NewDBClient("root:root@/finance")
	if err != nil {
		fmt.Println(err)
		return 1
	}

	if err := ui.Start(db); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}
