package app

import (
	"fmt"

	"github.com/alfonzm/lazydb/internal/ui"
)

func Start() int {
	if err := ui.Start(); err != nil {
		fmt.Println(err)
		return 1
	}

	return 0
}
