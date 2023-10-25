package handlers

import (
	"context"
	"fmt"

	"github.com/pavansokkenagaraj/rapid-agent/pkg/store/types"
)

type mockErrStore struct {
}

func (m *mockErrStore) WaitForReady(ctx context.Context) error {
	return nil
}

func (m *mockErrStore) ListSites(ctx context.Context) ([]*types.Site, error) {
	return nil, fmt.Errorf("error listing services")
}

func (m *mockErrStore) GetSite(ctx context.Context, id string) (*types.Site, error) {
	return nil, fmt.Errorf("error getting service")
}
