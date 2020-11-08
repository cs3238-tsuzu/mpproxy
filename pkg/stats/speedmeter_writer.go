package stats

import (
	"io"
	"time"
)

// SpeedmeterWriter is the interface for io.Writer with Speedmeter
type SpeedmeterWriter interface {
	io.Writer
	Speedmeter
}

type speedmeterWriter struct {
	Writer     io.Writer
	speedmeter *speedmeter
}

// NewSpeedmeterWriter returns a io.Writer wrapper to calculate the transfer speed
func NewSpeedmeterWriter(
	Writer io.Writer,
	tick time.Duration,
) io.Writer {
	return &speedmeterWriter{
		Writer:     Writer,
		speedmeter: newSpeedmeter(tick),
	}
}

var _ SpeedmeterWriter = &speedmeterWriter{}

func (r *speedmeterWriter) Write(p []byte) (int, error) {
	n, err := r.Writer.Write(p)

	r.speedmeter.update(int64(n))

	return n, err
}

// Current returns the current transfer speed
func (r *speedmeterWriter) Current() int64 {
	return r.speedmeter.Current()
}
