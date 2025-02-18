package main

import (
	"fmt"
	"log"
	"myapp/funcs"
	"myapp/routes"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		registry := template.NewRegistry()

		se.Router.GET("/admin/", func(e *core.RequestEvent) error {
			return funcs.RenderAdminPage(app, registry, e)
		})

		se.Router.POST("/admin/path/{id}", func(e *core.RequestEvent) error {
			id := e.Request.PathValue("id")
			data := struct {
				HttpMethod      string `json:"httpMethod" form:"httpMethod"`
				Serve           string `json:"serve" form:"serve"`
				JSONMessage     string `json:"jsonMessage" form:"jsonMessage"`
				StringMessage   string `json:"stringMessage" form:"stringMessage"`
				TemplateMessage string `json:"templateMessage" form:"templateMessage"`
				FileServePath   string `json:"fileServePath" form:"fileServePath"`
				Path            string `json:"path" form:"path"`
			}{}

			if err := e.BindBody(&data); err != nil {
				fmt.Println(err)
				return funcs.RenderAdminPage(app, registry, e)
			}

			record, err := app.FindRecordById("routes", id)
			if err != nil {
				fmt.Println("s", err)
				return funcs.RenderAdminPage(app, registry, e)
			}

			record.Set("httpMethod", data.HttpMethod)
			record.Set("serve", data.Serve)
			record.Set("jsonMessage", data.JSONMessage)
			record.Set("stringMessage", data.StringMessage)
			record.Set("templateMessage", data.TemplateMessage)
			record.Set("fileServePath", data.FileServePath)
			record.Set("path", data.Path)

			err = app.Save(record)
			if err != nil {
				fmt.Println("f", err)
			}

			return funcs.RenderAdminPage(app, registry, e)
		})

		routes.BaseRoute(se, app, registry)

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
