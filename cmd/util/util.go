package util

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"

	customErr "github.com/eyko139/immich-notifier/internal/errors"
	"github.com/eyko139/immich-notifier/internal/models"
)

type Helper struct {
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
    AppVersion string
}

func New(templateCache map[string]*template.Template, errlog, infolog *log.Logger, appVersion string) *Helper {
	return &Helper{
		TemplateCache: templateCache,
        InfoLog: infolog,
        ErrorLog: errlog,
        AppVersion: appVersion,
	}
}

func (h *Helper) Render(w http.ResponseWriter, template string, data any) {
	if ts, ok := h.TemplateCache[template]; !ok {
		h.ServerError(w, customErr.NewTemplateError(errors.New("Could not fetch template from cache")))
	} else {

		// writing template to a buffer first catches runtime errors
		buf := new(bytes.Buffer)

		err := ts.ExecuteTemplate(buf, "base", data)
		if err != nil {
		    h.ServerError(w, customErr.NewTemplateError(err))
		}
		buf.WriteTo(w)
	}
}

func (h *Helper) ReturnHtml(w http.ResponseWriter, templateName string, data any) {
	cwd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	staticPath := filepath.Join(cwd, "ui/html/singles")

	ts, err := template.ParseFiles(staticPath + "/" + templateName)

	if err != nil {
        h.ServerError(w, err)
	}

	ts.ExecuteTemplate(w, "base", data)
}

func (h *Helper) NewTemplateData(albums []models.Album, email, name string, telegramAvailable bool, userId string) *models.TemplateData {
	return &models.TemplateData{
		Albums: albums,
		User:   models.UserContext{Email: email, Name: name, TelegramAvailable: telegramAvailable, Authenticated: true, ID: userId},
        AppVersion: h.AppVersion,
	}
}

func (h *Helper) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// set frame depth to 2, we don't want to see this line first on the stack trace
	// when error occurs
	h.ErrorLog.Output(2, trace)
	h.ErrorLog.Print(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
