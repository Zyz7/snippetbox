package main

import (
	"crypto/tls"
	"example/internal/config"
	"example/internal/models"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	debug          bool
	logger         *slog.Logger
	snippets       *models.SnippetModel
	users          *models.UserModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {

	// Define a new command line flag with the values
	addr := flag.String("addr", ":4000", "HTTP network address")
	debug := flag.Bool("debug", false, "Enable debug mode")

	// You need to call this before you use the addr variable
	flag.Parse()
	godotenv.Load("../../.env")

	//Use to initialize a new structured logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := config.OpenDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	// Initialize a new template cache
	templateCache, err := newTemplateCache()

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initialize a decoder instance
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour
	// Make sure that the Secure attribute is set for https connection
	sessionManager.Cookie.Secure = true

	// Initialize a new instance of our application struct
	app := &application{
		debug:          *debug,
		logger:         logger,
		snippets:       &models.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Print a log message to say that the server is starting
	//log.Printf("starting server on: %s", *addr)

	// In this case the only thing that we are changing in the curve preferences
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	srv := &http.Server{
		Addr:    *addr,
		Handler: app.routes(),
		// Writes los entries at error level and assign it to the errorlog field
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	/*
		Use the ListenAndServe to start a new web server
		We pass in two parameters: the TCP network address to listen on
		this case 4000 and the servermux we just created
		If returns an error we use the Fatal function to log the error
		message and exit
		Note that any error returned is always non-nil
	*/
	//err = srv.ListenAndServe()
	// Ocupar el comando echo 'tls/' >> .gitignore para evitar mandar los certificados
	err = srv.ListenAndServeTLS("../../tls/cert.pem", "../../tls/key.pem")

	// And we also use to log any error message returned
	logger.Error(err.Error())
	//log.Fatal(err)

	os.Exit(1)
}
