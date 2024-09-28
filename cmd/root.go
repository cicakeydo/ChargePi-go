package cmd

import (
	"github.com/ChargePi/ChargePi-go/internal/pkg/configuration"
	"github.com/ChargePi/ChargePi-go/pkg/observability"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "chargepi",
	Short: "ChargePi is an open-source Charge point project.",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	cobra.OnInitialize(func() {
		observability.SetupLogging(log.StandardLogger(), observability.Logging{}, viper.GetBool(configuration.Debug))
	})

	rootCmd.AddCommand(runCommand())
	rootCmd.AddCommand(versionCommand())
	rootCmd.AddCommand(exportCommand())
	rootCmd.AddCommand(importCommand())

	rootCmd.PersistentFlags().BoolP(configuration.DebugFlag, "d", false, "debug mode")
	_ = viper.BindPFlag(configuration.Debug, rootCmd.PersistentFlags().Lookup(configuration.DebugFlag))
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.WithError(err).Fatal("Unable to run")
	}
}
