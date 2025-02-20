package routes

import (
	"fmt"
	"myapp/funcs"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func BaseRoute(se *core.ServeEvent, registry *template.Registry, allRoutes *map[string]*core.Record) {
	se.Router.Any("/", func(e *core.RequestEvent) error {
		path := funcs.RemoveTrailingSlash(e.Request.URL.Path)

		val, ok := (*allRoutes)[path+e.Request.Method]
		if ok {
			return funcs.ReturnCorrectResponse(val, e, registry)
		} else {
			return e.FileFS(os.DirFS("./files/"), funcs.RemoveFirstSlash(path))
		}
	})
}

func AdminRoutes(se *core.ServeEvent, app *pocketbase.PocketBase, registry *template.Registry, allRoutes *map[string]*core.Record) {
	se.Router.GET("/admin/", func(e *core.RequestEvent) error {
		success := e.Request.URL.Query().Get("success")
		errorParam := e.Request.URL.Query().Get("error")

		return funcs.RenderAdminPage(app, registry, e, funcs.DetermineSuccessMessage(success), funcs.DetermineErrorMessage(errorParam))
	})

	se.Router.POST("/admin/", func(e *core.RequestEvent) error {
		success := e.Request.URL.Query().Get("success")
		errorParam := e.Request.URL.Query().Get("error")

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
			funcs.GetAllRoutes(app, allRoutes)
		}()
		return e.Redirect(http.StatusPermanentRedirect, "/admin?success=UPDATED_PATH")
	})

	se.Router.POST("/admin/deletePath/{id}", func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")

		fmt.Println(id)

		go func() {
			funcs.GetAllRoutes(app, allRoutes)
		}()
		return e.Redirect(http.StatusPermanentRedirect, "/admin?success=DELETED_PATH")
	})
}
