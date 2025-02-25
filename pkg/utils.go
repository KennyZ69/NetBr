package pkg

import (
	"io"
	"net"
	"net/netip"
	"os"
)

func netipIP(ip *net.IP) (*netip.Addr, error) {
	addr, err := netip.ParseAddr(ip.String())
	if err != nil {
		return nil, err
	}

	byt := addr.As4()
	ret := netip.AddrFrom4(byt)
	return &ret, nil
}

func ReadFile(file string) ([]byte, error) {
	f, err := os.Open(file) // this opens the file just for reading
	if err != nil {
		return nil, err
	}
	defer f.Close()

	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func IncIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}
