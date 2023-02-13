package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Divyue30597/snippetbox-lets-go/internal/models"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	// Add a snippets field to the application struct. This will allow us to
	// make the SnippetModel object available to our handlers.
	snippets *models.SnippetModel
	// initialize this cache in the main() function and make it available to our handlers as a
	// dependency via the application struct
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	dsn := flag.String("dsn", "host=localhost port=5432 dbname=snippetbox user=postgres password=1234", "Postgresql data source name")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr
	// variable. You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000". If any errors are
	// encountered during parsing the application will be terminated.
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Initialize the new template
	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	// Initializing a new instance of application struct, containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		// Initialize a models.SnippetModel instance and add it to the application dependencies.
		snippets: &models.SnippetModel{
			DB: db,
		},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
