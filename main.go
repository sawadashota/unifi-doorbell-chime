package main

import (
	"log"

	"github.com/sawadashota/unifi-doorbell-chime/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
