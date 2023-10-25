//go:generate swagger generate spec

package apiserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pavansokkenagaraj/rapid-agent/pkg/handlers"
	"github.com/pavansokkenagaraj/rapid-agent/pkg/middleware"
	"github.com/pavansokkenagaraj/rapid-agent/pkg/persistence"
	"github.com/pavansokkenagaraj/rapid-agent/pkg/store"
)

type APIServerParams struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string
}

func Start(params *APIServerParams) {
	_, err := persistence.InitDB(params.PostgresHost, params.PostgresPort, params.PostgresUser, params.PostgresPassword, params.PostgresDatabase, params.PostgresSSLMode)
	if err != nil {
		log.Println("error getting db session")
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	if err := store.GetStore().WaitForReady(ctx); err != nil {
		log.Println("error waiting for ready")
		panic(err)
	}
	cancel()

	rapidStore := store.GetStore()
	siteHandler := handlers.NewSiteHandler(rapidStore)

	r := mux.NewRouter()

	// Add CORS middleware
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.RecoveryMiddleware)

	unauthenticatedRouter := r.PathPrefix("").Subrouter()
	RegisterUnauthenticatedRoutes(unauthenticatedRouter)

	// TODO: Add authentication middleware
	authenticatedRouter := r.PathPrefix("").Subrouter()
	authenticatedRouter.Use(middleware.LoggingMiddleware)
	RegisterAuthenticatedRoutes(authenticatedRouter, rapidStore, siteHandler)

	srv := &http.Server{
		Handler: r,
		Addr:    ":3000",
	}

	log.Println("Starting server on port 3000")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func RegisterUnauthenticatedRoutes(r *mux.Router) {
	r.HandleFunc("/healthz", handlers.Healthz)
}

func RegisterAuthenticatedRoutes(r *mux.Router, store store.Store, siteHandler handlers.SiteHandler) {
	r.Name("ListServices").Path("/api/v1/sites").Methods("GET").HandlerFunc(siteHandler.ListSites)
}
