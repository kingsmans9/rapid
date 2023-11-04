//go:generate swagger generate spec

package apiserver

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/spectrocloud/rapid-agent/pkg/handlers"
	"github.com/spectrocloud/rapid-agent/pkg/middleware"
	"github.com/spectrocloud/rapid-agent/pkg/persistence"
	"github.com/spectrocloud/rapid-agent/pkg/store"
	"github.com/spectrocloud/rapid-agent/pkg/util/filesystem"
)

type APIServerParams struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	PostgresSSLMode  string
	RapidDataDir     string
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

	fs := filesystem.NewFileSystem()
	rapidStore := store.GetStore()
	siteHandler := handlers.NewSiteHandler(rapidStore, fs, params.RapidDataDir)

	r := mux.NewRouter()

	// Add CORS middleware
	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(middleware.RecoveryMiddleware)

	unauthenticatedRouter := r.PathPrefix("").Subrouter()
	RegisterUnauthenticatedRoutes(unauthenticatedRouter)

	// TODO: Add authentication middleware
	authenticatedRouter := r.PathPrefix("").Subrouter()
	authenticatedRouter.Use(middleware.LoggingMiddleware)
	RegisterAuthenticatedRoutes(authenticatedRouter, siteHandler)

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

func RegisterAuthenticatedRoutes(r *mux.Router, siteHandler handlers.SiteHandler) {
	r.Name("ListSites").Path("/api/v1/sites").Methods("GET").HandlerFunc(siteHandler.ListSites)
	r.Name("GetData").Path("/api/v1/data").Methods("GET").HandlerFunc(siteHandler.GetData)
}
