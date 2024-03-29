package broker

import (
	"github.com/aurora-is-near/relayer2-base/types/event"
	"github.com/aurora-is-near/relayer2-base/types/request"
)

type SubID string

type Subscription interface {
	GetId() SubID
	GetNewHeadsCh() chan event.Block
	GetLogsCh() chan event.Logs
	GetLogsSubOpts() request.LogSubscriptionOptions
}

type Broker interface {
	SubscribeNewHeads(chan event.Block) Subscription
	SubscribeLogs(request.LogSubscriptionOptions, chan event.Logs) Subscription
	UnsubscribeFromNewHeads(Subscription)
	UnsubscribeFromLogs(Subscription)
	PublishNewHeads(event.Block)
	PublishLogs(event.Logs)
}
