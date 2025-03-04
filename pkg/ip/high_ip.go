package ip

import (
	"log"
	"net"
	"time"

	"github.com/KennyZ69/netBr/pkg"
	netlibk "github.com/KennyZ69/netlibK"
)

// func HighListIPs(cidr *net.IPNet, ifi *net.Interface, srcIP string) ([]net.IP, error) {
func HighListIPs(cfg pkg.Config) ([]net.IP, error) {
	var activeIPs []net.IP

	for ip := cfg.CIDR.IP.Mask(cfg.CIDR.Mask); cfg.CIDR.Contains(ip); IncIP(ip) {
		if ip.Equal(cfg.CIDR.IP) || ip.String() == cfg.Gateway {
			continue
		}

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
		log.Printf("Had a problem detecting active hosts using netlibK\nMoving on the builtin libraries (fallback)\n")
		activeIPs = fallBackPinging(cfg.CIDR, cfg.LocalIP)
	}

	return activeIPs, nil
}
