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
	// region Templates functions
	functions := template.FuncMap{
		"attr": func(s string) template.HTMLAttr {
			// #nosec G203
			return template.HTMLAttr(s)
		},
		"html": func(s string) template.HTML {
			// #nosec G203
			return template.HTML(s)
		},
		"style": func(s string) template.CSS {
			// #nosec G203
			return template.CSS(s)
		},
		"script": func(s string) template.JS {
			// #nosec G203
			return template.JS(s)
		},
	}
	// endregion
	// region Init templates
	app.templates["index"], err = template.New("index").
		Funcs(functions).
		Parse(safe.Must(box.FindString("index.gohtml")).(string))
	if err != nil {
		return err
	}
	// endregion

	return nil
}
