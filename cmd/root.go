package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	Version = "dev-master"
)

var rootCmd = &cobra.Command{
	Use:     "unifi-doorbell-chime",
	Short:   "Notify UniFi UVC G4 Doorbell ringing",
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}

func configDir() string {
	return filepath.Join(os.Getenv("HOME"), ".unifi-doorbell-chime")
}

func configFilePath() string {
	if configFile == "" {
		return filepath.Join(configDir(), "config.yaml")
	}
	return configFile
}
