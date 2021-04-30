package configuration

import (
	"github.com/phayes/freeport"
	"github.com/spf13/viper"
)

type ViperProvider struct{}

var _ Provider = new(ViperProvider)

const (
	viperLogLevel = "log.level"

	viperUnifiSkipTLSVerify = "unifi.skip_tls_verify"
	viperUnifiIp            = "unifi.ip"
	viperUnifiUsername      = "unifi.username"
	viperUnifiPassword      = "unifi.password"

	viperWebPort = "web.port"
	viperApiPort = "api.port"

	viperMessageTemplates = "message.templates"

	viperBootOptionMacAddress = "boot_option.mac_address"
)

func getString(key string, defaultValue string) string {
	v := viper.GetString(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func getBool(key string, defaultValue bool) bool {
	v := viper.Get(key)
	if v == nil {
		return defaultValue
	}
	return viper.GetBool(key)
}

func NewViperProvider() Provider {
	return &ViperProvider{}
}

func (v *ViperProvider) LogLevel() string {
	return getString(viperLogLevel, "info")
}

func (v *ViperProvider) UnifiSkipTLSVerify() bool {
	return getBool(viperUnifiSkipTLSVerify, true)
}

func (v *ViperProvider) UnifiIp() string {
	return viper.GetString(viperUnifiIp)
}

func (v *ViperProvider) UnifiUsername() string {
	return viper.GetString(viperUnifiUsername)
}

func (v *ViperProvider) UnifiPassword() string {
	return viper.GetString(viperUnifiPassword)
}

func (v *ViperProvider) WebPort() int {
	port := viper.GetInt(viperWebPort)
	if port == 0 {
		port, _ := freeport.GetFreePort()
		viper.Set(viperWebPort, port)
		return v.WebPort()
	}

	return port
}

func (v *ViperProvider) ApiPort() int {
	port := viper.GetInt(viperApiPort)
	if port == 0 {
		port, _ := freeport.GetFreePort()
		viper.Set(viperApiPort, port)
		return v.ApiPort()
	}

	return port
}

func (v *ViperProvider) MessageList() []string {
	return viper.GetStringSlice(viperMessageTemplates)
}

func (v *ViperProvider) BootOptionMacAddress() string {
	return viper.GetString(viperBootOptionMacAddress)
}
