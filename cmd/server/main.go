package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/arturfil/aggregator-script/db"
	"github.com/arturfil/aggregator-script/helpers"
	"github.com/arturfil/aggregator-script/services/scrappers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/joho/godotenv/autoload"
)

type AppServer struct {
	addr string
	db   *sql.DB
}

// NewAppServer - constructor for AppServer
func NewAppServer(addr string, db *sql.DB) *AppServer {
	return &AppServer{
		addr: addr,
		db:   db,
	}
}

func (app *AppServer) Serve() error {
	helpers.MessageLogs.InfoLog.Println("API listenting on port", app.addr)

	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

    // add services & routes here
    newsStore := scrappers.NewStore(app.db)
    newsHandler := scrappers.NewHandler(newsStore)
    newsHandler.RegisterRoutes(router)

	srv := &http.Server{Addr: fmt.Sprintf("%s", app.addr),
		Handler: router,
	}

	return srv.ListenAndServe()
}

func main() {
	dsn := os.Getenv("DSN")
	port := os.Getenv("PORT")

	db, err := db.NewDatabase(dsn)
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}

	server := NewAppServer(fmt.Sprintf(":%s", port), db.Client)
	if err := server.Serve(); err != nil {
		log.Fatal(err)
	}
}
