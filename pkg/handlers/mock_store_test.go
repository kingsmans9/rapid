package handlers

import (
	"context"

	"github.com/spectrocloud/rapid-agent/pkg/store/types"
)

type mockStore struct {
}

func (m *mockStore) WaitForReady(ctx context.Context) error {
	return nil
}

func (m *mockStore) ListSites(ctx context.Context) ([]*types.Site, error) {
	return []*types.Site{
		{
			SiteID:          "site1",
			SiteName:        "site1",
			EmailAddress:    "site1@site.com",
			SiteDescription: "site1",
			ISVSiteID:       "site1",
		},
	}, nil
}

func (m *mockStore) GetSite(ctx context.Context, id string) (*types.Site, error) {
	return &types.Site{
		SiteID:          "site1",
		SiteName:        "site1",
		EmailAddress:    "site1@site.com",
		SiteDescription: "site1",
		ISVSiteID:       "site1",
	}, nil
}

func (m *mockStore) DumpData(ctx context.Context) (map[string]interface{}, error) {
	return nil, nil
}
