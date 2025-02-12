package main

import (
	"log"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.Any("/", func(e *core.RequestEvent) error {
			path := e.Request.URL.Path
			if path[len(path)-1:] == "/" && len(path) > 1 {
				path = path[0 : len(path)-1]
			}

			record, err := app.FindFirstRecordByFilter(
				"routes",
				"path = {:path} && httpMethod = {:httpMethod}",
				dbx.Params{"path": path, "httpMethod": e.Request.Method},
			)
			if err != nil {
				return e.String(http.StatusNotFound, "Page not found.")
			} else {
				return e.String(http.StatusOK, record.GetString("serve"))
			}
		})

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
