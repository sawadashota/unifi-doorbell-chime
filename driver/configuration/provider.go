package configuration

type Provider interface {
	LogLevel() string

	UnifiSkipTLSVerify() bool
	UnifiIp() string
	UnifiUsername() string
	UnifiPassword() string

	NotificationIcon() string

	WebPort() uint64
	ApiPort() uint64

	MessageList() []string
}
