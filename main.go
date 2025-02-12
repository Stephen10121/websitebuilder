package main

import (
	"log"
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/", func(e *core.RequestEvent) error {
			record, err := app.FindFirstRecordByData("routes", "path", "/")
			if err != nil {
				return e.String(http.StatusNotFound, "Page not found.")
			} else {
				return e.String(http.StatusOK, record.GetString("serve"))
			}
		})

		se.Router.GET("/{name}/", func(e *core.RequestEvent) error {
			path := e.Request.URL.Path
			if path[len(path)-1:] == "/" {
				path = path[0 : len(path)-1]
			}

			record, err := app.FindFirstRecordByData("routes", "path", path)
			if err != nil {
				return e.String(http.StatusNotFound, "Page not found.")
			} else {
				return e.String(http.StatusOK, record.GetString("serve"))
			}
		})
		// serves static files from the provided public dir (if exists)
		// se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		return se.Next()
	})

	log.Println("Starting the server.")
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
