package ratelimiter

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

type LimitReader struct {
	ctx     context.Context
	reader  io.Reader
	limiter *rate.Limiter
}

func NewReader(ctx context.Context, r io.Reader, l *rate.Limiter) *LimitReader {
	return &LimitReader{
		ctx:     ctx,
		reader:  r,
		limiter: l,
	}
}

func (lr *LimitReader) Read(p []byte) (int, error) {
	// WaitN errors immediately if asked to wait for more tokens than the
	// limiter's burst allows, so cap the read to the burst size instead of
	// passing the caller's (potentially larger) buffer straight through.
	if burst := lr.limiter.Burst(); burst > 0 && len(p) > burst {
		p = p[:burst]
	}

	if err := lr.limiter.WaitN(lr.ctx, len(p)); err != nil {
		return 0, err
	}
	return lr.reader.Read(p)
}
