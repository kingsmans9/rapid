// swagger:route GET /sites service listSites
//
// # List all sites
//
// List all sites
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Schemes: http
//
//	Deprecated: false
//
//	Security:
//
//	Responses:
package handlers

import (
	"fmt"
	"net/http"
)

type ListSitesResponse struct {
	Sites []Site `json:"sites"`
}

type Site struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	EmailAddress string `json:"email_address"`
	Description  string `json:"description"`
	ISVSiteID    string `json:"isv_site_id"`
}

func (s *site) ListSites(w http.ResponseWriter, r *http.Request) {
	sites, err := s.store.ListSites(r.Context())
	if err != nil {
		JSONWithError(w, http.StatusInternalServerError, fmt.Sprintf("error listing sites: %v", err))
		return
	}

	response := ListSitesResponse{
		Sites: make([]Site, len(sites)),
	}

	for i, site := range sites {
		response.Sites[i] = Site{
			ID:           site.SiteID,
			Name:         site.SiteName,
			EmailAddress: site.EmailAddress,
			Description:  site.SiteDescription,
			ISVSiteID:    site.ISVSiteID,
		}
	}

	JSON(w, http.StatusOK, response)
}
