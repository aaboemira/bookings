package render

import (
	"bytes"
	"fmt"
	"github.com/aaboemira/bookings/internal/config"
	"github.com/aaboemira/bookings/internal/models"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}
var app *config.AppConfig

//New Templates sets the config for the templates

func NewTemplates(a *config.AppConfig) {
	app = a
}
func AddTempData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders a template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	// get the template cache from the app config
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = TemplateCreate()
	}
	temp, ok := tc[tmpl]
	if !ok {
		log.Fatal("couldn't get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = AddTempData(td, r)

	_ = temp.Execute(buf, td)

	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("error parsing template:", err)
	}
}
func TemplateCreate() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}
	pages, err := filepath.Glob("./templates/*page.tmpl")
	if err != nil {
		fmt.Println("1")
		return myCache, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println("2")
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*layout.tmpl")
		if err != nil {
			fmt.Println("3")
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				fmt.Println("4")
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
