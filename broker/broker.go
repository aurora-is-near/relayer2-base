package broker

import (
	"aurora-relayer-go-common/types/event"
	"aurora-relayer-go-common/types/request"
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
