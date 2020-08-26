package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/spyzhov/safe"
	"go.uber.org/zap"
)

// Start HTTP handler
func (app *Application) RunHttp(handler *http.ServeMux, port int, name string) error {
	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()
		app.Logger.Info(fmt.Sprintf("http handler '%s' started on [::]:%d", name, port))
		server := &http.Server{
			Addr:    ":" + strconv.Itoa(port),
			Handler: handler,
		}
		server.RegisterOnShutdown(app.CtxCancel)
		defer safe.Close(server, "http handler '"+name+"' close error")

		app.WaitGroup.Add(1)
		go func() {
			defer app.WaitGroup.Done()
			app.Error(server.ListenAndServe())
			app.Logger.Debug("http handler stops serve", zap.String("name", name))
		}()

		<-app.Ctx.Done()

		app.Logger.Debug("http stops", zap.String("name", name))
	}()
	return nil
}
