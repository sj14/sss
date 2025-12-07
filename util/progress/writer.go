package progress

import "io"

type Writer struct {
	writer  io.WriterAt
	tracker *tracker
}

func NewWriter(w io.WriterAt, total uint64, verbosity uint8, key string) *Writer {
	return &Writer{
		writer:  w,
		tracker: newTracker(total, verbosity, key),
	}
}

func (p *Writer) WriteAt(b []byte, off int64) (int, error) {
	n, err := p.writer.WriteAt(b, off)
	if n > 0 {
		p.tracker.add(uint64(n))
	}
	return n, err
}

func (p *Writer) Finish() {
	p.tracker.finish()
}
