package stats

import (
	"testing"
	"time"
)

func initSpeedmeter(t *testing.T, tick time.Duration) (*speedmeter, func(time.Time)) {
	t.Helper()

	meter := &speedmeter{}
	var now time.Time

	meter.now = func() time.Time {
		return now
	}
	meter.tick = 5 * time.Second

	return meter, func(t time.Time) {
		now = t
	}
}

func TestSpeedmeter(t *testing.T) {
	assert := func(t *testing.T, meter *speedmeter) func(expected int64) {
		return func(wanted int64) {
			t.Helper()

			if got := meter.Current(); got != wanted {
				t.Errorf("current should be %d, but got %d", wanted, got)
			}
		}
	}

	t.Run("success", func(t *testing.T) {
		tick := 5 * time.Second
		meter, setNow := initSpeedmeter(t, tick)
		assert := assert(t, meter)

		now := time.Date(2020, time.November, 9, 2, 7, 40, 0, time.FixedZone("Azia/Tokyo", 60*60*9))
		setNow(now)

		meter.update(2000)
		assert(0)

		now = now.Add(3 * time.Second)
		setNow(now)

		meter.update(1000)
		assert(0)

		now = now.Add(3 * time.Second)
		setNow(now)

		meter.update(3000)
		assert(1000)

		now = now.Add(5 * time.Second)
		setNow(now)

		meter.update(2500)
		assert(500)
	})
}
