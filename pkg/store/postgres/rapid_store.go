package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/spectrocloud/rapid-agent/pkg/persistence"
	"github.com/spectrocloud/rapid-agent/pkg/util/logger"
)

type RapidStore struct {
}

func (r *RapidStore) WaitForReady(ctx context.Context) error {
	err := waitForPostgres(ctx)
	if err != nil {
		fmt.Printf("error waiting for postgres: %v\n", err)
		return err
	}
	return nil
}

func waitForPostgres(ctx context.Context) error {
	logger.Debug("waiting for database to be ready")

	period := 1 * time.Second
	for {
		db, err := persistence.GetDBSession()
		if err != nil {
			logger.Debugf("error getting db session: %v", err)
		} else {

			query := `select count(1) from sites`
			row := db.QueryRow(query)
			var count int
			err = row.Scan(&count)
			if err != nil {
				logger.Debugf("error scanning row: %v", err)
			} else {
				logger.Debug("database is ready")
				return nil
			}
		}

		select {
		case <-time.After(period):
			continue
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for database to be ready")
		}
	}
}
