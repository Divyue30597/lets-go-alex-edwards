package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/Divyue30597/snippetbox-lets-go/internal/models"
)

// Define a templateData type to act as the holding structure for
// any dynamic data that we want to pass to our HTML templates.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	// adding form field with type any
	Form  any
	Flash string
	// 1st, to pass this information to our HTML templates, so that we can toggle the contents of the navigation bar appropriately.
	IsAuthenticated bool
	CSRFToken       string
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// create an in-memory map with the type map[string]*template.Template to cache the parsed templates.
func newTemplateCache() (map[string]*template.Template, error) {
	templateCache := make(map[string]*template.Template)

	// get all the files using the filepath.Glob
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	if err != nil {
		return nil, err
	}

	// Loop through the pages
	for _, page := range pages {
		// Extracting the file name like 'home.tmpl' from full filepath.
		name := filepath.Base(page)

		templateSet, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.tmpl")
		if err != nil {
			return nil, err
		}

		templateSet, err = templateSet.ParseGlob("./ui/html/partials/*.tmpl")
		if err != nil {
			return nil, err
		}

		// files := []string{
		// 	"./ui/html/base.tmpl",
		// 	"./ui/html/partials/nav.tmpl",
		// 	page,
		// }

		// templateSet, err := template.ParseFiles(files...)
		// if err != nil {
		// 	return nil, err
		// }

		templateSet, err = templateSet.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add the template to the map, using the page name.
		templateCache[name] = templateSet
	}
	// return map
	return templateCache, nil
}
