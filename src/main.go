package main

import (
	"fmt"
	"log"

	"github.com/KennyZ69/netBr/pkg"
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
		ip, cidr, err := pkg.GetIpRange(ifi)
		if err != nil {
			log.Fatalf("Error: Could not get local ip address: %s\n", err)
		}

		cfg = &pkg.Config{
			NetIfi:  ifi,
			LocalIP: ip,
			Mac:     ownMac,
			CIDR:    cidr,
		}
		if err = pkg.SaveConf(cfg); err != nil {
			log.Fatalf("Error saving config file: %s\n", err)
		}
	}

	fmt.Printf("Gotten ifi: %s\n", cfg.NetIfi.Name)
	fmt.Printf("Gotten ip: %s\n", cfg.LocalIP)
	fmt.Printf("Gotten mac: %s\n", cfg.Mac)
	fmt.Printf("Gotten CIDR: %s\n", cfg.CIDR)

}
