package cmd

import (
	"context"
	"fmt"
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
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
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
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
		"Config file. Default is $HOME/.unifi-doorbell-chime/config.yaml",
	)
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if configFile == "" {
		configFile = filepath.Join(os.Getenv("HOME"), ".unifi-doorbell-chime/config.yaml")
		if _, err := os.Stat(configFile); err != nil {
			fmt.Println("not found $HOME/.unifi-doorbell-chime/config.yaml")
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
