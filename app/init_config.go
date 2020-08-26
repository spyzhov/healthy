package app

import (
	"fmt"
	"io/ioutil"

	"github.com/spyzhov/healthy/config"
)

func (app *Application) setConfig() (err error) {
	var content []byte
	if app.Config.ConfigFile != "" {
		content, err = ioutil.ReadFile(app.Config.ConfigFile)
		if err != nil {
			return fmt.Errorf("can't read config file[%s]: %w", app.Config.ConfigFile, err)
		}
	} else if app.Config.ConfigYaml != "" {
		content = []byte(app.Config.ConfigYaml)
	} else {
		return fmt.Errorf("config was not specified, use env: CONFIG_FILE or CONFIG_YAML")
	}
	app.StepConfig, err = config.NewConfig(content)
	return err
}
