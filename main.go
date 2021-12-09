package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	monitor "github.com/butuzov/stackoverflow/internal/app"
)

func main() {
	app := monitor.New()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		cli.HandleExitCoder(err)
	}
}
