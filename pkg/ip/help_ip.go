package ip

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/KennyZ69/netBr/pkg"
)

const (
	RouteFile = "/proc/net/route"
)

// GetIpRange iterates over ip addresses and returns the local one along with the CIDR notation
func GetIpRange(ifi *net.Interface) (string, *net.IPNet, error) {
	addrs, err := ifi.Addrs()
	if err != nil {
		return "", nil, err
	}

	for _, addr := range addrs {
		ip, err := getLocalIP(addr)
		if err != nil {
			continue
		}
		net, err := getLocalCIDR(addr)
		if err != nil {
			continue
		}
		return ip, net, nil

	}
	return "", nil, fmt.Errorf("Did not find any ip address\n")
}

func getLocalIP(addr net.Addr) (string, error) {
	var ip net.IP

	switch t := addr.(type) {
	case *net.IPNet:
		ip = t.IP
	case *net.IPAddr:
		ip = t.IP
	}

	if ip != nil && ip.To4() != nil && !ip.IsLoopback() {
		return ip.String(), nil
	}

	return "", fmt.Errorf("No valid IPv4 could be found as local")
}
func getLocalCIDR(addr net.Addr) (*net.IPNet, error) {
	if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
		return ipNet, nil
	}
	return nil, fmt.Errorf("No valid IPNet could be found")
}

// GetGateway returns the ip of the gateway as a string (from /proc/net/route file) and error
func GetGateway() (net.IP, error) {
	// fmt.Println("Getting gateway")
	bytes, err := pkg.ReadFile(RouteFile)
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

func IncIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}
