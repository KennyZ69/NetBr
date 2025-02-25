package ip

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func fallBackPing(srcIP, targetIP string, timeout time.Duration) (bool, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", srcIP)
	if err != nil {
		return false, fmt.Errorf("error creating ICMP connection: %v", err)
	}
	defer conn.Close()

	dst, err := net.ResolveIPAddr("ip4", targetIP)
	if err != nil {
		return false, fmt.Errorf("error resolving target address: %v", err)
	}

	// Create an ICMP echo request message
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{ID: 1, Seq: 1, Data: []byte("PING")},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return false, fmt.Errorf("error marshaling ICMP message: %v", err)
	}

	start := time.Now()

	_, err = conn.WriteTo(msgBytes, dst)
	if err != nil {
		return false, fmt.Errorf("error sending ICMP request: %v", err)
	}

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return false, fmt.Errorf("error setting read deadline: %v", err)
	}

	reply := make([]byte, 1500)

	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return false, nil
	}

	duration := time.Since(start)

	receivedMsg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return false, fmt.Errorf("error parsing ICMP response: %v", err)
	}

	if receivedMsg.Type == ipv4.ICMPTypeEchoReply {
		fmt.Printf("Received reply from %s in %v\n", peer, duration)
		return true, nil
	}

	return false, nil
}

func fallBackPinging(cidr *net.IPNet, srcIP string) []net.IP {
	log.Println("Moving to fallback ping functionality")
	var activeIPs []net.IP
	var count = 0

	for ip := cidr.IP.Mask(cidr.Mask); cidr.Contains(ip); IncIP(ip) {
		if ip.Equal(cidr.IP) {
			continue
		}

		log.Println("Pinging", ip.String())

		active, err := fallBackPing(srcIP, ip.String(), time.Duration(time.Second*2))
		if err != nil {
			continue
		}

		if active {
			count++
			activeIPs = append(activeIPs, ip)
			log.Printf("Found active IP: %s (%d)\n", ip.String(), count)
		}
	}

	return activeIPs
}
