package app

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/database"
)

type App struct{}

func Start() {
	database := database.NewConnection()

	queryResult := database.Query()

  fmt.Println(queryResult)
	fmt.Println(database)
}
