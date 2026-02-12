package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/vlah-software-house/etsy-mock-api/internal/handlers"
	"github.com/vlah-software-house/etsy-mock-api/internal/middleware"
	"github.com/vlah-software-house/etsy-mock-api/internal/seed"
	"github.com/vlah-software-house/etsy-mock-api/internal/store"
)

func main() {
	port := flag.Int("port", 8080, "Server port")
	noAuth := flag.Bool("no-auth", false, "Disable authentication checks")
	noSeed := flag.Bool("no-seed", false, "Start with empty data store")
	seedConfig := flag.String("seed-config", "", "Path to JSON seed config for generated content")
	flag.Parse()

	s := store.New()
	if *seedConfig != "" {
		cfg, err := seed.LoadConfig(*seedConfig)
		if err != nil {
			log.Fatalf("Failed to load seed config: %v", err)
		}
		seed.GenerateFromConfig(s, cfg)
		log.Printf("Generated seed data from %s (%d shops)", *seedConfig, cfg.Shops)
	} else if !*noSeed {
		seed.Load(s)
		log.Println("Seed data loaded")
	}

	tokenStore := middleware.NewTokenStore()
	h := handlers.New(s, tokenStore)
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	var handler http.Handler = mux
	handler = middleware.JSONContent(handler)
	handler = middleware.RateLimitHeaders(handler)
	if !*noAuth {
		handler = middleware.MockAuth(tokenStore)(handler)
	}
	handler = middleware.CORS(handler)
	handler = middleware.RequestLogger(handler)

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("Etsy Mock API server starting on %s", addr)
	log.Printf("Base URL: http://localhost:%d/v3/application", *port)
	log.Printf("Health check: http://localhost:%d/ping", *port)
	log.Printf("OAuth token: POST http://localhost:%d/v3/public/oauth/token", *port)
	if *noAuth {
		log.Println("Authentication: DISABLED")
	} else {
		log.Println("Authentication: enabled (x-api-key format: keystring:shared_secret)")
		log.Println("Pre-seeded OAuth tokens: test-token-alice (user 1001), test-token-bob (user 1002)")
	}

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
