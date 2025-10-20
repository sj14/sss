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
	percent := float64(p.done) / float64(p.total) * 100

	var remaining time.Duration
	if speed > 0 {
		remaining = time.Duration(float64(p.total-p.done)/speed) * time.Second
	}

	fmt.Printf("\r%-50s\r%.0f%% of %s | %s/s | ETA %v", "", percent, humanize.Bytes(p.total), humanize.Bytes(uint64(speed)), remaining.Round(time.Second))
}

func (p *tracker) finish() {
	if p.verbosity < 1 {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	totalTime := time.Since(p.startTime)
	avgSpeed := float64(p.done) / totalTime.Seconds()

	fmt.Printf("\r%-50v\r%s in %v | %s/s\n", "", humanize.Bytes(p.total), totalTime.Round(time.Second), humanize.Bytes(uint64(avgSpeed)))
}
