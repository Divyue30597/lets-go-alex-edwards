package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// func (app *application) routes() *http.ServeMux {
// Update the signature for the routes() method so that it returns a
// http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.HandlerFunc(http.MethodGet, "/", app.homeHandler)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetViewHandler)
	router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreateHandler)
	router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePostHandler)

	// mux.HandleFunc("/", app.homeHandler)
	// mux.HandleFunc("/snippet/view", app.snippetViewHandler)
	// mux.HandleFunc("/snippet/create", app.snippetCreateHandler)

	// Pass the servemux as the 'next' parameter to the secureHeaders middleware.
	// Because secureHeaders is just a function, and the function returns a
	// http.Handler we don't need to do anything else.
	// return secureHeaders(mux)

	// We want the requests to be logged first, so logRequest middleware is attached to the servemux
	// return app.recoverPanic(app.logRequest(secureHeaders(mux)))

	/*  using external package to do the same above But the real power lies in the fact that you
	can use it to create middleware chains that can be assigned to variables, appended to, and reused.
		myChain := alice.New(myMiddlewareOne, myMiddlewareTwo)
		myOtherChain := myChain.Append(myMiddleware3)
		return myOtherChain.Then(myHandler)
	*/
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}