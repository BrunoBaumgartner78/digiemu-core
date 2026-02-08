package memory

import "time"

type RealClock struct{}

func (RealClock) NowUnix() int64 {
	return time.Now().UTC().Unix()
}

type FakeClock struct {
	Now int64
}

func (c FakeClock) NowUnix() int64 { return c.Now }
