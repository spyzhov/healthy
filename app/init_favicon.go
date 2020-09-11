package app

import (
	"github.com/gobuffalo/packr/v2"
)

// setFavicon initialize favicon var
func (app *Application) setFavicon() (err error) {
	box := packr.New("root", "../")
	app.favicon, err = box.Find("favicon.ico")
	if err != nil {
		return err
	}

	return nil
}
