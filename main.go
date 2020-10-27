package main

import (
	"os"

	"github.com/sawadashota/unifi-doorbell-chime/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
