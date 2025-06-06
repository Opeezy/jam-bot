package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var (
	port            = flag.String("port", "", "The server port to listen on")
	environment     = flag.String("env", "", "The current app environment")
	spotifyClientId = flag.String("spclientid", "", "Spotify web app client id")

	infoLog       *log.Logger
	errorLog      *log.Logger
	warningLog    *log.Logger
	middlewareLog *log.Logger

	templates = template.Must(template.ParseGlob("public/*.html"))
)

func init() {
	flag.Parse()
}

func main() {
	log.Println("creating loggers...")
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	infoLog = log.New(multiWriter, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(multiWriter, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(multiWriter, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	middlewareLog = log.New(multiWriter, "MIDDLEWARE: ", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("environment=%s", *environment)

	if *port == "" {
		errorLog.Fatal("port not specified")
	} else if *spotifyClientId == "" {
		errorLog.Fatal("spotify client id not specified")
	}

	router := chi.NewRouter()
	router.Use(MiddleWareLogger)

	fileServer := http.FileServer(http.Dir("./styles"))

	router.Route("/api/v1/", func(r chi.Router) {
		r.Get("/register", registerHandler)
		r.Get("/login", loginHandler)
		r.Get("/callback", callbackHandler)
	})
	router.Handle("/styles/*", http.StripPrefix("/styles/", fileServer))

	infoLog.Printf("server listening on port: %s", *port)
	http.ListenAndServe(fmt.Sprintf(":%s", *port), router)
}

func MiddleWareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		middlewareLog.Printf("%s %s %s from %s", r.Method, r.RequestURI, duration, r.RemoteAddr)
	})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "register.html", nil)
	if err != nil {
		errorLog.Printf("error at indexHandler: %s", err)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()

	values := url.Values{}
	values.Add("client_id", *spotifyClientId)
	values.Add("response_type", "code")
	values.Add("redirect_uri", "http://127.0.0.1:8080/api/v1/callback")
	values.Add("scope", "user-read-private")
	values.Add("state", state)

	queryString := values.Encode()

	url := "https://accounts.spotify.com/authorize?" + queryString
	infoLog.Println(url)

	http.Redirect(w, r, "https://google.com/", http.StatusSeeOther)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Callback!"))
}
