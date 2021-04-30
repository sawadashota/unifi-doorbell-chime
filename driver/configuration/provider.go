package configuration

type Provider interface {
	LogLevel() string

	UnifiSkipTLSVerify() bool
	UnifiIp() string
	UnifiUsername() string
	UnifiPassword() string

	WebPort() int
	APIPort() int

	MessageList() []string

	BootOptionMacAddress() string
}
