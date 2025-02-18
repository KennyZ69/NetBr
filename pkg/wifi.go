package pkg

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

type ParsedPacket struct {
	WifiFrame WifiFrame
	// Bacon string `json:"bacon"` // meant as beacon (or an access point)
	Bacon *layers.Dot11MgmtBeacon
	// data  []byte
	data *layers.Dot11Data
}

type WifiFrame struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
}

func CapturePackets(ifi string) ([]gopacket.Packet, error) {
	handle, err := pcap.OpenLive(ifi, 1024, true, pcap.BlockForever)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	var packets []gopacket.Packet
	pSrc := gopacket.NewPacketSource(handle, handle.LinkType())
	for p := range pSrc.Packets() {
		// log.Println("Getting a packet")
		packets = append(packets, p)
		ParseWifiPacket(p)
		fmt.Println(p)
	}

	return packets, nil
}

func ParseWifiPacket(packet gopacket.Packet) ParsedPacket {
	// log.Println("Parsing a packet")
	var retPacket ParsedPacket
	if baconLayer := packet.Layer(layers.LayerTypeDot11MgmtBeacon); baconLayer != nil {
		bacon, ok := baconLayer.(*layers.Dot11MgmtBeacon)
		if ok {
			log.Printf("Bacon: %v\n", bacon)
			retPacket.Bacon = bacon
		}
	}

	if wifiLayer := packet.Layer(layers.LayerTypeDot11); wifiLayer != nil {
		wifi, ok := wifiLayer.(*layers.Dot11)
		if ok {
			log.Printf("Wifi frame: %v\n", wifi)
			frame := WifiFrame{
				Source:      wifi.Address2.String(),
				Destination: wifi.Address1.String(),
			}
			retPacket.WifiFrame = frame
		}
	}

	if dataLayer := packet.Layer(layers.LayerTypeDot11Data); dataLayer != nil {
		data, ok := dataLayer.(*layers.Dot11Data)
		if ok {
			log.Printf("Data: %v\n", data)
			retPacket.data = data
		}
	}

	return retPacket
}

func printPacket(packet ParsedPacket) {
	// make this func to nicely print parsed packets
}
