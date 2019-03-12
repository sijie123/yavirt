package netx

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/juju/errors"
)

var (
	host   string
	nicIps = []string{}
)

func init() {
	var err error
	host, err = os.Hostname()
	if err != nil {
		panic(err)
	}

	var addrs []net.Addr
	addrs, err = net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}

	for _, ifaddr := range addrs {
		var ip net.IP
		switch typ := ifaddr.(type) {
		case *net.IPNet:
			ip = typ.IP
		case *net.IPAddr:
			ip = typ.IP
		}

		if ip.IsGlobalUnicast() {
			nicIps = append(nicIps, ip.String())
		}
	}
}

func GetLocalIP(network, laddr string) (string, error) {
	switch network {
	case "tcp", "tcp4", "tcp6":
		var tcpaddr, err = net.ResolveTCPAddr(network, laddr)
		if err != nil {
			return "", errors.Trace(err)
		}
		if tcpaddr.Port < 1 {
			return "", errors.Errorf("unexpectedly, resolve %s addr %s to %s", network, laddr, tcpaddr.String())
		}
		if len(nicIps) < 1 {
			return "", errors.New("unknown IP addr")
		}
		return net.JoinHostPort(nicIps[0], strconv.Itoa(tcpaddr.Port)), nil

	case "unix", "unixpacket":
		return laddr, nil

	default:
		return "", net.UnknownNetworkError(network)
	}
}

func IntToIPv4(v int64) string {
	var ui32 = 0xffffffff & v
	var sb strings.Builder

	for i := uint(24); i <= 24; i -= 8 {
		var seg = ui32 >> i
		fmt.Fprintf(&sb, "%d.", 0xff&seg)
	}

	return strings.TrimRight(sb.String(), ".")
}

func IPv4ToInt(ipv4 string) (int64, error) {
	var ui32 int64

	for i, elem := range strings.Split(ipv4, ".") {
		var v, err = strconv.Atoi(elem)
		if err != nil {
			return 0, errors.Trace(err)
		}

		ui32 |= int64(v << uint(24-i*8))
	}

	return ui32, nil
}
