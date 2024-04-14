package pipe

import (
	"fmt"
	"hash"
	"io"
	"sync"
	"time"
)

// TODO: custom buffer sizes
const gauge_cache_size = 8 * 1024 * 1024

type ReaderGauge struct {
	m     sync.Mutex
	count uint64
}

func (r *ReaderGauge) Count() uint64 {
	r.m.Lock()
	defer r.m.Unlock()
	return r.count
}

func (r *ReaderGauge) Measure(in io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		b := [gauge_cache_size]byte{}
		pos := 0
		for {
			n, err := in.Read(b[pos:])
			if err != nil {
				pw.CloseWithError(err)
				break
			}
			pos += n
			r.m.Lock()
			r.count += uint64(n)
			r.m.Unlock()
			m, err := pw.Write(b[:pos])
			pos -= m
			if err != nil {
				pw.CloseWithError(fmt.Errorf("could not write: %w", err))
				break
			}
		}
	}()
	return pr
}

type TimeRangedReaderGauge struct {
	m     sync.Mutex
	count uint64
	first time.Time
	last  time.Time
}

func (r *TimeRangedReaderGauge) Count() uint64 {
	r.m.Lock()
	defer r.m.Unlock()
	return r.count
}

func (r *TimeRangedReaderGauge) TimeRange() (time.Time, time.Time) {
	r.m.Lock()
	defer r.m.Unlock()
	return r.first, r.last
}

func (r *TimeRangedReaderGauge) Measure(in io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		b := [gauge_cache_size]byte{}
		pos := 0
		for {
			n, err := in.Read(b[pos:])
			if err != nil {
				pw.CloseWithError(err)
				break
			}
			pos += n
			r.m.Lock()
			if r.count == 0 {
				r.first = time.Now()
			}
			r.count += uint64(n)
			r.last = time.Now()
			r.m.Unlock()
			m, err := pw.Write(b[:pos])
			pos -= m
			if err != nil {
				pw.CloseWithError(fmt.Errorf("could not write: %w", err))
				break
			}
		}
	}()
	return pr
}

type Hash32Gauge struct {
	m sync.Mutex
	h hash.Hash32
}

func NewHash32Gauge(h hash.Hash32) *Hash32Gauge {
	return &Hash32Gauge{h: h}
}

func (r *Hash32Gauge) Sum32() uint32 {
	r.m.Lock()
	defer r.m.Unlock()
	return r.h.Sum32()
}

func (r *Hash32Gauge) Measure(in io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		b := [gauge_cache_size]byte{}
		pos := 0
		for {
			n, err := in.Read(b[pos:])
			if err != nil {
				pw.CloseWithError(err)
				break
			}
			pos += n
			r.m.Lock()
			r.h.Write(b[:pos])
			r.m.Unlock()
			m, err := pw.Write(b[:pos])
			pos -= m
			if err != nil {
				pw.CloseWithError(fmt.Errorf("could not write: %w", err))
				break
			}
		}
	}()
	return pr
}
