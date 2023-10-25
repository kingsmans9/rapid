package store

import (
	"context"

	"github.com/pavansokkenagaraj/rapid-agent/pkg/store/types"
)

type Store interface {
	SiteStore

	WaitForReady(ctx context.Context) error
}

type SiteStore interface {
	ListSites(ctx context.Context) ([]*types.Site, error)
	GetSite(ctx context.Context, id string) (*types.Site, error)
}
