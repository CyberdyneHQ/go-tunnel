// utils.go - misc utilities used by HTTP and Socks proxies
//
// Author: Sudhi Herle <sudhi@herle.net>
//
// This software does not come with any express or implied
// warranty; it is provided "as is". No claim  is made to its
// suitability for any purpose.

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"time"
)

// Return true if the err represents a TCP PIPE or RESET error
func isReset(err error) bool {
	if oe, ok := err.(*net.OpError); ok {
		if se, ok := oe.Err.(*os.SyscallError); ok {
			if se.Err == syscall.EPIPE || se.Err == syscall.ECONNRESET {
				return true
			}
		}
	}
	return false
}

// Write all bytes in 'b' and return err
func WriteAll(fd io.Writer, b []byte) (int, error) {
	var z int
	n := len(b)
	for n > 0 {
		nw, err := fd.Write(b)
		if err != nil {
			return z, err
		}

		n -= nw
		z += nw
		b = b[nw:]
	}

	return z, nil
}

// Read upto len(b) bytes from fd
func ReadAll(fd io.Reader, b []byte) (int, error) {
	var z int
	n := len(b)
	for n > 0 {
		nr, err := fd.Read(b)
		if err != nil {
			return z, err
		}

		n -= nr
		z += nr
		b = b[nr:]
	}
	return z, nil
}

// Format a time duration
func format(t time.Duration) string {
	u0 := t.Nanoseconds() / 1000
	ma, mf := u0/1000, u0%1000

	if ma == 0 {
		return fmt.Sprintf("%3.3d us", mf)
	}

	return fmt.Sprintf("%d.%3.3d ms", ma, mf)
}

// Return true if the new connection 'conn' passes the ACL checks
// Return false otherwise
func AclOK(cfg *ListenConf, addr net.Addr) bool {
	var ip net.IP

	switch a := addr.(type) {
	case *net.TCPAddr:
		ip = a.IP

	case *net.UDPAddr:
		ip = a.IP

	default:
		return false // conservatively block
	}

	for _, n := range cfg.Deny {
		if n.Contains(ip) {
			return false
		}
	}

	if len(cfg.Allow) == 0 {
		return true
	}

	for _, n := range cfg.Allow {
		if n.Contains(ip) {
			return true
		}
	}

	return false
}
