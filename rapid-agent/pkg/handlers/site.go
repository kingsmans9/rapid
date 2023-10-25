package handlers

import (
	"net/http"

	"github.com/pavansokkenagaraj/rapid-agent/pkg/store"
)

type SiteHandler interface {
	ListSites(w http.ResponseWriter, r *http.Request)
}

type site struct {
	store store.Store
}

func NewSiteHandler(store store.Store) SiteHandler {
	return &site{
		store: store,
	}
}
