package lib

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

var (
	ErrInvalidIPAddress      = errors.New("invalid IP address")
	ErrFailedToSplitHostPort = errors.New("failed to split host and port")
)

func SplitHostPort(addr string) (string, string, error) {
	if !strings.ContainsRune(addr, ':') {
		return addr, "", nil
	}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", "", fmt.Errorf("%w (addr=%q): %w", ErrFailedToSplitHostPort, addr, err)
	}

	return host, port, nil
}

func DetectLocalNetwork(requestAddr string) (bool, error) {
	var requestIP string

	requestAddrs := strings.SplitN(requestAddr, ",", 2) //nolint:mnd

	requestIP, _, err := SplitHostPort(requestAddrs[0])
	if err != nil {
		return false, err
	}

	requestIPNet := net.ParseIP(requestIP)
	if requestIPNet == nil {
		return false, fmt.Errorf("%w (request_ip=%q)", ErrInvalidIPAddress, requestIP)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false, err //nolint:wrapcheck
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		if !ipNet.Contains(requestIPNet) {
			continue
		}

		if requestIPNet.IsLoopback() {
			return true, nil
		}
	}

	return false, nil
}
