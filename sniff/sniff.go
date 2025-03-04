package sniff

import "net"

const logFp = "~/.config/netBr/sniffer.log" // log file path

// TODO:
// make the sniff func that would sniff packets on my device but filter out for the intercepted ones from promisc mode; mitm
// and run the mitm with this sniffer along each other
func Sniff(ifi *net.Interface, attackMAC, victimIP, gatewayIP string) {

}
