package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sawadashota/unifi-doorbell-chime/assets"

	"github.com/spf13/cobra"
)

const configTemplate = `---
unifi:
  ip: "192.168.1.1"
  username: "username"
  password: "password"

#boot_option:
#  mac_address: 00:00:00:00:00:00

message:
  templates:
    - "I'm on my way"
    - "I'm busy now"
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate config file and assets",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := createConfigDir(); err != nil {
			return err
		}
		if err := generateConfigFile(); err != nil {
			return err
		}
		if err := copyAssets(); err != nil {
			return err
		}

		fmt.Printf("created %s successfully\n\n", configDir())
		fmt.Printf("$ vi %s\n\n", configFilePath())
		fmt.Printf("then exec\n\n")
		fmt.Println("$ unifi-doorbell-chime start")
		return nil
	},
}

func createConfigDir() error {
	if _, err := os.Stat(configDir()); err == nil {
		return fmt.Errorf("%s aleady exist\n", configDir())
	}
	if err := os.MkdirAll(configDir(), 0775); err != nil {
		return fmt.Errorf("failed to create directory at %s. %s\n", configDir(), err)
	}
	return nil
}

func generateConfigFile() error {
	file, err := os.OpenFile(configFilePath(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to create config file at %s. %s", configFilePath(), err)
	}
	defer file.Close()

	_, err = fmt.Fprint(file, configTemplate)
	return err
}

func copyAssets() error {
	dir := filepath.Join(configDir(), "assets")
	if err := os.MkdirAll(dir, 0775); err != nil {
		return fmt.Errorf("failed to create directory at %s. %s\n", configDir(), err)
	}
	icon := assets.New().AppIcon()

	dest := filepath.Join(dir, "AppIcon.png")
	file, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return fmt.Errorf("failed to put assets at %s. %s", dest, err)
	}
	defer file.Close()

	_, err = file.Write(icon)
	return err
}

func init() {
	rootCmd.AddCommand(initCmd)
}
