package postgres

import (
	"context"
	"database/sql"
	"fmt"

	persistence "github.com/spectrocloud/rapid-agent/pkg/persistence"
	"github.com/spectrocloud/rapid-agent/pkg/store/types"
)

func (r *RapidStore) listTables(ctx context.Context) ([]string, error) {
	db, err := persistence.GetDBSession()
	if err != nil {
		return nil, fmt.Errorf("error getting db session: %v", err)
	}

	query := `SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE';`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying tables: %v", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		tables = append(tables, table)
	}

	return tables, nil
}

func (r *RapidStore) DumpData(ctx context.Context) (map[string]interface{}, error) {
	tables, err := r.listTables(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing tables: %v", err)
	}

	data := make(map[string]interface{})
	for _, table := range tables {
		data[table], err = r.dumpTable(ctx, table)
		if err != nil {
			return nil, fmt.Errorf("error dumping table %s: %v", table, err)
		}
	}

	return data, nil
}

func (r *RapidStore) dumpTable(ctx context.Context, table string) ([]map[string]interface{}, error) {
	db, err := persistence.GetDBSession()
	if err != nil {
		return nil, fmt.Errorf("error getting db session: %v", err)
	}

	query := fmt.Sprintf(`SELECT * FROM %s`, table)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying table %s: %v", table, err)
	}
	defer rows.Close()

	var data []map[string]interface{}
	for rows.Next() {
		row, err := scanRow(rows)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		data = append(data, row)
	}

	return data, nil
}

func scanRow(rows *sql.Rows) (map[string]interface{}, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %v", err)
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, fmt.Errorf("error scanning row: %v", err)
	}

	row := make(map[string]interface{})
	for i, col := range columns {
		val := values[i]

		b, ok := val.([]byte)
		if ok {
			row[col] = string(b)
		} else {
			row[col] = val
		}
	}

	return row, nil
}

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
