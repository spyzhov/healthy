package app

import (
	"fmt"
	"os"

	"github.com/spyzhov/safe"
)

func (app *Application) printVersion() {
	safe.Must(fmt.Printf("healthy %s -- %s\nBuild at: %s\n", app.Info.Version, app.Info.Commit, app.Info.Created))
	os.Exit(0)
}
