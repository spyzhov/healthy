package app

import (
	"github.com/gobuffalo/packr/v2"
)

type publicFile struct {
	Content     []byte
	ContentType string
}

// setPublic initialize favicon variables
func (app *Application) setPublic() (err error) {
	box := packr.New("public", "../public")
	app.files = make(map[string]*publicFile, 0)

	app.files["/android-chrome-192x192.png"] = &publicFile{ContentType: "image/png"}
	app.files["/android-chrome-512x512.png"] = &publicFile{ContentType: "image/png"}
	app.files["/apple-touch-icon.png"] = &publicFile{ContentType: "image/png"}
	app.files["/favicon-16x16.png"] = &publicFile{ContentType: "image/png"}
	app.files["/favicon-32x32.png"] = &publicFile{ContentType: "image/png"}
	app.files["/favicon.ico"] = &publicFile{ContentType: "image/x-icon"}
	app.files["/site.webmanifest"] = &publicFile{ContentType: "application/json"}

	app.files["/android-chrome-192x192.png"].Content, err = box.Find("android-chrome-192x192.png")
	if err != nil {
		return err
	}
	app.files["/android-chrome-512x512.png"].Content, err = box.Find("android-chrome-512x512.png")
	if err != nil {
		return err
	}
	app.files["/apple-touch-icon.png"].Content, err = box.Find("apple-touch-icon.png")
	if err != nil {
		return err
	}
	app.files["/favicon-16x16.png"].Content, err = box.Find("favicon-16x16.png")
	if err != nil {
		return err
	}
	app.files["/favicon-32x32.png"].Content, err = box.Find("favicon-32x32.png")
	if err != nil {
		return err
	}
	app.files["/favicon.ico"].Content, err = box.Find("favicon.ico")
	if err != nil {
		return err
	}
	app.files["/site.webmanifest"].Content, err = box.Find("site.webmanifest")
	if err != nil {
		return err
	}

	return nil
}
