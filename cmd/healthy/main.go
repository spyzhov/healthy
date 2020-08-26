package main

import (
	"fmt"
	"os"

	app "github.com/spyzhov/healthy/app/cmd"
	"github.com/spyzhov/safe"
)

func main() {
	application, err := app.New()
	if err != nil {
		safe.Must(fmt.Fprintf(os.Stderr, "Error: %s\n", err))
		os.Exit(1)
	}
	defer application.Close()
	application.Start()
}
