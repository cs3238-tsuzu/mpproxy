package stats

import (
	"net"
	"time"
)

// SpeedmeterConn is the interface for io.Reader with Speedmeter
type SpeedmeterConn interface {
	net.Conn
	ReadCurrent() int64
	WriteCurrent() int64
}

type speedmeterConn struct {
	net.Conn
	reader SpeedmeterReader
	writer SpeedmeterWriter
}

// NewSpeedmeterConn returns a net.Conn wrapper to calculate the transfer speed
func NewSpeedmeterConn(
	conn net.Conn,
	tick time.Duration,
) SpeedmeterConn {
	return &speedmeterConn{
		Conn:   conn,
		reader: NewSpeedmeterReader(conn, tick),
		writer: NewSpeedmeterWriter(conn, tick),
	}
}

func (c *speedmeterConn) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func (c *speedmeterConn) Write(p []byte) (int, error) {
	return c.writer.Write(p)
}

func (c *speedmeterConn) ReadCurrent() int64 {
	return c.reader.Current()
}

func (c *speedmeterConn) WriteCurrent() int64 {
	return c.writer.Current()
}
