package main

import (
	"net/http"

	"github.com/Divyue30597/snippetbox-lets-go/ui"
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

	// fileServer := http.FileServer(http.Dir("./ui/static"))
	// router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))
	// mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// Using embed file systems
	// Take the ui.Files embedded filesystem and convert it to a http.FS type so
	// that it satisfies the http.FileSystem interface. We then pass that to the
	// http.FileServer() function to create the file server handler.
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// Adding a new GET /ping route to help
	router.HandlerFunc(http.MethodGet, "/ping", ping)

	// Use the nosurf middleware on all our 'dynamic' routes.
	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	// router.HandlerFunc(http.MethodGet, "/", app.homeHandler)
	// router.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetViewHandler)
	// router.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreateHandler)
	// router.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePostHandler)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.homeHandler))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetViewHandler))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignUpHandler))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignUpPostHandler))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLoginHandler))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPostHandler))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(app.snippetCreateHandler))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(app.snippetCreatePostHandler))
	router.Handler(http.MethodPost, "/user/logout", protected.ThenFunc(app.userLogoutPostHandler))

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
