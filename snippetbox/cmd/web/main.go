package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Divyue30597/snippetbox-lets-go/internal/models"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type application struct {
	// dependency injection
	errorLog *log.Logger
	infoLog  *log.Logger
	// Add a snippets field to the application struct. This will allow us to
	// make the SnippetModel object available to our handlers.

	// Instead of using the SnippetModel and UserModel, replace them with the interface
	// snippets *models.SnippetModel
	// users    *models.UserModel
	snippets models.SnippetModelInterface
	users    models.UserModelInterface
	// initialize this cache in the main() function and make it available to our handlers as a
	// dependency via the application struct
	templateCache  map[string]*template.Template
	sessionManager *scs.SessionManager
	formDecoder    *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP Network Address")

	// database connection -> flag is used because we get to create cmd line flags for the user to input. Default value is already added as second param.
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

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = postgresstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	// Initializing a new instance of application struct, containing the dependencies.
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		// Initialize a models.SnippetModel instance and add it to the application dependencies.
		snippets: &models.SnippetModel{
			DB: db,
		},
		users: &models.UserModel{
			DB: db,
		},
		templateCache:  templateCache,
		sessionManager: sessionManager,
		formDecoder:    formDecoder,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS13,
	}

	// initiating server with all the required fields.
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
		// Setting TLSConfig field
		TLSConfig: tlsConfig,
		// Add Idle, read and write timeouts to the server
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// err = srv.ListenAndServe()
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
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
