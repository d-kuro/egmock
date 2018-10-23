package main

import (
	"os"

	"github.com/d-kuro/egmock/cli"
)

func main() {
	os.Exit(cli.Run(os.Args))
}
