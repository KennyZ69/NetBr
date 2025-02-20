package pkg

import (
	"net"
	"net/netip"
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

// TODO
func incIP(ip net.IP) {

}
