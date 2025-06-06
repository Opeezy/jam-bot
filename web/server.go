package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	port        = flag.String("port", "", "The server port to listen on")
	environment = flag.String("env", "", "The current app environment")

	infoLog    *log.Logger
	errorLog   *log.Logger
	warningLog *log.Logger
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

	infoLog.Printf("environment=%s", *environment)

	if *port == "" {
		errorLog.Fatal("port not specified")
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	infoLog.Printf("server listening on port: %s", *port)
	http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
}
