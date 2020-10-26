package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
	"github.com/sawadashota/unifi-doorbell-chime/driver"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "unifi-proto-chime",
	Short: "Notify UniFi UVC G4 Doorbell ringing",
	RunE: func(cmd *cobra.Command, _ []string) error {
		d, err := driver.NewDefaultDriver()
		if err != nil {
			return errors.WithStack(err)
		}

		// wait for Ctrl-C
		sigCh := make(chan os.Signal, 2)
		signal.Notify(sigCh, os.Interrupt)
		signal.Notify(sigCh, syscall.SIGTERM)

		listenerErrCh := make(chan error, 1)
		go func() {
			if err := d.Registry().Listener().Start(); err != nil {
				listenerErrCh <- errors.Wrap(err, "unexpected error occurred")
			}
		}()

		webFrontendErrCh := make(chan error, 1)
		go func() {
			if err := d.Registry().WebFrontendServer().Start(); err != nil {
				webFrontendErrCh <- errors.Wrap(err, "unexpected error occurred")
			}
		}()

		webApiErrCh := make(chan error, 1)
		go func() {
			if err := d.Registry().WebApiServer().Start(); err != nil {
				webApiErrCh <- errors.Wrap(err, "unexpected error occurred")
			}
		}()

		select {
		case <-sigCh:
			var eg errgroup.Group
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			eg.Go(func() error {
				return d.Registry().Listener().Shutdown(ctx)
			})
			eg.Go(func() error {
				return d.Registry().WebFrontendServer().Shutdown(ctx)
			})
			eg.Go(func() error {
				return d.Registry().WebApiServer().Shutdown(ctx)
			})
			if err := eg.Wait(); err != nil {
				d.Registry().Logger().Error(err)
			}
		case err := <-listenerErrCh:
			var eg errgroup.Group
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			eg.Go(func() error {
				return d.Registry().WebFrontendServer().Shutdown(ctx)
			})
			eg.Go(func() error {
				return d.Registry().WebApiServer().Shutdown(ctx)
			})
			if err := eg.Wait(); err != nil {
				d.Registry().Logger().Error(err)
			}
			d.Registry().Logger().Error(err)
		case err := <-webFrontendErrCh:
			var eg errgroup.Group
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			eg.Go(func() error {
				return d.Registry().Listener().Shutdown(ctx)
			})
			eg.Go(func() error {
				return d.Registry().WebApiServer().Shutdown(ctx)
			})
			if err := eg.Wait(); err != nil {
				d.Registry().Logger().Error(err)
			}
			d.Registry().Logger().Error(err)
		case err := <-webApiErrCh:
			var eg errgroup.Group
			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
			eg.Go(func() error {
				return d.Registry().Listener().Shutdown(ctx)
			})
			eg.Go(func() error {
				return d.Registry().WebFrontendServer().Shutdown(ctx)
			})
			if err := eg.Wait(); err != nil {
				d.Registry().Logger().Error(err)
			}
			d.Registry().Logger().Error(err)
		}

		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}

var configFile string

func init() {
	rootCmd.PersistentFlags().StringVarP(
		&configFile,
		"config",
		"c",
		"",
		"Config file. Default is $HOME/.unifi-doorbell-chime.yaml",
	)
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if configFile == "" {
		configFile = filepath.Join(os.Getenv("HOME"), ".unifi-doorbell-chime.yaml")
		if _, err := os.Stat(configFile); err != nil {
			if err := createSampleConfig(configFile); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("created sample config file at %s successfully\n", configFile)
			fmt.Print("try after configuration\n\n")
			fmt.Printf("$ vi %s\n\n", configFile)
			os.Exit(1)
		}
		viper.SetConfigFile(configFile)
	}
	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf(`Config file not found because "%s"`, err)
		os.Exit(1)
	}
}

const sampleConfig = `---
unifi:
  ip: "192.168.1.1"
  username: "username"
  password: "password"

message:
  templates:
    - "I'm on my way"
    - "I'm busy now"
`

func createSampleConfig(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()

	src := bytes.NewBufferString(sampleConfig)
	if _, err := io.Copy(f, src); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
