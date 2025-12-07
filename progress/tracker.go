package progress

import (
	"fmt"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
)

type tracker struct {
	verbosity   uint8
	total       uint64
	done        uint64
	lastTime    time.Time
	startTime   time.Time
	updateEvery time.Duration
	mu          sync.Mutex
}

func newTracker(total uint64, verbosity uint8) *tracker {
	now := time.Now()
	return &tracker{
		verbosity:   verbosity,
		total:       total,
		startTime:   now,
		lastTime:    now,
		updateEvery: 1 * time.Second,
	}
}

func (p *tracker) add(n uint64) {
	if p.verbosity < 1 {
		return
	}
	now := time.Now()
	p.mu.Lock()
	defer p.mu.Unlock()

	p.done += n
	if now.Sub(p.lastTime) >= p.updateEvery {
		p.lastTime = now
		p.progress(now)
	}
}

func (p *tracker) progress(now time.Time) {
	totalElapsed := now.Sub(p.startTime).Seconds()
	if totalElapsed <= 0 {
		totalElapsed = 1
	}

	speed := float64(p.done) / totalElapsed

	var (
		total   = "?"
		percent = ""
		eta     = ""
	)

	if p.total > 0 {
		total = humanize.Bytes(p.total)
		per := float64(p.done) / float64(p.total) * 100
		percent = fmt.Sprintf(" (%.0f%%)", per)

		if speed > 0 {
			remaining := time.Duration(float64(p.total-p.done)/speed) * time.Second
			eta = fmt.Sprintf("| ETA %v", remaining.Round(time.Second).String())
		}
	}

	fmt.Printf("\r%-50s\r%s/%s%s | %s/s %s", "", humanize.Bytes(p.done), total, percent, humanize.Bytes(uint64(speed)), eta)
}

func (p *tracker) finish() {
	if p.verbosity < 1 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	totalTime := time.Since(p.startTime)
	avgSpeed := float64(p.done) / totalTime.Seconds()

	fmt.Printf("\r%-50v\r%s in %v | %s/s\n", "", humanize.Bytes(p.done), totalTime.Round(time.Second), humanize.Bytes(uint64(avgSpeed)))
}
