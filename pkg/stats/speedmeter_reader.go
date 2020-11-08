package stats

import (
	"io"
	"time"
)

// SpeedmeterReader is the interface for io.Reader with Speedmeter
type SpeedmeterReader interface {
	io.Reader
	Speedmeter
}

type speedmeterReader struct {
	reader     io.Reader
	speedmeter *speedmeter
}

// NewSpeedmeterReader returns a io.Reader wrapper to calculate the transfer speed
func NewSpeedmeterReader(
	reader io.Reader,
	tick time.Duration,
) SpeedmeterReader {
	return &speedmeterReader{
		reader:     reader,
		speedmeter: newSpeedmeter(tick),
	}
}

var _ SpeedmeterReader = &speedmeterReader{}

func (r *speedmeterReader) Read(p []byte) (int, error) {
	n, err := r.reader.Read(p)

	r.speedmeter.update(int64(n))

	return n, err
}

// Current returns the current transfer speed
func (r *speedmeterReader) Current() int64 {

}
