package model

import (
	"sync"
)

type Connection struct {
	Ip     string
	Name   string
	Active bool
	Time   string
}

type Connections struct {
	sync.RWMutex
	Connections map[string]*Connection
}

func (c *Connections) Add(con Connection) {
	c.Lock()
	defer c.Unlock()
	if c.Connections == nil {
		c.Connections = make(map[string]*Connection)
	}
	c.Connections[con.Ip] = &con
}

func (c *Connections) Remove(con Connection) {
	c.Lock()
	defer c.Unlock()
	delete(c.Connections, con.Ip)
}

func (c *Connections) Disable(con Connection) {
	c.Lock()
	defer c.Unlock()
	c.Connections[con.Ip].Active = false
}

func (c *Connections) Enable(con Connection) {
	c.Lock()
	defer c.Unlock()
	c.Connections[con.Ip].Active = true
}
