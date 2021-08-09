package model

import "sync"

var ScaledUp = &TSFlag{}

type TSFlag struct {
	sync.Mutex
	v bool
}

func (t *TSFlag) Set(v bool) {
	t.Lock()
	t.v = v
	t.Unlock()
}

func (t *TSFlag) Get() bool {
	t.Lock()
	defer t.Unlock()
	return t.v
}
