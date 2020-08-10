package app

import (
	"html/template"

	"github.com/gobuffalo/packr/v2"
	"github.com/spyzhov/safe"
)

// setTemplates initialize all necessary templates
func (app *Application) setTemplates() (err error) {
	box := packr.New("templates", "../templates")
	app.templates = make(map[string]*template.Template)
	// region Init templates
	app.templates["index"], err = template.New("index").Parse(safe.Must(box.FindString("index.gohtml")).(string))
	if err != nil {
		return err
	}
	// endregion

	return nil
}
