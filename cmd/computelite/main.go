package main

import (
	"os"

	"github.com/adnant1/computelite/cmd/computelite/app"
)

// The computelite binary is responsible for creating the cluster, starts controllers
// and begins printing cluster state continuously.
func main() {
	code := app.Run()
	os.Exit(code)
}	