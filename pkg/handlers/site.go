package handlers

import (
	"net/http"

	"github.com/spectrocloud/rapid-agent/pkg/rapid"
	"github.com/spectrocloud/rapid-agent/pkg/store"
	"github.com/spectrocloud/rapid-agent/pkg/util/filesystem"
)

type SiteHandler interface {
	ListSites(w http.ResponseWriter, r *http.Request)
	GetData(w http.ResponseWriter, r *http.Request)
}

type site struct {
	rapidDataDir    string
	store           store.Store
	fs              filesystem.Interface
	rapidDataCopier rapid.RapidDataCopier
}

func NewSiteHandler(store store.Store, fs filesystem.Interface, rapidDataDir string) SiteHandler {
	return &site{
		rapidDataDir:    rapidDataDir,
		store:           store,
		fs:              fs,
		rapidDataCopier: rapid.NewRapidDataCopier(rapidDataDir, fs),
	}
}
