package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mohammadmokh/Movie-API/internal/data"

	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {

	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "application port")
	flag.Parse()

	cfg.db.dsn = os.Getenv("POSTGRES_URI")

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg.db.dsn)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModel(db),
	}

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Handler:      app.routes(),
	}

	logger.Printf("starting application on %d", app.config.port)
	err = srv.ListenAndServe()
	logger.Fatal(err)
}
