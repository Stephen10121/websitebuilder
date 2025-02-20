package routes

import (
	"myapp/funcs"
	"os"

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
