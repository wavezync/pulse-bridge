package cmd

import (
	"wavezync/pulse-bridge/cmd/pulsebridge"
	"wavezync/pulse-bridge/internal/env"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	configPath string
	host       string
	port       int
)

var rootCmd = &cobra.Command{
	Use:   "pulse-bridge",
	Short: "PulseBridge is a powerful uptime monitoring tool",
	Long:  `PulseBridge exposes internal service status via HTTP, enabling seamless integration with external monitoring tools like Atlassian Statuspage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		envConfig := env.Init()

		if cmd.Flags().Changed("config") {
			envConfig.ConfigPath = configPath
		}
		if cmd.Flags().Changed("host") {
			envConfig.Host = host
		}
		if cmd.Flags().Changed("port") {
			envConfig.Port = port
		}

		return pulsebridge.Run(envConfig)
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "config.yml", "Path to configuration file")
	rootCmd.PersistentFlags().StringVar(&host, "host", "0.0.0.0", "Host address to bind the server")
	rootCmd.PersistentFlags().IntVar(&port, "port", 8080, "Port to run the server")
}

func Execute() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Failed to execute root command")
	}
}
