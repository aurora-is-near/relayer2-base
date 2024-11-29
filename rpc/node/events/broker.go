package events

import (
	"bytes"
	"strings"
	"time"

	"github.com/aurora-is-near/relayer2-base/broker"
	"github.com/aurora-is-near/relayer2-base/log"
	"github.com/aurora-is-near/relayer2-base/rpc"
	"github.com/aurora-is-near/relayer2-base/types/common"
	"github.com/aurora-is-near/relayer2-base/types/event"
	"github.com/aurora-is-near/relayer2-base/types/request"
)

type Type byte

const (
	// LogsChSize is the size of channel listening to Logs types.
	LogsChSize = 5
	// NewHeadsChSize is the size of channel listening to NewHeads types.
	NewHeadsChSize = 5
)

const (
	// UnknownSubscription indicates an unknown subscription type
	UnknownSubscription Type = iota
	// NewHeadsSubscription tracks the newly added block responses
	NewHeadsSubscription
	// LogsSubscription queries for new or removed (chain reorg) logs
	LogsSubscription
)

type EventSubscription struct {
	id         rpc.ID
	typ        Type
	created    time.Time
	logOpts    request.LogSubscriptionOptions
	newHeadsCh chan event.Block
	logsCh     chan event.Logs
}

// GetId returns identifier of the EventSubscription which implements broker.Subscription
func (es *EventSubscription) GetId() broker.SubID {
	return broker.SubID(es.id)
}

// GetLogsSubOpts returns log filters of the EventSubscription which implements broker.Subscription
func (es *EventSubscription) GetLogsSubOpts() request.LogSubscriptionOptions {
	return es.logOpts
}

// GetNewHeadsCh returns the newHeads types channel of the EventSubscription which implements broker.Subscription
func (es *EventSubscription) GetNewHeadsCh() chan event.Block {
	return es.newHeadsCh
}

// GetLogsCh returns the logs types channel of the EventSubscription which implements broker.Subscription
func (es *EventSubscription) GetLogsCh() chan event.Logs {
	return es.logsCh
}

// EventBroker offers support to manage types subscriptions and broadcast the incoming events to
// subscribed objects.
type EventBroker struct {
	l                 *log.Logger
	stopCh            chan bool
	publishNewHeadsCh chan event.Block
	publishLogsCh     chan event.Logs
	subNewHeadsCh     chan broker.Subscription
	subLogsCh         chan broker.Subscription
	unsubNewHeadsCh   chan broker.Subscription
	unsubLogsCh       chan broker.Subscription
	DebugInfo         chan int
}

// NewEventBroker creates a new EventBroker object
func NewEventBroker() *EventBroker {
	return &EventBroker{
		l:                 log.Log(),
		stopCh:            make(chan bool, 1),
		publishNewHeadsCh: make(chan event.Block, NewHeadsChSize),
		publishLogsCh:     make(chan event.Logs, LogsChSize),
		subNewHeadsCh:     make(chan broker.Subscription),
		subLogsCh:         make(chan broker.Subscription),
		unsubNewHeadsCh:   make(chan broker.Subscription),
		unsubLogsCh:       make(chan broker.Subscription),
		DebugInfo:         make(chan int),
	}
}

// SubscribeNewHeads creates a new subscription and signals the
// EventBroker subscription channel to handle the subscription map
func (eb *EventBroker) SubscribeNewHeads(ch chan event.Block) broker.Subscription {
	sub := &EventSubscription{
		id:         rpc.NewID(),
		typ:        NewHeadsSubscription,
		created:    time.Now(),
		newHeadsCh: ch,
		logsCh:     make(chan event.Logs),
	}
	eb.subNewHeadsCh <- sub
	eb.l.Debug().Msgf("new subscription request to New Heads with Id: [%s]", sub.id)
	return sub
}

// SubscribeLogs creates a new subscription and signals the
// EventBroker subscription channel to handle the subscription map
func (eb *EventBroker) SubscribeLogs(opts request.LogSubscriptionOptions, ch chan event.Logs) broker.Subscription {
	sub := &EventSubscription{
		id:         rpc.NewID(),
		typ:        LogsSubscription,
		created:    time.Now(),
		logOpts:    opts,
		newHeadsCh: make(chan event.Block),
		logsCh:     ch,
	}
	eb.subLogsCh <- sub
	eb.l.Debug().Msgf("new subscription request to Logs with Id: [%s]", sub.id)
	return sub
}

