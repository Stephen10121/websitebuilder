package routes

import (
	"myapp/funcs"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func BaseRoute(se *core.ServeEvent, app *pocketbase.PocketBase, registry *template.Registry) {
	se.Router.Any("/", func(e *core.RequestEvent) error {
		path := funcs.RemoveTrailingSlash(e.Request.URL.Path)

		record, err := app.FindFirstRecordByFilter(
			"routes",
			"path = {:path} && httpMethod = {:httpMethod}",
			dbx.Params{"path": path, "httpMethod": e.Request.Method},
		)

		if err != nil {
			return e.Error(http.StatusNotFound, "Page not found.", nil)
		} else {
			return funcs.ReturnCorrectResponse(record, e, registry)
		}
	})
}
