package protocol

import (
	"sync/atomic"
)

var connId int64

type Conn struct {
	id    int64
	state int
	conn  ReadWriteCloser
}

func Wrapper(conn ReadWriteCloser) *Conn {
	return &Conn{
		id:   atomic.AddInt64(&connId, 1),
		conn: conn,
	}
}

func (c *Conn) ID() int64 {
	return c.id
}

func (c *Conn) IP() string {
	return c.conn.IP()
}

func (c *Conn) Close() error {
	c.state = 1
	return c.conn.Close()
}

func (c *Conn) State() int {
	return c.state
}

func (c *Conn) Accept(apply func(msg Packet, c *Conn)) error {
	for {
		packet, err := c.conn.Read()
		if err != nil {
			return err
		}
		apply(*packet, c)
	}
}

func (c *Conn) Write(packet Packet) error {
	return c.conn.Write(packet)
}

func (c *Conn) Read() (*Packet, error) {
	return c.conn.Read()
}
