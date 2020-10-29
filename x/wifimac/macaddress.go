package wifimac

import (
	"net"

	"github.com/pkg/errors"
)

func GetMacAddress() (*net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			netInterface, err := net.InterfaceByName(iface.Name)
			if err != nil {
				return nil, errors.WithStack(err)
			}
			return &netInterface.HardwareAddr, nil
		}
	}
	return nil, errors.New("not found network")
}
