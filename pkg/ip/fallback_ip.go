package ip

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func fallBackPing(srcIP, targetIP string, timeout time.Duration) (bool, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", srcIP)
	if err != nil {
		return false, fmt.Errorf("Error creating ICMP connection: %s", err)
	}
	defer conn.Close()

	dst, err := net.ResolveIPAddr("ip4", targetIP)
	if err != nil {
		return false, fmt.Errorf("Error resolving target: %s", err)
	}

	randID := rand.Intn(65535)
	randSeq := rand.Intn(65535)

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{ID: randID, Seq: randSeq, Data: []byte("Hello my precious!")},
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return false, fmt.Errorf("Error marshaling ICMP message: %s", err)
	}

	_, err = conn.WriteTo(msgBytes, dst)
	if err != nil {
		return false, fmt.Errorf("Error sending ICMP request: %s", err)
	}

	err = conn.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return false, fmt.Errorf("Error setting deadline: %s", err)
	}

	reply := make([]byte, 2048)

	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		return false, nil
	}

	receivedMsg, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return false, fmt.Errorf("Error parsing ICMP response: %s", err)
	}

	if echoRep, ok := receivedMsg.Body.(*icmp.Echo); ok {
		if echoRep.ID == randID && echoRep.Seq == randSeq && peer.String() == targetIP {
			return true, nil
		}
	}

	return false, nil
}

func fallBackPinging(cidr *net.IPNet, srcIP string) []net.IP {
	var activeIPs []net.IP
	var count = 0
	var wg sync.WaitGroup
	chanIP := make(chan net.IP)

	for ip := cidr.IP.Mask(cidr.Mask); cidr.Contains(ip); IncIP(ip) {

		ipCopy := make(net.IP, len(ip))
		copy(ipCopy, ip)

		wg.Add(1)

		go func(targetIP net.IP) {

			defer wg.Done()

			if targetIP.Equal(cidr.IP) {
				return
			}

			active, err := fallBackPing(srcIP, targetIP.String(), 2*time.Second)
			if err != nil {
				return
			}

			if active {
				count++
				chanIP <- targetIP
				log.Printf("Found active IP: %s (%d)\n", targetIP.String(), count)
			}

		}(ipCopy)
	}

	go func() {
		wg.Wait()
		close(chanIP)
	}()

	for ip := range chanIP {
		activeIPs = append(activeIPs, ip)
	}

	return activeIPs
}
