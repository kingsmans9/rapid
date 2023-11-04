package main

import (
	"fmt"
	"strings"

	"github.com/spectrocloud/rapid-agent/pkg/apiserver"
	"github.com/spectrocloud/rapid-agent/pkg/util/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func APICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api",
		Short: "Starts the Rapid Agent server",
		Long:  ``,
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()

			if v.GetString("log-level") == "debug" {
				logger.SetDebug()
			}

			params := apiserver.APIServerParams{
				PostgresHost:     v.GetString("postgres-host"),
				PostgresPort:     v.GetString("postgres-port"),
				PostgresUser:     v.GetString("postgres-user"),
				PostgresPassword: v.GetString("postgres-password"),
				PostgresDatabase: v.GetString("postgres-database"),
				PostgresSSLMode:  v.GetString("postgres-sslmode"),
				RapidDataDir:     v.GetString("rapid-data-dir"),
			}

			switch {
			case params.PostgresHost == "":
				return fmt.Errorf("postgres-host is required")
			case params.PostgresPort == "":
				return fmt.Errorf("postgres-port is required")
			case params.PostgresUser == "":
				return fmt.Errorf("postgres-user is required")
			case params.PostgresPassword == "":
				return fmt.Errorf("postgres-password is required")
			case params.PostgresDatabase == "":
				return fmt.Errorf("postgres-database is required")
			case params.PostgresSSLMode == "":
				return fmt.Errorf("postgres-sslmode is required")
			case params.RapidDataDir == "":
				return fmt.Errorf("rapid-data-dir is required")
			}

			apiserver.Start(&params)
			return nil
		},
	}

	cmd.Flags().String("postgres-host", "", "postgres host")
	cmd.Flags().String("postgres-port", "", "postgres port")
	cmd.Flags().String("postgres-user", "", "postgres user")
	cmd.Flags().String("postgres-password", "", "postgres password")
	cmd.Flags().String("postgres-database", "", "postgres database")
	cmd.Flags().String("postgres-sslmode", "", "postgres sslmode")
	cmd.Flags().String("rapid-data-dir", "", "rapid data directory")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return cmd
}
