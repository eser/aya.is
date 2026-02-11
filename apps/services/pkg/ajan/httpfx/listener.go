package httpfx

import (
	"context"
	"net"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

// ListenerConfig holds configuration for the high-performance TCP listener.
type ListenerConfig struct {
	// KeepAlivePeriod sets the TCP keep-alive period. Zero disables keep-alive.
	KeepAlivePeriod time.Duration

	// MaxConnections limits concurrent connections. Zero means unlimited.
	MaxConnections int

	// TCPNoDelay disables Nagle's algorithm (TCP_NODELAY) for lower latency.
	TCPNoDelay bool

	// KeepAlive enables TCP keep-alive on accepted connections.
	KeepAlive bool
}

// HighPerfListener wraps a net.Listener with connection limiting and optimized socket options.
type HighPerfListener struct {
	net.Listener

	config      *ListenerConfig
	connSem     chan struct{} // Semaphore for connection limiting
	activeConns int64         // Atomic counter for active connections
}

// NewHighPerfListener creates a high-performance TCP listener with optimized socket options.
// It configures TCP_NODELAY, SO_REUSEADDR, and optionally limits concurrent connections.
func NewHighPerfListener(
	ctx context.Context,
	addr string,
	config *ListenerConfig,
) (*HighPerfListener, error) {
	// Create listener config with socket options
	lc := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var sockErr error

			err := c.Control(func(fd uintptr) {
				// Enable SO_REUSEADDR for faster restarts
				sockErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEADDR, 1)
				if sockErr != nil {
					return
				}

				// Enable TCP_NODELAY for lower latency (disable Nagle's algorithm)
				if config.TCPNoDelay {
					sockErr = unix.SetsockoptInt(int(fd), unix.IPPROTO_TCP, unix.TCP_NODELAY, 1)
					if sockErr != nil {
						return
					}
				}

				// SO_REUSEPORT allows multiple listeners on same port (load balancing)
				// Note: Only enable if you're running multiple server instances
				// sockErr = unix.SetsockoptInt(int(fd), unix.SOL_SOCKET, unix.SO_REUSEPORT, 1)
			})
			if err != nil {
				return err
			}

			return sockErr
		},
		KeepAlive: config.KeepAlivePeriod,
	}

	// Create the listener
	ln, err := lc.Listen(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	hpl := &HighPerfListener{
		Listener:    ln,
		config:      config,
		activeConns: 0,
	}

	// Initialize connection semaphore if limiting is enabled
	if config.MaxConnections > 0 {
		hpl.connSem = make(chan struct{}, config.MaxConnections)
	}

	return hpl, nil
}

// Accept waits for and returns the next connection with optimized settings.
// If MaxConnections is set, it blocks when the limit is reached.
func (l *HighPerfListener) Accept() (net.Conn, error) {
	// Acquire semaphore slot if connection limiting is enabled
	if l.connSem != nil {
		l.connSem <- struct{}{}
	}

	conn, err := l.Listener.Accept()
	if err != nil {
		// Release semaphore on error
		if l.connSem != nil {
			<-l.connSem
		}

		return nil, err
	}

	// Apply TCP optimizations to the connection
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		// Enable keep-alive
		if l.config.KeepAlive {
			_ = tcpConn.SetKeepAlive(true)

			if l.config.KeepAlivePeriod > 0 {
				_ = tcpConn.SetKeepAlivePeriod(l.config.KeepAlivePeriod)
			}
		}

		// TCP_NODELAY is set at listener level, but we can also set it per-connection
		if l.config.TCPNoDelay {
			_ = tcpConn.SetNoDelay(true)
		}
	}

	atomic.AddInt64(&l.activeConns, 1)

	// Wrap connection to track when it's closed
	return &trackedConn{
		Conn:     conn,
		listener: l,
	}, nil
}

// ActiveConnections returns the current number of active connections.
func (l *HighPerfListener) ActiveConnections() int64 {
	return atomic.LoadInt64(&l.activeConns)
}

// trackedConn wraps a net.Conn to track connection lifecycle.
type trackedConn struct {
	net.Conn

	listener *HighPerfListener
	closed   int32 // Atomic flag to prevent double-close
}

// Close closes the connection and releases resources.
func (c *trackedConn) Close() error {
	// Prevent double-close
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}

	// Decrement active connections
	atomic.AddInt64(&c.listener.activeConns, -1)

	// Release semaphore slot
	if c.listener.connSem != nil {
		<-c.listener.connSem
	}

	return c.Conn.Close()
}
