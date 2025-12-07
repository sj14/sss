package ratelimiter

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

type LimitReader struct {
	r io.Reader
	l *rate.Limiter
}

func NewReader(r io.Reader, l *rate.Limiter) *LimitReader {
	return &LimitReader{
		r: r,
		l: l,
	}
}

func (rr *LimitReader) Read(p []byte) (int, error) {
	_ = rr.l.WaitN(context.Background(), len(p))
	return rr.r.Read(p)
}
