package main

import (
	"log"
	"myapp/funcs"
	"myapp/initializers"
	"myapp/routes"
	"time"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func main() {
	initializers.SetupEnv()
	app := pocketbase.New()
	allRoutes := make(map[string]*core.Record)

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		registry := template.NewRegistry()

		routes.AdminRoutes(se, app, registry, &allRoutes)
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
