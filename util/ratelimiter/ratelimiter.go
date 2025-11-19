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
	err := lr.limiter.WaitN(lr.ctx, len(p))
	if err != nil {
		return 0, err
	}
	return lr.reader.Read(p)
}
