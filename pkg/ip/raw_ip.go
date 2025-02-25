package ip

import (
	"log"
	"net"
	"time"

	netlibk "github.com/KennyZ69/netlibK"
)

func ICMPClient(ifi *net.Interface) (*netlibk.Client, error) {
	c, err := netlibk.ICMPSetClient(ifi)
	if err != nil {
		return c, err
	}
	defer c.Close()

	if err = c.Conn.SetDeadline(time.Now().Add(time.Duration(2))); err != nil {
		return c, err
	}

	return c, nil
}

func ListIPs(cidr *net.IPNet, ifi *net.Interface) ([]net.IP, error) {
	c, err := ICMPClient(ifi)
	if err != nil {
		return nil, err
	}

	var activeIPs []net.IP
	var count int = 0

	for ip := cidr.IP.Mask(cidr.Mask); cidr.Contains(ip); IncIP(ip) {
		if ip.Equal(cidr.IP) {
			continue
		}

		_, active, err := c.Ping(ip, []byte("Hello victim!"))
		if err != nil {
			continue
		}

		if active {
			count++
			activeIPs = append(activeIPs, ip)
			log.Printf("Found active IP: %s (%d)\n", ip.String(), count)
		}
	}

	// if len(activeIPs) < 1 {
	// fallBackPinging(cidr)
	// }

	return activeIPs, nil
}
