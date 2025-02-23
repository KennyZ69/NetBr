package main

import (
	"fmt"
	"net"
)

// chooseVictimIP let's the user choose from a list of IP addresses which one to attack and returns the chosen IP from an array
func chooseVictimIP(activeIPs []net.IP) net.IP {
	var idx int

	for {
		fmt.Println("Choose the number of IP to attack")
		fmt.Scanln(&idx)
		if idx > len(activeIPs) || idx < 1 {
			continue
		}
		return activeIPs[idx-1]
	}

}
