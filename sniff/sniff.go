package sniff

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

// TODO:
// make the sniff func that would sniff packets on my device but filter out for the intercepted ones from promisc mode; mitm
// and run the mitm with this sniffer along each other
func Sniff(ifi *net.Interface, attackMAC, victimIP, gatewayIP string) {
	handle, err := pcap.OpenLive(ifi.Name, 1600, false, pcap.BlockForever)
	if err != nil {
		log.Fatalf("Error opening handle to sniff packets: %s\n", err)
	}

	defer handle.Close()

	logFp := logFPath()

	logF, err := os.OpenFile(logFp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening %s log file: %s\n", logFp, err)
	}
	defer logF.Close()

	packetSrc := gopacket.NewPacketSource(handle, handle.LinkType())

	log.Println("Sniffing the intercepted packets ... ")

	for packet := range packetSrc.Packets() {
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}
		ipPacket, ok := ipLayer.(*layers.IPv4)
		if !ok {
			continue
		}

		ethLayer := packet.Layer(layers.LayerTypeEthernet)
		if ethLayer == nil {
			continue
		}
		ethPacket, ok := ethLayer.(*layers.Ethernet)
		if !ok {
			continue
		}

		if ipPacket.SrcIP.String() == victimIP && ethPacket.DstMAC.String() == attackMAC {
			t := time.Now().Format("01-02-2000 10:01:02")
			info := fmt.Sprintf("\n[%s] Intercepted packet\nSource: %s\nDestination: %s (Original dest = %s)\n)", t, ipPacket.SrcIP.String(), ipPacket.DstIP.String(), gatewayIP)

			fmt.Print(info)

			if _, err := logF.WriteString(info); err != nil {
				// do some if there is an error
			}
		}
	}
}

func logFPath() string {
	usr, err := user.Current()
	// if the current user throws an error, it falls back to a local dir config file
	if err != nil {
		return "./sniffer.log"
	}
	// else It creates the path for directory of the current app and its files
	return filepath.Join(usr.HomeDir, ".config", "netBr", "sniffer.log")
}
