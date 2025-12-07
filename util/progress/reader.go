package progress

import "io"

type Reader struct {
	reader  io.Reader
	tracker *tracker
}

func NewReader(r io.Reader, total uint64, verbosity uint8, key string) *Reader {
	return &Reader{
		reader:  r,
		tracker: newTracker(total, verbosity, key),
	}
}

func (r *Reader) Read(b []byte) (int, error) {
	n, err := r.reader.Read(b)
	if n > 0 {
		r.tracker.add(uint64(n))
	}
	return n, err
}

func (r *Reader) Finish() {
	r.tracker.finish()
}
