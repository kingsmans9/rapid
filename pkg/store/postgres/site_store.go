package postgres

import (
	"context"
	"database/sql"
	"fmt"

	persistence "github.com/pavansokkenagaraj/rapid-agent/pkg/persistence"
	"github.com/pavansokkenagaraj/rapid-agent/pkg/store/types"
)

func (r *RapidStore) ListSites(ctx context.Context) ([]*types.Site, error) {
	db, err := persistence.GetDBSession()
	if err != nil {
		return nil, fmt.Errorf("error getting db session: %v", err)
	}

	query := `SELECT site_id, site_name, email_address, site_description, isv_site_id, site_uid FROM sites`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying sites: %v", err)
	}
	defer rows.Close()

	var sites []*types.Site
	for rows.Next() {
		site, err := scanSitesRow(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		sites = append(sites, site)
	}

	return sites, nil
}

func scanSitesRow(rows *sql.Rows) (*types.Site, error) {
	var site types.Site
	var emailAddress, siteDescription, isvSiteID, siteUID sql.NullString
	if err := rows.Scan(&site.SiteID, &site.SiteName, &emailAddress, &siteDescription, &isvSiteID, &siteUID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	if emailAddress.Valid {
		site.EmailAddress = emailAddress.String
	}
	if siteDescription.Valid {
		site.SiteDescription = siteDescription.String
	}
	if isvSiteID.Valid {
		site.ISVSiteID = isvSiteID.String
	}
	if siteUID.Valid {
		site.SiteUID = siteUID.String
	}

	return &site, nil
}

func (r *RapidStore) GetSite(ctx context.Context, id string) (*types.Site, error) {
	db, err := persistence.GetDBSession()
	if err != nil {
		return nil, fmt.Errorf("error getting db session: %v", err)
	}

	query := `SELECT site_id, email_address, site_description, isv_site_id, site_uid FROM sites WHERE site_uid = $1`
	row := db.QueryRowContext(ctx, query, id)

	return scanSiteRow(row)
}

func scanSiteRow(row *sql.Row) (*types.Site, error) {
	var site types.Site
	var emailAddress, siteDescription, isvSiteID, siteUID sql.NullString
	if err := row.Scan(&site.SiteID, &site.SiteName, &emailAddress, &siteDescription, &isvSiteID, &siteUID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	if emailAddress.Valid {
		site.EmailAddress = emailAddress.String
	}
	if siteDescription.Valid {
		site.SiteDescription = siteDescription.String
	}
	if isvSiteID.Valid {
		site.ISVSiteID = isvSiteID.String
	}
	if siteUID.Valid {
		site.SiteUID = siteUID.String
	}

	return &site, nil
}
