package timer

import "time"

// NewTimer returns new started timer
func NewTimer() *Timer {
	c := &Timer{}
	c.Start()
	return c
}

type Timer struct {
	timerOn  bool
	started  int64
	stopped  int64
	duration int64
}

func (c *Timer) Start() {
	if !c.timerOn {
		c.started = time.Now().UnixNano()
		c.timerOn = true
	}
}

func (c *Timer) Stop() {
	if c.timerOn {
		c.duration += time.Now().UnixNano() - c.started
		c.timerOn = false
	}
}

func (c *Timer) Reset() {
	if c.timerOn {
		c.started = time.Now().UnixNano()
	}
	c.duration = 0
}

func (c *Timer) Seconds() float64 {
	if !c.timerOn {
		return float64(c.duration) / 1e9
	}
	return float64(c.duration+time.Now().UnixNano()-c.started) / 1e9
}
