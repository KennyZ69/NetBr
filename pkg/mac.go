package pkg

import (
	"fmt"
	"net"
)

// func GetOwnMac(ifi string) (string, error) {
func GetOwnMac(ifi *net.Interface) (string, error) {
	// netIfi, err := net.InterfaceByName(ifi)
	// if err != nil {
	// 	return "", err
	// }

	fmt.Println("Mac: ", ifi.HardwareAddr.String())
	return ifi.HardwareAddr.String(), nil
}

func GetMac() (map[string]string, error) {
	ifis, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// make a map for the mac : ip relations
	ret := make(map[string]string)

	for _, ifi := range ifis {
		mac := ifi.HardwareAddr.String()
		if mac != "" {
		}
		addrs, err := ifi.Addrs()
		if err != nil {
			continue
			// return nil, err
		}

		for _, addr := range addrs {
			switch addr.(type) {
			case *net.IPAddr:
				ret[mac] = addr.String()
			}
		}
	}

	return ret, nil
}

func GetOwnIP(ifi *net.Interface) (string, error) {
	addrs, err := ifi.Addrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipNet.IP.To4() != nil {
				fmt.Println("IP: ", ipNet.IP.String())
				return ipNet.String(), nil
			}
		}
	}
	return "", fmt.Errorf("Did not find any ip address\n")
}
