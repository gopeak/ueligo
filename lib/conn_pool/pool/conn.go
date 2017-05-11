package pool

import "net"

// poolConn is a wrapper around *net.TCPConn to modify the the behavior of
// *net.TCPConn's Close() method.
type poolConn struct {
	*net.TCPConn
	c *channelPool
}

// Close() puts the given connects back to the pool instead of closing it.
func (p poolConn) Close() error {
	return p.c.put(p.TCPConn)
}

// newConn wraps a standard *net.TCPConn to a poolConn *net.TCPConn.
func (c *channelPool) wrapConn(conn *net.TCPConn) *net.TCPConn {
	p := poolConn{c: c}
	p.TCPConn = conn
	return p.TCPConn
}
