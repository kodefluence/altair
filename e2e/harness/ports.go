//go:build e2e

package harness

import (
	"fmt"
	"net"
)

// FreeTCP asks the kernel for an unused TCP port, closes the listener, and
// returns the number. There's a microsecond race between Close() and a
// subsequent Listen on the same port — acceptable for smoke tests; if it
// ever flakes, swap for a retry loop.
func FreeTCP() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, fmt.Errorf("allocate free port: %w", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
