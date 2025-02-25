package main

import (
	"fmt"
	"log"

	"github.com/KennyZ69/netBr/pkg"
	ipkg "github.com/KennyZ69/netBr/pkg/ip"
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
			Gateway: gateway,
		}
		if err = pkg.SaveConf(cfg); err != nil {
			log.Fatalf("Error saving config file: %s\n", err)
		}
	}

	fmt.Printf("Gotten ifi: %s\n", cfg.NetIfi.Name)
	fmt.Printf("Gotten ip: %s\n", cfg.LocalIP)
	fmt.Printf("Gotten mac: %s\n", cfg.Mac)
	fmt.Printf("Gotten CIDR: %s\n", cfg.CIDR)
	fmt.Printf("Gotten Gateway IP: %s\n", cfg.Gateway.String())

	// activeIPs, err := pkg.ListIPs(cfg.CIDR, cfg.NetIfi)
	// // // activeIPs, err := pkg.HighListIPs(cfg.CIDR, cfg.NetIfi)
	// if err != nil {
	// 	log.Fatalf("Error in the ListIPs func: %s\n", err)
	// } else if activeIPs == nil {
	// 	log.Fatalf("No active IP addresses were found\n")
	// }

	activeIPs, err := ipkg.HighListIPs(cfg.CIDR, cfg.NetIfi, cfg.LocalIP)
	if err != nil {
		log.Fatalf("Error somewhere in HighListIPs: %s\n", err)
	} else if activeIPs == nil {
		log.Fatalf("No active IP addresses were found\n")
	}

	fmt.Println("Active IPs:", activeIPs)

	// TODO:
	// Now I could print out the active IP addresses and let the attacker choose which one to intercept
	// so I can ping the ips using my netlibk and map its corresponding mac to it
}
