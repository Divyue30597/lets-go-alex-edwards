package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Divyue30597/snippetbox-lets-go/internal/models"
	"github.com/julienschmidt/httprouter"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	// http.NotFound(w, r)
	// 	app.notFound(w)
	// 	return
	// }

	// panic("oops! something went wrong")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// for _, snippet := range snippets {
	// 	fmt.Fprintf(w, "%+v\n", snippet)
	// }

	// files := []string{
	// 	"./ui/html/pages/home.tmpl",
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// }
	// // templateSet, err := template.ParseFiles("./ui/html/pages/home.tmpl")

	// templateSet, err := template.ParseFiles(files...)
	// if err != nil {
	// 	// app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// 	return
	// }

	// data := &templateData{
	// 	Snippets: snippets,
	// }

	// // err = templateSet.Execute(w, nil)
	// err = templateSet.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	// app.errorLog.Print(err.Error())
	// 	app.serverError(w, err)
	// 	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	// }

	// replacing the above logic with render method on app
	// app.render(w, http.StatusOK, "home.tmpl", &templateData{
	// 	Snippets: snippets,
	// })

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)

	// w.Write([]byte("Hello from snippetbox."))
}

func (app *application) snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	// id, err := strconv.Atoi(r.URL.Query().Get("id"))
	id, err := strconv.Atoi(params.ByName("id"))
	// r.URL.Query().Get("id")
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}

	// Using the snippets to get a record
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// files := []string{
	// 	"./ui/html/base.tmpl",
	// 	"./ui/html/partials/nav.tmpl",
	// 	"./ui/html/pages/view.tmpl",
	// }

	// templateSet, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// data := &templateData{
	// 	Snippet: snippet,
	// }

	// err = templateSet.ExecuteTemplate(w, "base", data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }

	// Changing the above operations to 2 different functions one for reading and another for rendering from the cache
	// app.render(w, http.StatusOK, "view.tmpl", &templateData{
	// 	Snippet: snippet,
	// })

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)

	// w.Write([]byte("Display a specific snippet..."))
	// fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
	// fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) snippetCreatePostHandler(w http.ResponseWriter, r *http.Request) {
	// Can be removed since we are using httprouter to choose the request methods.
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	// w.Header()["Date"] = nil
	// 	// w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
	// 	// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	return
	// }

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	expires := r.PostForm.Get("expires")

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// w.Write([]byte("Create a snippet..."))
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (app *application) snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	// w.Write([]byte("Display the form for creating a new snippet..."))
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.tmpl", data)
}
