package main

import (
	"searcher/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	if err := app.InitProject(); err != nil {
		panic(err)
	}
}
