package eventbroker

import (
	"aurora-relayer-go-common/log"
	"aurora-relayer-go-common/utils"
	"bytes"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
)

type Type byte

const (
	// logsChSize is the size of channel listening to Logs event.
	LogsChSize = 5
	// newHeadsChSize is the size of channel listening to NewHeads event.
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

type Subscription struct {
	id         rpc.ID
	typ        Type
	created    time.Time
	logOpts    utils.LogSubscriptionOptions
	newHeadsCh chan *utils.BlockResponse
	logsCh     chan []*utils.LogResponse
}

// EventBroker offers support to manage event subscriptions and broadcast the incoming events to
// subscribed objects.
type EventBroker struct {
	l                 *log.Logger
	stopCh            chan bool
	publishNewHeadsCh chan *utils.BlockResponse
	publishLogsCh     chan []*utils.LogResponse
	subNewHeadsCh     chan *Subscription
	subLogsCh         chan *Subscription
	unsubNewHeadsCh   chan *Subscription
	unsubLogsCh       chan *Subscription
	DebugInfo         chan int
}

// Creates a new EventsForGoEth object
func NewEventBroker() *EventBroker {
	return &EventBroker{
		l:                 log.Log(),
		stopCh:            make(chan bool, 1),
		publishNewHeadsCh: make(chan *utils.BlockResponse, NewHeadsChSize),
		publishLogsCh:     make(chan []*utils.LogResponse, LogsChSize),
		subNewHeadsCh:     make(chan *Subscription),
		subLogsCh:         make(chan *Subscription),
		unsubNewHeadsCh:   make(chan *Subscription),
		unsubLogsCh:       make(chan *Subscription),
		DebugInfo:         make(chan int),
	}
}

// SubscribeNewHeads creates a new subscription and signals the
// EventBroker subscription channel to handle the subscription map
func (eb *EventBroker) SubscribeNewHeads(ch chan *utils.BlockResponse) *Subscription {
	sub := &Subscription{
		id:         rpc.NewID(),
		typ:        NewHeadsSubscription,
		created:    time.Now(),
		newHeadsCh: ch,
		logsCh:     make(chan []*utils.LogResponse),
	}
	eb.subNewHeadsCh <- sub
	eb.l.Info().Msgf("new subscription request to New Heads with Id: [%s]", sub.id)
	return sub
}

// SubscribeLogs creates a new subscription and signals the
// EventBroker subscription channel to handle the subscription map
func (eb *EventBroker) SubscribeLogs(opts utils.LogSubscriptionOptions, ch chan []*utils.LogResponse) *Subscription {
	sub := &Subscription{
		id:         rpc.NewID(),
		typ:        LogsSubscription,
		created:    time.Now(),
		logOpts:    opts,
		newHeadsCh: make(chan *utils.BlockResponse),
		logsCh:     ch,
	}
	eb.subLogsCh <- sub
	eb.l.Info().Msgf("new subscription request to Logs with Id: [%s]", sub.id)
	return sub
}

// UnsubscribeFromNewHeads signals the EventBroker's related channel
// to delete the subscription
func (eb *EventBroker) UnsubscribeFromNewHeads(sub *Subscription) {
	eb.unsubNewHeadsCh <- sub
	eb.l.Info().Msgf("unsubscription request to New Heads with Id: [%s]", sub.id)
}

// UnsubscribeFromLogs signals the EventBroker's related channel
// to delete the subscription
func (eb *EventBroker) UnsubscribeFromLogs(sub *Subscription) {
	eb.unsubLogsCh <- sub
	eb.l.Info().Msgf("unsubscription request to Logs with Id: [%s]", sub.id)
}

// This is the main loop of the EventBroker that receives and distributes the events.
func (eb *EventBroker) Start() {
	subsNewHeads := map[rpc.ID]*Subscription{}
	subsLogs := map[rpc.ID]*Subscription{}
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
			subsNewHeads[sub.id] = sub
		case sub := <-eb.subLogsCh:
			subsLogs[sub.id] = sub
		case sub := <-eb.unsubNewHeadsCh:
			delete(subsNewHeads, sub.id)
		case sub := <-eb.unsubLogsCh:
			delete(subsLogs, sub.id)
		case msg := <-eb.publishNewHeadsCh:
			for _, v := range subsNewHeads {
				// v.newHeadsCh is buffered, use non-blocking send to protect the broker:
				// timeout preferred instead of default to be able to tolerate slight delays
				select {
				case v.newHeadsCh <- msg:
				case <-time.After(10 * time.Millisecond):
					eb.l.Warn().Msg("Publishing to New Heads channel fall into DEFAULT!")
				}
			}
		case logs := <-eb.publishLogsCh:
			if len(logs) > 0 {
				for _, s := range subsLogs {
					matchedLogs := filterLogs(logs, s.logOpts)
					if len(matchedLogs) > 0 {
						// v.logsCh is buffered, use non-blocking send to protect the broker:
						// timeout preferred instead of default to be able to tolerate slight delays
						select {
						case s.logsCh <- logs:
						case <-time.After(10 * time.Millisecond):
							eb.l.Warn().Msg("Publishing to Logs channel fall into DEFAULT!")
						}
					}
				}
			}
		}
	}
}

// Publish API that EventBroker provides. The Indexer/DB will call this
// to publish the related event to the subscribers
func (eb *EventBroker) PublishNewHeads(br *utils.BlockResponse) {
	eb.publishNewHeadsCh <- br
}

// Publish API that EventBroker provides. The Indexer/DB will call this
// to publish the related event to the subscribers
func (eb *EventBroker) PublishLogs(lr []*utils.LogResponse) {
	eb.publishLogsCh <- lr
}

// filterLogs creates a slice of Log Response matching the given criteria.
func filterLogs(rLogs []*utils.LogResponse, opts utils.LogSubscriptionOptions) []*utils.LogResponse {
	var ret []*utils.LogResponse
Logs:
	for _, log := range rLogs {
		if len(opts.Address) > 0 && !includes(opts.Address, log.Address) {
			continue
		}

		// If the number of filtered topics provided is greater than the amount of topics in logs, skip.
		if len(opts.Topics) > len(log.Topics) {
			continue
		}
		for i, sub := range opts.Topics {
			match := len(sub) == 0 // empty rule set == wildcard
			for _, topic := range sub {
				if bytes.Equal(log.Topics[i], topic) {
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

func includes(addresses utils.Addresses, address utils.Address) bool {
	for _, addr := range addresses {
		if addr == address {
			return true
		}
	}
	return false
}
