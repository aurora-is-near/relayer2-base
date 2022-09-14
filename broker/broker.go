package broker

import (
	"aurora-relayer-go-common/utils"
)

type SubID string

type Subscription interface {
	GetId() SubID
	GetNewHeadsCh() chan *utils.BlockResponse
	GetLogsCh() chan []*utils.LogResponse
	GetLogsSubOpts() utils.LogSubscriptionOptions
}

type Broker interface {
	SubscribeNewHeads(chan *utils.BlockResponse) Subscription
	SubscribeLogs(utils.LogSubscriptionOptions, chan []*utils.LogResponse) Subscription
	UnsubscribeFromNewHeads(Subscription)
	UnsubscribeFromLogs(Subscription)
	PublishNewHeads(*utils.BlockResponse)
	PublishLogs([]*utils.LogResponse)
}
