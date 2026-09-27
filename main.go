package main

import (
	"log"
	"net/http"
	"os"
	"practice/internal/config"
	"practice/internal/database"
	"practice/internal/domain/notes"
	"practice/internal/handlers"
	"practice/internal/middleware"
	"practice/internal/repository/postgres"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	l := log.New(os.Stdout, "[API] ", log.Default().Flags())

	cfg := config.New()
	db, err := database.Connect(cfg.DSN())
	if err != nil {
		l.Fatalf("failed to connect to db: %v", err)
	}
	l.Println("connected to database")

	if err := db.AutoMigrate(&notes.Note{}); err != nil {
		l.Fatalf("failed to migrate database: %v", err)
	}

	l.Println("migrations successful")

	r := gin.New() // gin.Default() includes logger and recovery
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.AttachRequestIDMiddleware())

	notesRepo := postgres.NewNotesRepo(db)
	noteService := notes.NewService(notesRepo)
	notesHandler := handlers.NewNotesHandler(noteService, l)

	r.GET("/health", handlers.HealthHandler)
	r.GET("/notes/:id", notesHandler.GetNoteByIDHandler())
	r.POST("/notes", notesHandler.CreateNoteHandler())
	r.GET("/notes", notesHandler.ListNotesHandler())
	r.DELETE("/notes/:id", notesHandler.DeleteNoteHandler())

	srv := http.Server{
		Addr:         ":9090",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
	}

	l.Println("Starting server on :9090")
	if err := srv.ListenAndServe(); err != nil {
		l.Fatalf("Failed to start server: %v", err.Error())
	}
}
