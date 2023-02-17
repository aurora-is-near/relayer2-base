package events_test

import (
	"crypto/rand"
	"github.com/aurora-is-near/relayer2-base/broker"
	"github.com/aurora-is-near/relayer2-base/rpcnode/github-ethereum-go-ethereum/events"
	"github.com/aurora-is-near/relayer2-base/types/event"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/request"
	"github.com/aurora-is-near/relayer2-base/types/response"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	numClients          = 100
	eventTimeoutSeconds = 5
)

// Test all the flows of the Broker implementation
func TestBrokerFlows(t *testing.T) {
	eb := events.NewEventBroker()
	go eb.Start()

	// Handles eventCounter channel to calculate the number of received events (by all clients)
	rcvEventCounter := 0
	eventCounterCh := make(chan int)
	go func() {
		for range eventCounterCh {
			rcvEventCounter++
		}
	}()

	// Create client and subscribe to the events
	clientNHSubs := make([]broker.Subscription, numClients)
	clientLogSubs := make([]broker.Subscription, numClients)
	for i := 0; i < numClients; i++ {
		clientNHSubs[i], clientLogSubs[i] = createClientAndSubscribe(eb, eventCounterCh)
		time.Sleep(1 * time.Millisecond)
	}

	time.Sleep(1 * time.Second)

	// Check if subscriptions OK
	numSubsNH, numSubsLog := getNumberOfSubscriptions(eb)
	assert.Equal(t, numSubsNH, numClients)
	assert.Equal(t, numSubsLog, numClients)

	// Start publishing events
	sentNumMsgCh := make(chan int)
	go func() {
		sentMsgCounter := 0
		timeout := time.After(eventTimeoutSeconds * time.Second)
		for {
			eb.PublishNewHeads(&response.Block{
				Number: primitives.HexUint(sentMsgCounter),
			})
			tmpLRS := make(event.Logs, 1)
			tmpLRS[0] = GenerateLogResponse()
			eb.PublishLogs(tmpLRS)
			sentMsgCounter += 1

			time.Sleep(10 * time.Millisecond)
			select {
			case <-timeout:
				sentNumMsgCh <- sentMsgCounter
				return
			default:
			}

		}
	}()

	// Syncs the end of types dissemination
	sentEventCounter := <-sentNumMsgCh

	// Check if sent and received types counters OK
	assert.Equal(t, sentEventCounter*numClients*2, rcvEventCounter)

	// Now unsubscribe
	for i := 0; i < numClients; i++ {
		eb.UnsubscribeFromNewHeads(clientNHSubs[i])
		eb.UnsubscribeFromLogs(clientLogSubs[i])
	}

	time.Sleep(1 * time.Second)

	// Check if unsubscription OK
	numSubsNH, numSubsLog = getNumberOfSubscriptions(eb)
	assert.Equal(t, numSubsNH, 0)
	assert.Equal(t, numSubsLog, 0)

}

func createClientAndSubscribe(eb *events.EventBroker, eventCounterCh chan int) (broker.Subscription, broker.Subscription) {
	clientNHCh := make(chan event.Block)
	subsNH := eb.SubscribeNewHeads(clientNHCh)

	clientLogCh := make(chan event.Logs)
	subsLog := eb.SubscribeLogs(request.LogSubscriptionOptions{}, clientLogCh)

	go func() {
		for {
			select {
			case <-clientNHCh:
				eventCounterCh <- 1
			case <-clientLogCh:
				eventCounterCh <- 1
			}
		}
	}()

	return subsNH, subsLog
}

func getNumberOfSubscriptions(eb *events.EventBroker) (int, int) {
	eb.DebugInfo <- -1
	numSubsNH := <-eb.DebugInfo
	eb.DebugInfo <- -2
	numSubsLog := <-eb.DebugInfo

	return numSubsNH, numSubsLog
}

func GenerateLogResponse() *response.Log {
	return &response.Log{
		Address: randomAddress(),
		Topics:  randomTopics(),
	}
}

func randomAddress() primitives.Data20 {
	return primitives.Data20FromBytes(randomBytes(20))
}

func randomTopics() []primitives.Data32 {
	return []primitives.Data32{primitives.Data32FromBytes(randomBytes(10))}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
