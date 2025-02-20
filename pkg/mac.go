package pkg

import (
	"fmt"
	"net"
	"time"

	"github.com/mdlayher/arp"
)

// func GetOwnMac(ifi string) (string, error) {
func GetOwnMac(ifi *net.Interface) (string, error) {
	// netIfi, err := net.InterfaceByName(ifi)
	// if err != nil {
	// 	return "", err
	// }

	return ifi.HardwareAddr.String(), nil
}

func ArpScan(ifi *net.Interface, ipNet *net.IPNet) (map[string]string, error) {

	c, err := arp.Dial(ifi)
	if err != nil {
		return nil, fmt.Errorf("could not dial arp client: %s\n", err)
	}
	defer c.Close()

	// make a map for the mac : ip relations
	ret := make(map[string]string)

	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		if ip.Equal(ipNet.IP) {
			continue
		}

		nIP, err := netipIP(&ip)
		if err != nil {
			return nil, err
		}

		mac, err := c.Resolve(*nIP)
		if err != nil {
			continue
		}

		ret[ip.String()] = mac.String()
		fmt.Printf("Found %s : %s\n", ip.String(), mac.String())

		time.Sleep(50 * time.Millisecond)
	}

	return ret, nil
}

// GetIpRange iterates over ip addresses and returns the first one along with the CIDR notation
func GetIpRange(ifi *net.Interface) (string, *net.IPNet, error) {
	addrs, err := ifi.Addrs()
	if err != nil {
		return "", nil, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), ipNet, nil
			}
		}
	}
	return "", nil, fmt.Errorf("Did not find any ip address\n")
}
