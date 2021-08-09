package tools

import (
	"time"
)

func NewCooldown(seconds int, freeze func()) *cooldown {
	r := &cooldown{
		seconds: seconds,
		ping:    make(chan interface{}),
		freeze:  freeze,
	}
	go r.run()
	return r
}

type cooldown struct {
	seconds int
	ping    chan interface{}
	freeze  func()
}

func (c *cooldown) Ping() {
	c.ping <- true
}

func (c *cooldown) run() {
	for {
		select {
		case <-time.After(time.Duration(c.seconds)):
			c.freeze()
		case <-c.ping:
			//
		}
	}
}
