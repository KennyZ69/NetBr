package pkg

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	netlibk "github.com/KennyZ69/netlibK"
)

const (
	RouteFile = "/proc/net/route"
)

func ICMPClient(ifi *net.Interface) (*netlibk.Client, error) {
	c, err := netlibk.ICMPSetClient(ifi)
	if err != nil {
		return c, err
	}
	defer c.Close()

	if err = c.Conn.SetDeadline(time.Now().Add(time.Duration(2))); err != nil {
		return c, err
	}

	return c, nil
}

// TODO:
// I need to print out the list of active IP adresses
// and let the user type in one of them to attack

func ListIPs(cidr *net.IPNet, ifi *net.Interface) ([]net.IP, error) {
	c, err := ICMPClient(ifi)
	if err != nil {
		return nil, err
	}

	var activeIPs []net.IP
	var count int = 0

	for ip := cidr.IP.Mask(cidr.Mask); cidr.Contains(ip); incIP(ip) {
		if ip.Equal(cidr.IP) {
			continue
		}

		_, active, err := c.Ping(ip, []byte("Hello victim!"))
		if err != nil {
			continue
		}

		if active {
			count++
			activeIPs = append(activeIPs, ip)
			log.Printf("Found active IP: %s (%d)\n", ip.String(), count)
		}
	}

	return activeIPs, nil
}

func HighListIPs(cidr *net.IPNet, ifi *net.Interface) ([]net.IP, error) {
	var activeIPs []net.IP

	for ip := cidr.IP.Mask(cidr.Mask); cidr.Contains(ip); incIP(ip) {
		if ip.Equal(cidr.IP) {
			continue
		}
		_, active, err := netlibk.HigherLvlPing(ip, []byte("Hello victim!"), time.Duration(2))
		if err != nil {
			continue
		}

		if active {
			activeIPs = append(activeIPs, ip)
			log.Printf("Found active IP: %s\n", ip.String())
		}
	}

	if len(activeIPs) == 0 {
		return nil, fmt.Errorf("No active IPs found on the network\n")
	}

	return activeIPs, nil
}

// GetIpRange iterates over ip addresses and returns the local one along with the CIDR notation
func GetIpRange(ifi *net.Interface) (string, *net.IPNet, error) {
	addrs, err := ifi.Addrs()
	if err != nil {
		return "", nil, err
	}

	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), ipNet, nil
			}
		}
	}
	return "", nil, fmt.Errorf("Did not find any ip address\n")
}

// GetGateway returns the ip of the gateway as a string (from /proc/net/route file) and error
func GetGateway() (net.IP, error) {
	// fmt.Println("Getting gateway")
	bytes, err := readFile(RouteFile)
	if err != nil {
		return nil, err
	}

	return parseGatewayFile(bytes)
}

// parseGatewayFile takes in the bytes of /proc/net/route file and parses them to return the Gateway IP and error
func parseGatewayFile(file []byte) (net.IP, error) {
	// fmt.Printf("Parsing the %s file\n", RouteFile)
	s := bufio.NewScanner(bytes.NewReader(file))

	if !s.Scan() {
		err := s.Err()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("No gateway found\n")
	}

	for s.Scan() {
		row := s.Text()
		// split each row (token) by tabs ("\t")
		fields := strings.Split(row, "\t")
		// fmt.Printf("Got fields: %s\n", fields)

		if len(fields) < 11 {
			return nil, fmt.Errorf("Invalid file format for %s\n", RouteFile)
		}

		// 1 is the destination and 7 is the mask
		// Iface(0)	Destination(1)	Gateway(2)	Flags	RefCnt	Use	Metric	Mask(7)		MTU	Window	IRTT
		if !(fields[1] == "00000000" && fields[7] == "00000000") {
			continue
		}

		// fmt.Printf("Found the gateway field: %s\n", fields[2])

		// Now return the found IP
		return parseGatewayIPBytes(fields[2])
	}

	return nil, fmt.Errorf("No gateway found\n")
}

// Gets the gateway ip as a string in hex and returns it as net.IP and error
func parseGatewayIPBytes(gateway string) (net.IP, error) {
	// fmt.Println("Parsing the bytes of gateway string")
	ip32, err := strconv.ParseUint(gateway, 16, 32)
	if err != nil {
		return nil, err
	}

	ip := make(net.IP, 4)
	binary.LittleEndian.PutUint32(ip, uint32(ip32))
	return ip, nil
}
