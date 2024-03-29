package app

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/db"
)

func Start() int {
	dbClient, err := db.NewDBClient("root:root@/finance")
	if err != nil {
		fmt.Println(err)
		return 1
	}

	tables, err := dbClient.GetTables();
	if err != nil {
		fmt.Println(err)
		return 1
	}

	fmt.Println(tables)

  rows, err := dbClient.GetRecords("accounts")
  if err != nil {
    fmt.Println(err)
    return 1
  }

  fmt.Println(rows)

	return 0
}
