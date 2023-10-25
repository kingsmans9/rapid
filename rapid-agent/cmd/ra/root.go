package main

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rapid-agent",
		Short: "Rapid Agent",
		Long:  ``,
		Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	cobra.OnInitialize(initConfig)
	cmd.PersistentFlags().String("log-level", "info", "set the log level")
	cmd.AddCommand(APICmd())
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	return cmd
}

func initConfig() {
	viper.SetEnvPrefix("RAPID_AGENT")
	viper.AutomaticEnv()
}
