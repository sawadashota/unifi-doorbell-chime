package wifimac

import (
	"net"

	"golang.org/x/xerrors"
)

func GetMacAddress() (*net.HardwareAddr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
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
			return nil, xerrors.Errorf(": %w", err)
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
				return nil, xerrors.Errorf(": %w", err)
			}
			return &netInterface.HardwareAddr, nil
		}
	}
	return nil, xerrors.New("not found network")
}
