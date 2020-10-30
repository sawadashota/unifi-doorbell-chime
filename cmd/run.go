package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/pkg/errors"
	"github.com/sawadashota/unifi-doorbell-chime/driver"
	"github.com/sawadashota/unifi-doorbell-chime/x/wifimac"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

		var eg errgroup.Group
		defer func() {
			if err := eg.Wait(); err != nil {
				d.Registry().Logger().Error(err)
			}
		}()

		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()

		errCh := make(chan error, 1)
		eg.Go(func() error {
			if d.Configuration().BootOptionMacAddress() != "" {
				return i.bootWithMacAddressObservation(ctx)
			}
			return i.boot(ctx)
		})

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
	desired := i.d.Configuration().BootOptionMacAddress()
	i.d.Registry().Logger().Infof("boot only when mac address is %s", desired)

	for {
		ma, err := wifimac.GetMacAddress()
		if err != nil {
			i.d.Registry().Logger().Infof("it seem no network connected. waiting until mac address to be %s", desired)
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
					desired,
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
				desired,
			)
			time.Sleep(checkMacAddressInterval)
			continue
		}

		i.d.Registry().Logger().Infof("current mac address is %s. waiting to be %s", ma.String(), desired)
		time.Sleep(checkMacAddressInterval)
	}
}

func (i *instance) boot(ctx context.Context) error {
	var eg errgroup.Group
	defer func() {
		if err := eg.Wait(); err != nil {
			i.d.Registry().Logger().Error(err)
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, len(i.d.Registry().Services()))
	for _, svc := range i.d.Registry().Services() {
		s := svc
		eg.Go(func() error {
			if err := s.Start(ctx); err != nil {
				errCh <- err
				return err
			}
			return nil
		})
	}

	select {
	case <-ctx.Done():
		return errors.WithStack(ctx.Err())
	case err := <-errCh:
		cancel()
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
