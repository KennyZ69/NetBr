package ip

import (
	"log"
	"net"
	"time"

	netlibk "github.com/KennyZ69/netlibK"
)

func HighListIPs(cidr *net.IPNet, ifi *net.Interface, srcIP string) ([]net.IP, error) {
	var activeIPs []net.IP

	for ip := cidr.IP.Mask(cidr.Mask); cidr.Contains(ip); IncIP(ip) {
		if ip.Equal(cidr.IP) {
			continue
		}

		log.Printf("Pinging %s\n", ip.String())
		_, active, err := netlibk.HigherLvlPing(ip, []byte("Hello victim!"), time.Duration(2))
		if err != nil {
			continue
		}

		if active {
			activeIPs = append(activeIPs, ip)
			log.Printf("Found active IP: %s\n", ip.String())
		}
	}

	if len(activeIPs) < 1 {
		fallBackPinging(cidr, srcIP)
	}

	return activeIPs, nil
}
