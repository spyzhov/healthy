package main

import (
	"github.com/spyzhov/healthy/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		panic(err)
	}
	defer application.Close()
	application.Start()
}
