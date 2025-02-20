package main

import (
	"fmt"
	"log"
	"myapp/funcs"
	"myapp/routes"
	"net/http"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func main() {
	app := pocketbase.New()
	allRoutes := make(map[string]*core.Record)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		registry := template.NewRegistry()

		se.Router.GET("/admin/", func(e *core.RequestEvent) error {
			success := e.Request.URL.Query().Get("success")
			errorParam := e.Request.URL.Query().Get("error")

			return funcs.RenderAdminPage(app, registry, e, funcs.DetermineSuccessMessage(success), funcs.DetermineErrorMessage(errorParam))
		})

		se.Router.POST("/admin/", func(e *core.RequestEvent) error {
			success := e.Request.URL.Query().Get("success")
			errorParam := e.Request.URL.Query().Get("error")

			fmt.Println(funcs.DetermineSuccessMessage(success))
			return funcs.RenderAdminPage(app, registry, e, funcs.DetermineSuccessMessage(success), funcs.DetermineErrorMessage(errorParam))
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
				return e.Redirect(http.StatusPermanentRedirect, "/admin?error=UPDATE_PATH_INVALID_PARAMS")
			}

			record, err := app.FindRecordById("routes", id)
			if err != nil {
				fmt.Println("s", err)
				return e.Redirect(http.StatusPermanentRedirect, "/admin?error=UPDATE_PATH_NONEXISTANT")
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
				return e.Redirect(http.StatusPermanentRedirect, "/admin?error=UPDATE_PATH")
			}

			go func() {
				funcs.GetAllRoutes(app, &allRoutes)
			}()
			return e.Redirect(http.StatusPermanentRedirect, "/admin?success=UPDATED_PATH")
		})

		se.Router.POST("/admin/deletePath/{id}", func(e *core.RequestEvent) error {
			id := e.Request.PathValue("id")

			fmt.Println(id)

			go func() {
				funcs.GetAllRoutes(app, &allRoutes)
			}()
			return e.Redirect(http.StatusPermanentRedirect, "/admin?success=DELETED_PATH")
		})

		routes.BaseRoute(se, registry, &allRoutes)

		return se.Next()
	})

	go func() {
		time.Sleep(5 * time.Second)
		funcs.GetAllRoutes(app, &allRoutes)
	}()

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
