package configuration

import (
	"os"
	"path/filepath"

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

	viperNotificationIcon = "notification.icon"

	viperWebPort = "proto.port"
	viperApiPort = "api.port"

	viperMessageTemplates = "message.templates"
)

func getString(key string, defaultValue string) string {
	v := viper.GetString(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func getUnit64(key string, defaultValue uint64) uint64 {
	v := viper.Get(key)
	if v == nil {
		return defaultValue
	}
	return viper.GetUint64(key)
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

func (v *ViperProvider) NotificationIcon() string {
	return getString(
		viperNotificationIcon,
		filepath.Join(os.Getenv("HOME"), ".unifi-doorbell-chime/assets/AppIcon.png"),
	)
}

func (v *ViperProvider) WebPort() uint64 {
	return getUnit64(viperWebPort, 33319)
}

func (v *ViperProvider) ApiPort() uint64 {
	return getUnit64(viperApiPort, 33320)
}

func (v *ViperProvider) MessageList() []string {
	return viper.GetStringSlice(viperMessageTemplates)
}
