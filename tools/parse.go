package tools

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const DefaultPort = "18467"

func ParseAddress(addr string) (host string, port string, err error) {
	return ParseAddressWithDefault(addr, "0.0.0.0", DefaultPort)
}

func ParseAddressWithDefault(addr, defaultHost, defaultPort string) (host string, port string, err error) {
	// Si s est un nombre (port seul)
	if v, errPort := strconv.Atoi(addr); errPort == nil {
		if v < 1 || v > 65536 {
			return "", "", fmt.Errorf("port number %s is out of range", addr)
		}

		return defaultHost, addr, nil
	}

	// Essaie split host:port
	host, p, errSplit := net.SplitHostPort(addr)
	if errSplit == nil {
		if host == "" {
			host = defaultHost
		}
		return host, p, nil
	}

	// Pas de port, juste un host
	if strings.Contains(addr, ":") {
		// Par exemple IPv6 sans port comme [::1]
		// On enlève crochets si présents
		host = strings.Trim(addr, "[]")
		if host == "" {
			host = defaultHost
		}
		return host, defaultPort, nil
	}

	return addr, defaultPort, nil
}
