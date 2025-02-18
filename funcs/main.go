package funcs

import (
	"fmt"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

/*
This function removes the trailing / from the string
*/
func RemoveTrailingSlash(path string) string {
	if path[len(path)-1:] == "/" && len(path) > 1 {
		return path[0 : len(path)-1]
	}
	return path
}

func RemoveFirstSlash(path string) string {
	if path[0:1] == "/" && len(path) > 1 {
		return path[1:]
	}
	return path
}

func getString(record *core.Record, e *core.RequestEvent) error {
	return e.String(http.StatusOK, record.GetString("stringMessage"))
}

func getFile(record *core.Record, e *core.RequestEvent) error {
	return e.FileFS(os.DirFS("./files"), record.GetString("fileServePath"))
}

func getJSON(record *core.Record, e *core.RequestEvent) error {
	var result map[string]interface{}
	err := record.UnmarshalJSONField("jsonMessage", &result)
	if err != nil {
		return e.Error(http.StatusNotFound, "Page not found.", nil)
	}
	return e.JSON(http.StatusOK, result)
}

type TemplateResType struct {
	TemplatePath string         `json:"templatePath"`
	Data         map[string]any `json:"data"`
}

func getTemplate(record *core.Record, e *core.RequestEvent, registry *template.Registry) error {
	var result TemplateResType
	err := record.UnmarshalJSONField("templateMessage", &result)
	if err != nil {
		return e.Error(http.StatusNotFound, "Page not found.", nil)
	}
	html, err := registry.LoadFiles("./files/" + result.TemplatePath).Render(result.Data)

	if err != nil {
		return e.Error(http.StatusNotFound, "Page not found.", nil)
	}

	return e.HTML(http.StatusOK, html)
}

func ReturnCorrectResponse(record *core.Record, e *core.RequestEvent, registry *template.Registry) error {
	serveType := record.GetString("serve")
	switch serveType {
	case "STRING":
		return getString(record, e)
	case "FILE":
		return getFile(record, e)
	case "JSON":
		return getJSON(record, e)
	case "TEMPLATE":
		return getTemplate(record, e, registry)
	default:
		return e.Error(http.StatusNotFound, "Page not found.", nil)
	}
}

func RenderAdminPage(app *pocketbase.PocketBase, registry *template.Registry, e *core.RequestEvent) error {
	records, _ := app.FindAllRecords("routes")
	newRecords := []any{}
	for _, record := range records {
		newRecord := map[string]any{
			"id":              record.Id,
			"path":            record.GetString("path"),
			"serve":           record.GetString("serve"),
			"jsonMessage":     record.GetString("jsonMessage"),
			"stringMessage":   record.GetString("stringMessage"),
			"templateMessage": record.GetString("templateMessage"),
			"fileServePath":   record.GetString("fileServePath"),
			"httpMethod":      record.GetString("httpMethod"),
		}
		newRecords = append(newRecords, newRecord)
	}

	data := map[string]any{"records": newRecords}

	html, err := registry.LoadFiles("./admin/index.html").Render(data)

	if err != nil {
		fmt.Println(err)
		return e.Error(http.StatusNotFound, "Page not found.", nil)
	}

	return e.HTML(http.StatusOK, html)
}
