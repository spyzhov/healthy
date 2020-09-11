package main

import (
	"fmt"
	"os"

	"github.com/spyzhov/healthy/app"
	"github.com/spyzhov/safe"
)

//go:generate go-license -format md -output INHERITED_LICENSES.md

func main() {
	application, err := app.New()
	if err != nil {
		safe.Must(fmt.Fprintf(os.Stderr, "Error: %s\n", err))
		os.Exit(1)
	}
	defer application.Close()
	application.Start()
}
