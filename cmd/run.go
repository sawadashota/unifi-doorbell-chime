package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sawadashota/unifi-doorbell-chime/driver"
	"github.com/sawadashota/unifi-doorbell-chime/x/wifimac"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var runCmd = &cobra.Command{
	Use:   "start",
	Short: "start listen to doorbell ringing",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		viper.SetConfigFile(configFilePath())

		if err := viper.ReadInConfig(); err != nil {
			return fmt.Errorf(`config file not found because "%s"`, err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, _ []string) error {
		d := driver.NewDefaultDriver()
		i := newInstance(d)

		// wait for Ctrl-C
		sigCh := make(chan os.Signal, 2)
		signal.Notify(sigCh, os.Interrupt)
		signal.Notify(sigCh, syscall.SIGTERM)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		errCh := make(chan error, 1)
		go func() {
			if d.Configuration().BootOptionMacAddress() != "" {
				if err := i.bootWithMacAddressObservation(ctx); err != nil {
					errCh <- err
				}
				return
			}
			if err := i.boot(ctx); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-sigCh:
			return nil
		case err := <-errCh:
			d.Registry().Logger().Debugf("%+v", err)
			return err
		}
	},
}

type instance struct {
	d driver.Driver
}

func newInstance(d driver.Driver) *instance {
	return &instance{
		d: d,
	}
}

func (i *instance) bootWithMacAddressObservation(ctx context.Context) error {
	const checkMacAddressInterval = 1 * time.Minute

	i.d.Registry().Logger().Infof("boot only when mac address is %s", i.d.Configuration().BootOptionMacAddress())

	for {
		ma, err := wifimac.GetMacAddress()
		if err != nil {
			i.d.Registry().Logger().Infof(
				"it seem no network connected. waiting until mac address to be %s",
				i.d.Configuration().BootOptionMacAddress(),
			)
			time.Sleep(checkMacAddressInterval)
			continue
		}
		if ma.String() == i.d.Configuration().BootOptionMacAddress() {
			err := i.boot(ctx)
			if err == nil {
				return nil
			}
			ma2, err2 := wifimac.GetMacAddress()
			if err2 != nil {
				i.d.Registry().Logger().Infof(
					"listener stopped because no network connected. waiting until mac address to be %s",
					i.d.Configuration().BootOptionMacAddress(),
				)
				time.Sleep(checkMacAddressInterval)
				continue
			}

			if ma2.String() == i.d.Configuration().BootOptionMacAddress() {
				return errors.WithMessage(err, "unexpected error occurred")
			}

			i.d.Registry().Logger().Infof(
				"listener stopped because current mac address is %s. waiting to be %s",
				ma2.String(),
				i.d.Configuration().BootOptionMacAddress(),
			)
			time.Sleep(checkMacAddressInterval)
			continue
		}

		i.d.Registry().Logger().Infof(
			"current mac address is %s. waiting to be %s",
			ma.String(),
			i.d.Configuration().BootOptionMacAddress(),
		)
		time.Sleep(checkMacAddressInterval)
	}
}

func (i *instance) boot(ctx context.Context) error {
	listenerErrCh := make(chan error, 1)
	go func() {
		if err := i.d.Registry().Listener().Start(); err != nil {
			listenerErrCh <- errors.Wrap(err, "unexpected error occurred")
		}
	}()

	webFrontendErrCh := make(chan error, 1)
	go func() {
		if err := i.d.Registry().WebFrontendServer().Start(); err != nil {
			webFrontendErrCh <- errors.Wrap(err, "unexpected error occurred")
		}
	}()

	webApiErrCh := make(chan error, 1)
	go func() {
		if err := i.d.Registry().WebApiServer().Start(); err != nil {
			webApiErrCh <- errors.Wrap(err, "unexpected error occurred")
		}
	}()

	select {
	case <-ctx.Done():
		var eg errgroup.Group
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		eg.Go(func() error {
			return i.d.Registry().Listener().Shutdown(ctx)
		})
		eg.Go(func() error {
			return i.d.Registry().WebFrontendServer().Shutdown(ctx)
		})
		eg.Go(func() error {
			return i.d.Registry().WebApiServer().Shutdown(ctx)
		})
		if err := eg.Wait(); err != nil {
			i.d.Registry().Logger().Error(err)
		}
		return errors.WithStack(ctx.Err())

	case err := <-listenerErrCh:
		var eg errgroup.Group
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		eg.Go(func() error {
			return i.d.Registry().WebFrontendServer().Shutdown(ctx)
		})
		eg.Go(func() error {
			return i.d.Registry().WebApiServer().Shutdown(ctx)
		})
		if err := eg.Wait(); err != nil {
			i.d.Registry().Logger().Error(err)
		}
		return errors.WithStack(err)

	case err := <-webFrontendErrCh:
		var eg errgroup.Group
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		eg.Go(func() error {
			return i.d.Registry().Listener().Shutdown(ctx)
		})
		eg.Go(func() error {
			return i.d.Registry().WebApiServer().Shutdown(ctx)
		})
		if err := eg.Wait(); err != nil {
			i.d.Registry().Logger().Error(err)
		}
		return errors.WithStack(err)

	case err := <-webApiErrCh:
		var eg errgroup.Group
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		eg.Go(func() error {
			return i.d.Registry().Listener().Shutdown(ctx)
		})
		eg.Go(func() error {
			return i.d.Registry().WebFrontendServer().Shutdown(ctx)
		})
		if err := eg.Wait(); err != nil {
			i.d.Registry().Logger().Error(err)
		}
		return errors.WithStack(err)
	}
}

var configFile string

func init() {
	runCmd.PersistentFlags().StringVarP(
		&configFile,
		"config",
		"c",
		"",
		"Config file. Default is $HOME/.unifi-doorbell-chime/config.yaml",
	)
	rootCmd.AddCommand(runCmd)
}
