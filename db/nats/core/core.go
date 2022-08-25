package core

import (
	"github.com/nats-io/nats.go"
	"sync"
)

var lock = &sync.Mutex{}
var nc *nats.Conn

func Open(options nats.Options) (*nats.Conn, error) {
	var err error
	if nc == nil {
		lock.Lock()
		defer lock.Unlock()
		if nc == nil {
			nc, err = options.Connect()
		}
	}
	return nc, err
}

func Close() error {
	return nc.Drain()
}
