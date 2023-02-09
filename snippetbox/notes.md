# Let's go Notes

## Chapter 2.2

### Handlers

Handlers -> think of it as controllers in MVC pattern. They are used to implement the logic and for writing HTTP response headers and bodies.

http.ResponseWriter -> provides methods for assembling a HTTP response and sending it to the user

http.Request -> a pointer to a struct which holds information about the current request (like the HTTP method and the URL being requested)

Servemux is a router in go terminology. This stores a mapping between the URL patterns for your application and the corresponding handlers.

## Chapter 2.3

### Fixed Path and subtree pattern

Servemux supports two types of URL pattern: _fixed paths_ and _subtree paths_. Fixed paths don't end with trailing slash, whereas subtree paths do end with trailing slash.

`/snippet/view` and `snippet/create/` -> fixed paths -> In Go’s servemux, fixed path patterns like these are only matched (and the corresponding handler is called) when the request URL path exactly matches the fixed path.

`/` -> subtree path -> Subtree path patterns are matched (and the corresponding handler called) whenever the _start_ of a request URL path matches the subtree path. This helps explain why the "/" pattern is acting like a catch-all. The pattern essentially means match a single slash, followed by anything (or nothing at all)

## The `DefaultServeMux`

`http.Handle` and `http.HandleFunc` function allows to register routes without declaring a servemux. Since there is a global variable `var DefaultServeMux = NewServeMux()` which is initialized by default in the `net/http` package. Since it is a global variable, any package can access it and register a route - including third party packages.

```go
func main() {
    http.HandleFunc("/", home)
    http.HandleFunc("/snippet/view", snippetView)
    http.HandleFunc("/snippet/create", snippetCreate)

    log.Print("Starting server on :4000")
    err := http.ListenAndServe(":4000", nil)
    log.Fatal(err)
}
```

## Chapter 2.4

### Customizing HTTP Headers `(look at additional information for more info)`

```go
func snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// We tell the user which request method are supported for the particular url.
		w.Header().Set("Allow", http.MethodPost)

		// These can be replaced by http.Error since behind the scenes it calls w.Write() and w.WriteHeader()
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		// Here this will have no effect on the output, the map has to be updated before calling the w.Write() and w.WriteHeader() functions
		// w.Header().Set("Allow", "POST")
	}

	w.Write([]byte("Create a snippet..."))
}
```

## Chapter 2.5

### The `io.Writer` interface

```go
func snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// r.URL.Query().Get("id")
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// w.Write([]byte("Display a specific snippet..."))
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
```

The code above introduced another new thing behind-the-scenes. If you take a look at the documentation for the fmt.Fprintf() function you’ll notice that it takes an io.Writer as the first parameter…

```go
func Fprintf(w io.Writer, format string, a ...any) (n int, err error)
```

…but we passed it our http.ResponseWriter object instead — and it worked fine. We’re able to do this because the io.Writer type is an interface, and the http.ResponseWriter object satisfies the interface because it has a w.Write() method. [Read more about concept of interfaces](https://www.alexedwards.net/blog/interfaces-explained)

## Chapter 2.6

**_Read once again to understand why we follow this folder structure better._**

Code before refactoring in single main file.

main.go - single file

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Handlers -> think of it as controllers in MVC pattern. They are used to implement the logic and for writing HTTP response headers and bodies.
// http.ResponseWriter ->  provides methods for assembling a HTTP response and sending it to the user
// http.Request -> a pointer to a struct which holds information about the current request (like the HTTP method and the URL being requested)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	// Check if the current request URL path exactly matches "/". If it doesn't, use
	// the http.NotFound() function to send a 404 response to the client.
	// Importantly, we then return from the handler. If we don't return the handler
	// would keep executing and also write the "Hello from SnippetBox" message.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from snippetbox."))
}

func snippetViewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	// r.URL.Query().Get("id")
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// w.Write([]byte("Display a specific snippet..."))
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func snippetCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		// We tell the user which request method are supported for the particular url.
		w.Header().Set("Allow", http.MethodPost)
		w.Header()["Date"] = nil
		w.Header()["X-XSS-Protection"] = []string{"1; mode=block"}
		// These can be replaced by http.Error since behind the scenes it calls w.Write() and w.WriteHeader()
		// w.WriteHeader(405)
		// w.Write([]byte("Method Not Allowed"))

		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)

		// Here this will have no effect on the output, the map has to be updated before calling the w.Write() and w.WriteHeader() functions
		// w.Header().Set("Allow", "POST")
	}

	w.Write([]byte("Create a snippet..."))
}

func main() {
	fmt.Println("Hello World")

	// servemux is routing in go terminology
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/snippet/view", snippetViewHandler)
	mux.HandleFunc("/snippet/create", snippetCreateHandler)

	log.Print("Starting server on: 4000")
	// ListenAndServe accepts a TCP network address and a handler.
	err := http.ListenAndServe("localhost:4000", mux)
	log.Fatal(err)

}
```

## Chapter 2.7
