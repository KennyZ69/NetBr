package pkg

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/mdlayher/arp"
)

func GetOwnMac(ifi *net.Interface) (string, error) {
	return ifi.HardwareAddr.String(), nil
}

func ArpScan(ifi *net.Interface, ipNet *net.IPNet) (map[string]string, error) {

	c, err := arp.Dial(ifi)
	if err != nil {
		return nil, fmt.Errorf("could not dial arp client: %s\n", err)
	}
	defer c.Close()

	log.Println("Dialed the arp client")

	// make a map for the mac : ip relations
	ret := make(map[string]string)

	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		if ip.Equal(ipNet.IP) {
			continue
		}

		nIP, err := netipIP(&ip)
		if err != nil {
			return nil, err
		}

		mac, err := c.Resolve(*nIP)
		if err != nil {
			continue
		}

		ret[ip.String()] = mac.String()
		fmt.Printf("Found %s : %s\n", ip.String(), mac.String())

		time.Sleep(50 * time.Millisecond)
	}

	return ret, nil
}

// TODO:
// finish these possible implementations or the arp mapping
// func ScanARP(ifi *net.Interface) error {
// 	handle, err := pcap.OpenLive(ifi.Name, 65536, true, pcap.BlockForever)
// 	if err != nil {
// 		return err
// 	}
// 	defer handle.Close()
//
// 	stop := make(chan struct{})
// 	go readARP(handle, ifi, stop)
// 	defer close(stop)
// }
//
// func readARP(handle *pcap.Handle, ifi *net.Interface, stop chan struct{}) {
// 	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
// }
