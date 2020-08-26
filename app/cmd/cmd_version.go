package cmd

import (
	"fmt"
	"os"

	base "github.com/spyzhov/healthy/app"
	"github.com/spyzhov/safe"
)

func (app *Application) printVersion() {
	info := base.NewBuildInfo()
	safe.Must(fmt.Printf("healthy %s -- %s\nBuild at: %s\n", info.Version, info.Commit, info.Created))
	os.Exit(0)
}
