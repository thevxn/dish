//go:build linux || darwin

package netrunner

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"runtime"
	"syscall"
	"time"

	"go.vxn.dev/dish/pkg/logger"
	"go.vxn.dev/dish/pkg/socket"
)

type ICMPType int

const (
	echoReply   ICMPType = 0
	echoRequest ICMPType = 8
)

const (
	ipStripHdr = 23
	testID     = 0x1234
	testSeq    = 0x0001
)

type icmpRunner struct {
	logger logger.Logger
}

// RunTest is used to test ICMP sockets. It sends an ICMP Echo Request to the given socket using
// non-privileged ICMP and verifies the reply. The test passes if the reply has the same payload
// as the request. Returns an error if the socket host cannot be resolved to an IPv4 address. If
// the host resolves to more than one address, only the first one is used.
func (runner *icmpRunner) RunTest(ctx context.Context, sock socket.Socket) socket.Result {
	runner.logger.Debugf("Resolving host '%s' to an IP address", sock.Host)

	addr, err := net.DefaultResolver.LookupIPAddr(ctx, sock.Host)
	if err != nil {
		return socket.Result{
			Socket: sock,
			Error:  fmt.Errorf("failed to resolve socket host: %w", err),
		}
	}

	ip := addr[0].IP.To4()
	if ip == nil {
		return socket.Result{Socket: sock, Error: errors.New("not a valid IPv4 address")}
	}

	sockAddr := &syscall.SockaddrInet4{Addr: [4]byte(ip)}

	// When using ICMP over DGRAM, Linux Kernel automatically sets (overwrites) and
	// validates the id, seq and checksum of each incoming and outgoing ICMP message.
	// This is largely non-documented in the linux man pages. The closest I found is:
	// - (Linux news) lwn.net/Articles/420800/
	// - (MacOS man) https://www.manpagez.com/man/4/icmp/
	// - (Third-party article) https://inc0x0.com/icmp-ip-packets-ping-manually-create-and-send-icmp-ip-packets/
	// "[...] most Linux systems use a unique identifier for every ping process, and sequence
	// number is an increasing number within that process. Windows uses a fixed identifier, which
	// varies between Windows versions, and a sequence number that is only reset at boot time."
	sysSocket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_ICMP)
	if err != nil {
		return socket.Result{
			Socket: sock,
			Error:  fmt.Errorf("failed to create a non-privileged icmp socket: %w", err),
		}
	}

	defer func() {
		if cerr := syscall.Close(sysSocket); cerr != nil {
			runner.logger.Errorf(
				"error closing ICMP socket (fd %d) for %s:%d: %v",
				sysSocket, sock.Host, sock.Port, cerr,
			)
		}
	}()

	if runtime.GOOS == "darwin" {
		if err := syscall.SetsockoptInt(sysSocket, syscall.IPPROTO_IP, ipStripHdr, 1); err != nil {
			return socket.Result{
				Socket: sock,
				Error:  fmt.Errorf("failed to set ip strip header: %w", err),
			}
		}
	}

	if d, ok := ctx.Deadline(); ok {
		// Set a socket receive timeout.
		t := syscall.NsecToTimeval(time.Until(d).Nanoseconds())
		if err := syscall.SetsockoptTimeval(sysSocket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &t); err != nil {
			return socket.Result{
				Socket: sock,
				Error: fmt.Errorf(
					"failed to set a timeout on a non-privileged icmp socket: %w",
					err,
				),
			}
		}
	}

	payload := []byte("ICMP echo")

	// ICMP Header size is 8 bytes.
	reqBuf := make([]byte, 8+len(payload))

	// ICMP Header.
	// ID, Seq and Checksum are filled in automatically by the kernel on linux machines, not on darwin ipv4
	reqBuf[0] = byte(echoRequest) // Type: Echo
	copy(reqBuf[8:], payload)

	// Set the ID, Seq and Checksum for the darwin based machines
	if runtime.GOOS == "darwin" {
		binary.BigEndian.PutUint16(reqBuf[4:6], testID)
		binary.BigEndian.PutUint16(reqBuf[6:8], testSeq)
		csum := checksum(reqBuf)
		reqBuf[2] ^= byte(csum)
		reqBuf[3] ^= byte(csum >> 8)
	}

	runner.logger.Debug("ICMP runner: send to " + ip.String())

	if err := syscall.Sendto(sysSocket, reqBuf, 0, sockAddr); err != nil {
		return socket.Result{
			Socket: sock,
			Error:  fmt.Errorf("failed to send an echo request: %w", err),
		}
	}

	// Maximum Transmission Unit (MTU) equals 1500 bytes.
	// Recvfrom before writing to the buffer, checks its length (not capacity).
	// If the length of the buffer is too small to fit the data then it's silently truncated.
	replyBuf := make([]byte, 1500)

	runner.logger.Debug("ICMP runner: recv from " + ip.String())

	n, _, err := syscall.Recvfrom(sysSocket, replyBuf, 0)
	if err != nil {
		return socket.Result{
			Socket: sock,
			Error:  fmt.Errorf("failed to receive a reply from a socket: %w", err),
		}
	}

	if n < 8 {
		return socket.Result{
			Socket: sock,
			Error:  fmt.Errorf("reply is too short: received %d bytes ", n),
		}
	}

	if replyBuf[0] != byte(echoReply) {
		return socket.Result{Socket: sock, Error: errors.New("received unexpected reply type")}
	}

	if !bytes.Equal(reqBuf[8:], replyBuf[8:n]) {
		return socket.Result{
			Socket: sock,
			Error:  errors.New("failed to validate echo reply: payloads are not equal"),
		}
	}

	return socket.Result{Socket: sock, Passed: true}
}

// checksum calculates the internet checksum for the given byte slice.
// This function was taken from the x/net/icmp package, which is not available in the standard library.
// https://godoc.org/golang.org/x/net/icmp
func checksum(b []byte) uint16 {
	csumcv := len(b) - 1 // checksum coverage
	s := uint32(0)
	for i := 0; i < csumcv; i += 2 {
		s += uint32(b[i+1])<<8 | uint32(b[i])
	}
	if csumcv&1 == 0 {
		s += uint32(b[csumcv])
	}
	s = s>>16 + s&0xffff
	s = s + s>>16
	return ^uint16(s)
}
