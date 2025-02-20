package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"myapp/funcs"
	"myapp/initializers"
	"net/http"
	"os"
	"time"

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
	g := se.Router.Group("/admin")

	g.BindFunc(func(e *core.RequestEvent) error {
		auth, err := e.Request.Cookie("auth")
		if err != nil {
			return funcs.RenderLoginPage(registry, e)
		}

		mockData := e.Request.UserAgent() + initializers.AdminPassword + initializers.AdminSalt
		hash := sha256.New()
		hash.Write([]byte(mockData))
		bs := hash.Sum(nil)
		hexString := hex.EncodeToString(bs)

		if auth.Value != hexString {
			return funcs.RenderLoginPage(registry, e)
		}

		return e.Next()
	})

	se.Router.POST("/admin/login", func(e *core.RequestEvent) error {
		data := struct {
			Username string `json:"username" form:"username"`
			Password string `json:"password" form:"password"`
		}{}

		if err := e.BindBody(&data); err != nil {
			fmt.Println(err)
			return e.Redirect(http.StatusPermanentRedirect, "/admin?error=INVALID_LOGIN_PARAMETERS")
		}

		if data.Username == initializers.AdminUsername && data.Password == initializers.AdminPassword {
			data := e.Request.UserAgent() + data.Password + initializers.AdminSalt
			hash := sha256.New()
			hash.Write([]byte(data))
			bs := hash.Sum(nil)
			hexString := hex.EncodeToString(bs)

			cookie := new(http.Cookie)
			cookie.Name = "auth"
			cookie.Value = hexString
			cookie.Expires = time.Now().Add(24 * time.Hour)

			e.SetCookie(cookie)
			return e.Redirect(http.StatusPermanentRedirect, "/admin?success=LOGGED_IN")
		}

		return e.Redirect(http.StatusPermanentRedirect, "/admin?error=INVALID_LOGIN_PARAMETERS")
	})

	se.Router.GET("/admin/logout", func(e *core.RequestEvent) error {
		cookie := new(http.Cookie)
		cookie.Name = "auth"
		cookie.Value = ""
		cookie.Expires = time.Now().Add(0 * time.Hour)
		e.SetCookie(cookie)
		return e.Redirect(http.StatusPermanentRedirect, "/admin")
	})

	g.GET("/", func(e *core.RequestEvent) error {
		success := e.Request.URL.Query().Get("success")
		errorParam := e.Request.URL.Query().Get("error")

		return funcs.RenderAdminPage(app, registry, e, funcs.DetermineSuccessMessage(success), funcs.DetermineErrorMessage(errorParam))
	})

	g.POST("/", func(e *core.RequestEvent) error {
		success := e.Request.URL.Query().Get("success")
		errorParam := e.Request.URL.Query().Get("error")

		return funcs.RenderAdminPage(app, registry, e, funcs.DetermineSuccessMessage(success), funcs.DetermineErrorMessage(errorParam))
	})

	g.POST("/path/{id}", func(e *core.RequestEvent) error {
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

	g.POST("/deletePath/{id}", func(e *core.RequestEvent) error {
		id := e.Request.PathValue("id")

		fmt.Println(id)

		go func() {
			funcs.GetAllRoutes(app, allRoutes)
		}()
		return e.Redirect(http.StatusPermanentRedirect, "/admin?success=DELETED_PATH")
	})
}
