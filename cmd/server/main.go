package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"context"
	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/mdns"
	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/persistence"
	transporthttp "github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/transport/http"
	"github.com/Woord-En-Lewe/nmos_registry/internal/infrastructure/transport/websocket"
	"github.com/Woord-En-Lewe/nmos_registry/internal/registry"
	_ "modernc.org/sqlite"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Initialize DB
	db, err := sql.Open("sqlite", "file::memory:?cache=shared&_pragma=foreign_keys(1)")
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
	subscriptionEngine := registry.NewSubscriptionEngine(repo)
	wsManager := websocket.NewManager(subscriptionEngine)
	subscriptionEngine.SetNotifier(wsManager)

	go wsManager.Run()

	// Start Garbage Collector
	go heartbeatEngine.StartGarbageCollector(ctx)

	// 4. Setup Router
	regHandlers := transporthttp.NewRegistrationHandlers(resourceManager, heartbeatEngine)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	queryHandlers := transporthttp.NewQueryHandlers(repo, subscriptionEngine, wsManager, port)
	router := transporthttp.NewRouter(regHandlers, queryHandlers)

	// 5. Start mDNS Announcer
	hostname, _ := os.Hostname()
	mdnsConfig := mdns.NewConfig(hostname, 8080, 8080)
	mdnsAnnouncer, err := mdns.NewAnnouncer(mdnsConfig)
	if err != nil {
		log.Printf("Warning: failed to create mDNS announcer: %v", err)
	} else {
		go func() {
			if err := mdnsAnnouncer.Start(ctx); err != nil {
				log.Printf("mDNS announcement error: %v", err)
			}
		}()
		defer mdnsAnnouncer.Stop()
	}

	// 6. Start Server
	addr := ":" + port
	log.Printf("Starting NMOS Registry on %s", addr)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
