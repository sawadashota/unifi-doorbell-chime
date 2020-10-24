package driver

import (
	"github.com/pkg/errors"
	"github.com/sawadashota/unifi-doorbell-chime/driver/configuration"
)

type Driver interface {
	Registry() Registry
	Configuration() configuration.Provider
}

type DefaultDriver struct {
	r Registry
	c configuration.Provider
}

var _ Driver = new(DefaultDriver)

func (d *DefaultDriver) Registry() Registry {
	return d.r
}

func (d *DefaultDriver) Configuration() configuration.Provider {
	return d.c
}

func NewDefaultDriver() (Driver, error) {
	c := configuration.NewViperProvider()
	r, err := NewDefaultRegistry(c)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &DefaultDriver{
		c: c,
		r: r,
	}, nil
}
