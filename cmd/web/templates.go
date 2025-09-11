package main

import (
	"example/internal/models"
	"example/ui"
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
)

type templateData struct {
	CurrentYear     int
	Snippet         models.Snippet
	Snippets        []models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	User            models.User
}

func humanDate(t time.Time) string {

	if t.IsZero() {
		return ""
	}

	loc, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		return t.Format("02 Jan 2006 at 15:04")
	}
	return t.In(loc).Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	// Initialize a new map to act as the cache
	cache := map[string]*template.Template{}

	//pages, err := filepath.Glob("../../ui/html/pages/*.html")
	// To get a slice of all filepaths in the ui.Files embedded
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")

	if err != nil {
		return nil, err
	}

	// Loop through the page filepaths one by one
	for _, page := range pages {

		// Extract the file name from the full filepath
		name := filepath.Base(page)

		/*
			files := []string{
				"../../ui/html/base.html",
				"../../ui/html/partials/nav.html",
				page,
			}

			ts, err := template.New(name).Funcs(functions).ParseFiles("../../ui/html/base.html")
			if err != nil {
				return nil, err
			}

			ts, err = ts.ParseGlob("../../ui/html/partials/*.html")
			if err != nil {
				return nil, err
			}

			ts, err = ts.ParseFiles(page)
			if err != nil {
				return nil, err
			}
		*/

		// Create a slice containing the filepath patterns for the templates we want to parse
		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		// To parse the template files from the ui.Files embedded filesystem
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map using the name of the page as the key
		cache[name] = ts
	}

	// Return the map
	return cache, nil
}
