package timer_test

import (
	. "github.com/stojg/pants/timer"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	c := NewTimer()
	c.Stop()
	time.Sleep(time.Second / 10)
	actual := c.Seconds()
	if 0.1 < actual {
		t.Errorf("Clock should be less than 0.1 seconds, got %f", actual)
	}
	if 0.0 > actual {
		t.Errorf("Clock should be more than 0.0 seconds, got %f", actual)
	}
	c.Reset()
	actual = c.Seconds()
	if actual != 0.0 {
		t.Errorf("Time should be zero, got %f", actual)
	}
	c.Start()
	time.Sleep(time.Second / 10)
	actual = c.Seconds()
	if 0.1 > actual {
		t.Errorf("Clock should more than 0.1 seconds, got %f", actual)
	}
	if 0.2 < actual {
		t.Errorf("Clock should be less than 0.2 seconds, got %f", actual)
	}
	c.Reset()
	actual = c.Seconds()

	if 0.1 < actual {
		t.Errorf("Clock should less than 0.1 seconds, got %f", actual)
	}
	if 0.0 > actual {
		t.Errorf("Clock should be more than 0.0 seconds, got %f", actual)
	}
}

func BenchmarkGetTime(b *testing.B) {
	t := NewTimer()
	var now float64
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		now = t.Seconds()
	}
	_ = now
}
