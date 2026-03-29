package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"context"
	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/persistence"
	transporthttp "github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/transport/http"
	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
	_ "modernc.org/sqlite"
)

func main() {
	// 1. Initialize DB
	db, err := sql.Open("sqlite", ":memory:?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatalf("failed to open sqlite: %v", err)
	}
	defer db.Close()

	// 2. Setup schema
	if err := persistence.InitDB(db); err != nil {
		log.Fatalf("failed to create schema: %v", err)
	}
	log.Println("Database schema initialized successfully")

	// 3. Initialize Repositories and Services
	repo := persistence.NewSQLiteRepository(db)
	resourceManager := registry.NewResourceManager(repo)
	heartbeatEngine := registry.NewHeartbeatEngine(repo, 0, 0)

	// Start Garbage Collector
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go heartbeatEngine.StartGarbageCollector(ctx)

	// 4. Setup Router
	regHandlers := transporthttp.NewRegistrationHandlers(resourceManager, heartbeatEngine)
	queryHandlers := transporthttp.NewQueryHandlers(repo)
	router := transporthttp.NewRouter(regHandlers, queryHandlers)

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("Starting NMOS Registry on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
