package app

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/database"
)

type App struct{}

func Start() int {
	database := database.NewConnection()

	queryResult := database.Query()

  fmt.Println(queryResult)
	fmt.Println(database)

  return 0
}
