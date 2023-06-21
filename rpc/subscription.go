package rpc

import (
	"context"
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"strings"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// A Subscription is created by a notifier and tied to that notifier. The client can use
// this subscription to wait for an unsubscribe request for the client, see Err().
type Subscription struct {
	ID     ID
	method string
	err    chan error // closed on unsubscribe
}

// Err returns a channel that is closed when the client send an unsubscribe request.
func (s *Subscription) Err() <-chan error {
	return s.err
}

var globalGen = randomIDGenerator()

// ID defines a pseudo random number that is used to identify RPC subscriptions.
type ID string

// NewID returns a new, random ID.
func NewID() ID {
	return globalGen()
}

// randomIDGenerator returns a function generates a random IDs.
func randomIDGenerator() func() ID {
	var buf = make([]byte, 8)
	var seed int64
	if _, err := crand.Read(buf); err == nil {
		seed = int64(binary.BigEndian.Uint64(buf))
	} else {
		seed = int64(time.Now().Nanosecond())
	}

	var (
		mu  sync.Mutex
		rng = rand.New(rand.NewSource(seed))
	)
	return func() ID {
		mu.Lock()
		defer mu.Unlock()
		id := make([]byte, 16)
		rng.Read(id)
		return encodeID(id)
	}
}

func encodeID(b []byte) ID {
	id := hex.EncodeToString(b)
	id = strings.TrimLeft(id, "0")
	if id == "" {
		id = "0" // ID's are RPC quantities, no leading zero's and 0 is 0x0.
	}
	return ID("0x" + id)
}

type notifierKey struct{}

// PutNotifierKey is a helper function to put Notifier in the context so that subscription handlers can use it
func PutNotifierKey(ctx context.Context, n *Notifier) context.Context {
	return context.WithValue(ctx, notifierKey{}, n)
}

// NotifierFromContext returns the Notifier value stored in ctx, if any.
func NotifierFromContext(ctx context.Context) (*Notifier, bool) {
	n, ok := ctx.Value(notifierKey{}).(*Notifier)
	return n, ok
}

// Notifier is tied to a RPC connection that supports subscriptions.
// Server callbacks use the notifier to send notifications.
type Notifier struct {
	wsCtx  *WebSocketContext
	h      *handler
	method string

	mu  sync.Mutex
	sub *Subscription
}

// CreateSubscription returns a new subscription that is coupled to the
// RPC connection. By default subscriptions are inactive and notifications
// are dropped until the subscription is marked as active. This is done
// by the RPC server after the subscription ID is send to the client.
func (n *Notifier) CreateSubscription() *Subscription {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.sub != nil {
		panic("can't create multiple subscriptions with Notifier")
	}
	n.sub = &Subscription{ID: globalGen(), method: n.method, err: make(chan error, 1)}
	n.wsCtx.subscriptions[n.sub.ID] = n.sub
	return n.sub
}

// Notify sends a notification to the client with the given data as payload.
func (n *Notifier) Notify(id ID, data interface{}) error {
	enc, err := jsoniter.Marshal(data)
	if err != nil {
		return err
	}

	if n.sub == nil {
		panic("can't Notify before subscription is created")
	} else if n.sub.ID != id {
		panic("Notify with wrong ID")
	}

	return n.send(n.sub, enc)
}

// send generates the response and writes is to the websocket connection's output channel
func (n *Notifier) send(sub *Subscription, data jsoniter.RawMessage) error {
	resp := createEventResponse([]byte(sub.ID), data)
	if resp != nil {
		n.wsCtx.output <- resp
	}

	return nil
}
