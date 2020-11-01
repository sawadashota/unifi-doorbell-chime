package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sawadashota/unifi-doorbell-chime/driver"
	"github.com/sawadashota/unifi-doorbell-chime/x/wifimac"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
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

		var eg errgroup.Group
		defer func() {
			if err := eg.Wait(); err != nil {
				d.Registry().Logger().Error(err)
			}
		}()

		ctx, stop := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		errCh := make(chan error, 1)
		eg.Go(func() error {
			if d.Configuration().BootOptionMacAddress() != "" {
				return i.bootWithMacAddressObservation(ctx)
			}
			return i.boot(ctx)
		})

		select {
		case <-ctx.Done():
			return nil
		case err := <-errCh:
			d.Registry().Logger().Debugf("%+v", err)
			return xerrors.Errorf("exit command: %w", err)
		}
	},
}

type instance struct {
	d      driver.Driver
	logger logrus.FieldLogger
}

func newInstance(d driver.Driver) *instance {
	return &instance{
		d:      d,
		logger: d.Registry().AppLogger("instance"),
	}
}

func (i *instance) bootWithMacAddressObservation(ctx context.Context) error {
	desired := i.d.Configuration().BootOptionMacAddress()

	i.logger.Infof("boot only when mac address is %s", desired)
	for {
		if i.isDesiredMacAddress() {
			if err := i.checkWorthToContinue(i.boot(ctx)); err != nil {
				if xerrors.Is(err, complete) {
					return nil
				}
				return err
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func (i *instance) isDesiredMacAddress() bool {
	desired := i.d.Configuration().BootOptionMacAddress()

	ma, err := wifimac.GetMacAddress()
	if err != nil {
		i.logger.Infof("it seem no network connected. waiting until mac address to be %s", desired)
		return false
	}
	if ma.String() != desired {
		i.logger.Infof("current mac address is %s. waiting to be %s", ma.String(), desired)
		return false
	}
	return true
}

var (
	complete = xerrors.New("complete successfully")
)

func (i *instance) checkWorthToContinue(err error) error {
	desired := i.d.Configuration().BootOptionMacAddress()

	if err == nil {
		return complete
	}
	ma2, err2 := wifimac.GetMacAddress()
	if err2 != nil {
		i.logger.Infof(
			"listener stopped because no network connected. waiting until mac address to be %s",
			desired,
		)
		return nil
	}

	if ma2.String() == desired {
		return xerrors.Errorf("unexpected error occurred: %w", err)
	}

	i.logger.Infof("listener stopped because current mac address is %s. waiting to be %s",
		ma2.String(),
		desired,
	)
	return nil
}

func (i *instance) boot(ctx context.Context) error {
	var eg errgroup.Group
	defer func() {
		if err := eg.Wait(); err != nil {
			i.logger.Error(err)
		}
	}()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errCh := make(chan error, len(i.d.Registry().Services()))
	for _, svc := range i.d.Registry().Services() {
		s := svc
		eg.Go(func() error {
			if err := s.Start(ctx); err != nil {
				if xerrors.Is(err, context.Canceled) {
					return nil
				}
				errCh <- err
				return xerrors.Errorf("an error occurred at %T : %w", s, err)
			}
			return nil
		})
	}

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		cancel()
		return xerrors.Errorf("all process canceled: %w", err)
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
