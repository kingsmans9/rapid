package store

import "github.com/pavansokkenagaraj/rapid-agent/pkg/store/postgres"

var s Store

var _ Store = (*postgres.RapidStore)(nil)

func GetStore() Store {
	if s == nil {
		s = &postgres.RapidStore{}
	}
	return s
}
