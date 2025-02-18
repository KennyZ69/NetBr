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
		ifi, err := pkg.GetInterface()
		if err != nil {
			log.Fatalf("Error: Could not get interface in main func: %s\n", err)
		}

		ownMac, err := pkg.GetOwnMac(ifi)
		if err != nil {
			log.Fatalf("Error: Could not get own mac address: %s\n", err)
		}

		ip, err := pkg.GetOwnIP(ifi)
		if err != nil {
			log.Fatalf("Error: Could not get local ip address: %s\n", err)
		}

		cfg = &pkg.Config{
			NetIfi:  ifi.Name,
			LocalIP: ip,
			Mac:     ownMac,
		}
		if err = pkg.SaveConf(cfg); err != nil {
			log.Fatalf("Error saving config file: %s\n", err)
		}
	}

	fmt.Printf("Gotten ifi: %s\n", cfg.NetIfi)

	// _, err = pkg.CapturePackets(cfg.NetIfi)
	// if err != nil {
	// 	log.Fatalf("Error capturing packets: %s\n", err)
	// }

	macToIp, err := pkg.GetMac()
	if err != nil {
		log.Fatalf("Error getting macs to ip: %s\n", err)
	}

	fmt.Println(macToIp)
}
