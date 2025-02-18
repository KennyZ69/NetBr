package pkg

import (
	"fmt"
	"log"
	"net"
)

// getInterface returns a pointer to the user's network interface or nil and error
func GetInterface() (*net.Interface, error) {
	ifis, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Error getting list of net interfaces: %s\n", err)
	}

	for _, ifi := range ifis {
		// check whether the interface is down or is a loopback
		if ifi.Flags&net.FlagLoopback != 0 || ifi.Flags&net.FlagUp == 0 {
			continue
		}
		return &ifi, nil
	}

	return nil, fmt.Errorf("Error: Couldn't get the network interface\n")
}
