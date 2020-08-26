package main

import (
	app "github.com/spyzhov/healthy/app/web"
)

func main() {
	application, err := app.New()
	if err != nil {
		panic(err)
	}
	defer application.Close()
	application.Start()
}
