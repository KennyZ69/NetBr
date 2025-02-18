package pkg

import "net"

// func GetOwnMac(ifi string) (string, error) {
func GetOwnMac(ifi *net.Interface) (string, error) {
	// netIfi, err := net.InterfaceByName(ifi)
	// if err != nil {
	// 	return "", err
	// }

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

}
