package model

import "sync"

var Counter = counter{}

type counter struct {
	sync.Mutex
	v int
}

func (c *counter) Add() {
	c.Lock()
	c.v++
	c.Unlock()
}

func (c *counter) Sub() {
	c.Lock()
	c.v--
	c.Unlock()
}

func (c *counter) Get() int {
	c.Lock()
	defer c.Unlock()
	return c.v
}
