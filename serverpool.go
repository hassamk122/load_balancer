package main

import "sync/atomic"

type ServerPool struct {
	backends []*Backend
	current  uint64
}

// increment atomically. taking mod by length for cycle
func (s *ServerPool) NextIndex() int {
	return int(atomic.AddUint64(&s.current, uint64(1)) % uint64(len(s.backends)))
}

// returns next active peer to take a conn
// loops entire backends to find a alive backend
// if backend is alive, use it and mark it
func (s *ServerPool) GetNextPeer() *Backend {
	next := s.NextIndex()
	l := len(s.backends) + next
	for i := next; i < l; i++ {
		idx := i % len(s.backends)

		if s.backends[idx].isAlive() {
			if i != next {
				atomic.StoreUint64(&s.current, uint64(idx))
			}
			return s.backends[idx]
		}
	}
	return nil
}
