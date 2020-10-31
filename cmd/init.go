package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
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
			return xerrors.Errorf("failed to create config directory: %w", err)
		}
		if err := generateConfigFile(); err != nil {
			return xerrors.Errorf("failed to generate config file: %w", err)
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
		return xerrors.Errorf("%s already exist", configDir())
	}
	if err := os.MkdirAll(configDir(), 0775); err != nil {
		return xerrors.Errorf("failed to create directory at %s. %s", configDir(), err)
	}
	return nil
}

func generateConfigFile() error {
	file, err := os.OpenFile(configFilePath(), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return xerrors.Errorf("failed to create config file at %s. %s", configFilePath(), err)
	}
	defer file.Close()

	_, err = fmt.Fprint(file, configTemplate)
	return xerrors.Errorf("failed to write config file: %w", err)
}

func init() {
	rootCmd.AddCommand(initCmd)
}
