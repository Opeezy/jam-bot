package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	infoLog    *log.Logger
	errorLog   *log.Logger
	warningLog *log.Logger
	traceLog   *log.Logger
)

func main() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Create loggers for each level
	infoLog = log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warningLog = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	traceLog = log.New(logFile, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)

	err = godotenv.Load()
	if err != nil {
		errorLog.Println("Error loading env file.")
		traceLog.Fatal(err)
	}

	token := os.Getenv("TOKEN")
	if token == "" {
		errorLog.Fatal("Empty token.")
	}

	infoLog.Println("Initializing discord client...")
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		errorLog.Println("Closing application due to failed client initialization.")
		traceLog.Fatal(err)
	}

	infoLog.Println("Opening connection...")
	err = discord.Open()
	if err != nil {
		errorLog.Println("Unable to open connection.")
		traceLog.Fatal(err)
	}

	log.Println("Bot is now running. Press CTRL-C to exit.")
	infoLog.Println("Bot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}
