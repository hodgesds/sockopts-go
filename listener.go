package sockopts

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

func SockoptListener(network, address string, sockopts ...int) (net.Listener, error) {
	var err error
	var fd int
	var addr syscall.Sockaddr

	switch network {
	case "tcp", "tcp4":
		fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
		if err != nil {
			return nil, fmt.Errorf("could not open socket")
		}
		netAddr, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, fmt.Errorf("could not open socket")
		}
		var ipAddr [4]byte
		copy(ipAddr[:], netAddr.IP)
		addr = &syscall.SockaddrInet4{
			Port: netAddr.Port,
			Addr: ipAddr,
		}
	case "tcp6":
		fd, err = syscall.Socket(syscall.AF_INET6, syscall.SOCK_STREAM, 0)
		if err != nil {
			return nil, fmt.Errorf("could not open socket")
		}
		netAddr, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, fmt.Errorf("could not open socket")
		}
		ipAddr := [16]byte{}
		copy(ipAddr[:], netAddr.IP)
		addr = &syscall.SockaddrInet6{
			Port: netAddr.Port,
			Addr: ipAddr,
		}
	case "unix":
		fd, err = syscall.Socket(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
		addr = &syscall.SockaddrUnix{
			Name: address,
		}
	case "udp", "udp4":
		fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
		if err != nil {
			return nil, fmt.Errorf("could not open socket")
		}
		netAddr, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, fmt.Errorf("could not open socket")
		}
		ipAddr := [4]byte{}
		copy(ipAddr[:], netAddr.IP)
		addr = &syscall.SockaddrInet4{
			Port: netAddr.Port,
			Addr: ipAddr,
		}
	case "udp6":
		fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
		if err != nil {
			return nil, fmt.Errorf("could not open socket")
		}
		netAddr, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, fmt.Errorf("could not open socket")
		}
		ipAddr := [16]byte{}
		copy(ipAddr[:], netAddr.IP)
		addr = &syscall.SockaddrInet6{
			Port: netAddr.Port,
			Addr: ipAddr,
		}
	default:
		return nil, fmt.Errorf("unknown network family: %s", network)
	}
	if err != nil {
		syscall.Close(fd)
		return nil, err
	}

	for _, sockopt := range sockopts {
		if sockopt == SO_REUSEPORT {
			if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, sockopt, 1); err != nil {
				syscall.Close(fd)
				return nil, err
			}
		} else if socktop == TCP_FASTOPEN {
			if err := allowFastOpen(); err != nil {
				return nil, err
			}
			if err := syscall.SetsockoptInt(fd, syscall.SOL_TCP, sockopt, 1); err != nil {
				syscall.Close(fd)
				return nil, err
			}
		}
	}

	if err := syscall.Bind(fd, addr); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	if err := syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		syscall.Close(fd)
		return nil, err
	}

	f := os.NewFile(uintptr(fd), "l")

	return net.FileListener(f) // or net.FileConn, net.FilePacketConn
}

/*
 The changes required to a server in order to support TFO are minimal, and are highlighted in the code template below.

    sfd = socket(AF_INET, SOCK_STREAM, 0);   // Create socket

    bind(sfd, ...);                          // Bind to well known address

    int qlen = 5;                            // Value to be chosen by application
    setsockopt(sfd, SOL_TCP, TCP_FASTOPEN, &qlen, sizeof(qlen));

    listen(sfd, ...);                        // Mark socket to receive connections

    cfd = accept(sfd, NULL, 0);              // Accept connection on new socket

    // read and write data on connected socket cfd

    close(cfd);
*/