// UnsubscribeFromNewHeads signals the EventBroker's related channel
// to delete the subscription
func (eb *EventBroker) UnsubscribeFromNewHeads(sub broker.Subscription) {
	eb.unsubNewHeadsCh <- sub
	eb.l.Debug().Msgf("unsubscription request to New Heads with Id: [%s]", sub.GetId())
}

// UnsubscribeFromLogs signals the EventBroker's related channel
// to delete the subscription
func (eb *EventBroker) UnsubscribeFromLogs(sub broker.Subscription) {
	eb.unsubLogsCh <- sub
	eb.l.Debug().Msgf("unsubscription request to Logs with Id: [%s]", sub.GetId())
}

// Start main loop of the EventBroker that receives and distributes the events.
func (eb *EventBroker) Start() {
	subsNewHeads := map[broker.SubID]broker.Subscription{}
	subsLogs := map[broker.SubID]broker.Subscription{}
	for {
		select {
		case <-eb.stopCh:
			return
		// this case is only for testing purposes
		case req := <-eb.DebugInfo:
			switch req {
			case -1:
				eb.DebugInfo <- len(subsNewHeads)
			case -2:
				eb.DebugInfo <- len(subsLogs)
			}
		case sub := <-eb.subNewHeadsCh:
			subsNewHeads[sub.GetId()] = sub
		case sub := <-eb.subLogsCh:
			subsLogs[sub.GetId()] = sub
		case sub := <-eb.unsubNewHeadsCh:
			delete(subsNewHeads, sub.GetId())
		case sub := <-eb.unsubLogsCh:
			delete(subsLogs, sub.GetId())
		case msg := <-eb.publishNewHeadsCh:
			for _, v := range subsNewHeads {
				// v.newHeadsCh is buffered, use non-blocking send to protect the broker:
				// timeout preferred instead of default to be able to tolerate slight delays
				select {
				case v.GetNewHeadsCh() <- msg:
				case <-time.After(10 * time.Millisecond):
					eb.l.Warn().Msg("Publishing to New Heads channel fall into DEFAULT!")
				}
			}
		case logs := <-eb.publishLogsCh:
			if len(logs) > 0 {
				for _, s := range subsLogs {
					matchedLogs := filterLogs(logs, s.GetLogsSubOpts())
					if len(matchedLogs) > 0 {
						// v.logsCh is buffered, use non-blocking send to protect the broker:
						// timeout preferred instead of default to be able to tolerate slight delays
						select {
						case s.GetLogsCh() <- matchedLogs:
						case <-time.After(10 * time.Millisecond):
							eb.l.Warn().Msg("Publishing to Logs channel fall into DEFAULT!")
						}
					}
				}
			}
		}
	}
}

// PublishNewHeads provides publish API for new block head types. Implements broker.Broker interface
func (eb *EventBroker) PublishNewHeads(b event.Block) {
	eb.publishNewHeadsCh <- b
}

// PublishLogs provides publish API for logs types. Implements broker.Broker interface
func (eb *EventBroker) PublishLogs(l event.Logs) {
	eb.publishLogsCh <- l
}

// filterLogs creates a slice of Log Response matching the given criteria.
func filterLogs(logs event.Logs, opts request.LogSubscriptionOptions) event.Logs {
	var ret event.Logs
Logs:
	for _, log := range logs {
		if len(opts.Address) > 0 && !includes(opts.Address, log.Address.Hex()) {
			continue
		}
		for i, sub := range opts.Topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				// empty rule set == wildcard. Otherwise, check if topic index of opts fits in the number of topics in received log and if it fits check topics for equality
				if len(topic.Content) == 0 || (i < len(log.Topics) && bytes.Equal(log.Topics[i].Bytes(), topic.Bytes())) {
					match = true
					break
				}
			}
			if !match {
				continue Logs
			}
		}
		ret = append(ret, log)
	}
	return ret
}

func includes(addresses []common.Address, address string) bool {
	for _, addr := range addresses {
		if strings.EqualFold(addr.Hex(), address) {
			return true
		}
	}
	return false
}
