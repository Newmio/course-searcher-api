package main

import (
	"searcher/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	app.InitProject().Run()
}
