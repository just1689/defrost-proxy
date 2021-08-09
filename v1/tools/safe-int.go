package tools

import "sync"

type SafeInt struct {
	v int
	l sync.Mutex
}

func (s *SafeInt) Set(v int) {
	s.l.Lock()
	defer s.l.Unlock()
	s.v = v
}

func (s *SafeInt) Get() int {
	s.l.Lock()
	defer s.l.Unlock()
	return s.v
}
