//go:build linux || darwin

package netrunner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

	"go.vxn.dev/dish/pkg/socket"
)

type icmpRunner struct {
	verbose bool
}

// RunTest is used to test ICMP sockets. It sends an ICMP Echo Request to the given socket using
// non-privileged ICMP and verifies the reply. The test passes if the reply has the same payload
// as the request. Returns an error if the socket host cannot resolve to an IP address. If the host
// resolves to more than one address, only the first one is tested.
func (runner icmpRunner) RunTest(ctx context.Context, sock socket.Socket) socket.Result {
	if runner.verbose {
		log.Printf("Resolving host '%s' to an IP address", sock.Host)
	}

	addr, err := net.DefaultResolver.LookupIP(ctx, "ip", sock.Host)
	if err != nil {
		return socket.Result{Socket: sock, Error: fmt.Errorf("failed to resolve socket host: %w", err)}
	}

	ip := addr[0]

	sockAddr := &syscall.SockaddrInet4{Addr: [4]byte(ip)}

	if runner.verbose {
		log.Println("ICMP runner: send to " + ip.String())
	}

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
		return socket.Result{Socket: sock, Error: fmt.Errorf("failed to create a non-privileged icmp socket: %w", err)}
	}
	defer syscall.Close(sysSocket)

	if d, ok := ctx.Deadline(); ok {
		// Set a socket receive timeout.
		t := syscall.NsecToTimeval(time.Until(d).Nanoseconds())
		if err := syscall.SetsockoptTimeval(sysSocket, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &t); err != nil {
			return socket.Result{Socket: sock, Error: fmt.Errorf("failed to set a timeout on a non-privileged icmp socket: %w", err)}
		}
	}

	payload := []byte("ICMP echo")

	// ICMP Header size is 8 bytes.
	reqBuf := make([]byte, 8+len(payload))

	// ICMP Header.
	// ID, Seq and Checksum are filled in automatically by the kernel.
	reqBuf[0] = 8 // Type: Echo

	copy(reqBuf[8:], payload)

	if err := syscall.Sendto(sysSocket, reqBuf, 0, sockAddr); err != nil {
		return socket.Result{Socket: sock, Error: fmt.Errorf("failed to send an echo request: %w", err)}
	}

	// Maximum Transmission Unit (MTU) equals 1500 bytes.
	// Recvfrom before writing to the buffer, checks its length (not capacity).
	// If the length of the buffer is too small to fit the data then it's silently truncated.
	replyBuf := make([]byte, 1500)

	n, _, err := syscall.Recvfrom(sysSocket, replyBuf, 0)
	if err != nil {
		return socket.Result{Socket: sock, Error: fmt.Errorf("failed to receive a reply from a socket: %w", err)}
	}

	if n < 8 {
		return socket.Result{Socket: sock, Error: fmt.Errorf("reply is too short: received %d bytes ", n)}
	}

	if replyBuf[0] != 0 {
		return socket.Result{Socket: sock, Error: errors.New("received unexpected reply type")}
	}

	if !bytes.Equal(reqBuf[8:], replyBuf[8:n]) {
		return socket.Result{Socket: sock, Error: errors.New("failed to validate echo reply: payloads are not equal")}
	}

	return socket.Result{Socket: sock, Passed: true}
}
