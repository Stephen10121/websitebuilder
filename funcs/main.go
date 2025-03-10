package funcs

import (
	"fmt"
	"math/rand"
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

func RenderAdminPage(app *pocketbase.PocketBase, registry *template.Registry, e *core.RequestEvent, successMsg string, errorMsg string) error {
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

	data := map[string]any{
		"records":    newRecords,
		"files":      FetchAllPublicFiles(),
		"success":    len(successMsg) > 0 && len(errorMsg) == 0,
		"successMsg": successMsg,
		"error":      len(errorMsg) > 0,
		"errorMsg":   errorMsg,
	}

	html, err := registry.LoadFiles("./admin/index.html").Render(data)

	if err != nil {
		fmt.Println(err)
		return e.Error(http.StatusNotFound, "Page not found.", nil)
	}

	return e.HTML(http.StatusOK, html)
}

func FetchAllPublicFiles() []string {
	filesInDir, err := os.ReadDir("./files/")
	if err != nil {
		fmt.Println(err)
		return []string{}
	}

	var files []string
	for _, file := range filesInDir {
		files = append(files, file.Name())
	}
	return files
}

func DetermineSuccessMessage(msg string) string {
	switch msg {
	case "DELETED_PATH":
		return "Successfully Deleted a path."
	case "UPDATED_PATH":
		return "Successfully Updated a path."
	case "LOGGED_IN":
		return "Successfully logged in!"
	default:
		return ""
	}
}

func DetermineErrorMessage(msg string) string {
	switch msg {
	case "UPDATE_PATH_NONEXISTANT":
		return "Path doesn't exist."
	case "UPDATE_PATH_INVALID_PARAMS":
		return "Invalid parameters."
	case "UPDATE_PATH":
		return "Failed to update the path."
	case "INVALID_LOGIN_PARAMETERS":
		return "Invalid Login parameters."
	case "RECORD_NONEXISTANT":
		return "The record doesn't exist!"
	case "RECORD_DELETE":
		return "Failed to delete record!"
	default:
		return ""
	}
}

func GetAllRoutes(app *pocketbase.PocketBase, allRoutes *map[string]*core.Record) {
	records, err := app.FindAllRecords("routes")
	*allRoutes = make(map[string]*core.Record)

	if err != nil {
		fmt.Println(err)
		return
	}

	for _, record := range records {
		(*allRoutes)[record.GetString("path")+record.GetString("httpMethod")] = record
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func RenderLoginPage(registry *template.Registry, e *core.RequestEvent) error {
	successMsg := DetermineSuccessMessage(e.Request.URL.Query().Get("success"))
	errorMsg := DetermineErrorMessage(e.Request.URL.Query().Get("error"))

	data := map[string]any{
		"success":    len(successMsg) > 0 && len(errorMsg) == 0,
		"successMsg": successMsg,
		"error":      len(errorMsg) > 0,
		"errorMsg":   errorMsg,
	}
	html, err := registry.LoadFiles("./admin/login.html").Render(data)

	if err != nil {
		fmt.Println(err)
		return e.Error(http.StatusInternalServerError, "Internal Server Error.", nil)
	}

	return e.HTML(http.StatusOK, html)
}
