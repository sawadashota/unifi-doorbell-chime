package configuration

type Provider interface {
	LogLevel() string

	UnifiSkipTLSVerify() bool
	UnifiIp() string
	UnifiUsername() string
	UnifiPassword() string

	WebPort() uint64
	ApiPort() uint64

	MessageList() []string

	BootOptionMacAddress() string
}
