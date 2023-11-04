package handlers

import (
	"context"
	"fmt"

	"github.com/spectrocloud/rapid-agent/pkg/store/types"
)

type mockErrStore struct {
}

func (m *mockErrStore) WaitForReady(ctx context.Context) error {
	return nil
}

func (m *mockErrStore) ListSites(ctx context.Context) ([]*types.Site, error) {
	return nil, fmt.Errorf("error listing sites")
}

func (m *mockErrStore) GetSite(ctx context.Context, id string) (*types.Site, error) {
	return nil, fmt.Errorf("error getting sites")
}

func (m *mockErrStore) DumpData(ctx context.Context) (map[string]interface{}, error) {
	return nil, nil
}
