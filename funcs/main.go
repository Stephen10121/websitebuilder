package funcs

import (
	"net/http"
	"os"

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
