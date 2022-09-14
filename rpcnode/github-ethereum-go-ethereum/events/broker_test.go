package eventbroker_test

import (
	"aurora-relayer-go-common/broker"
	eventbroker "aurora-relayer-go-common/rpcnode/github-ethereum-go-ethereum/events"
	"aurora-relayer-go-common/utils"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

const (
	numClients          = 100
	eventTimeoutSeconds = 5
)

// Test all the flows of the Broker implementation
func TestBrokerFlows(t *testing.T) {
	eb := eventbroker.NewEventBroker()
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
			eb.PublishNewHeads(&utils.BlockResponse{
				Number: utils.IntToUint256(sentMsgCounter),
			})
			tmpLRS := make([]*utils.LogResponse, 1)
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

	// Syncs the end of event dissemination
	sentEventCounter := <-sentNumMsgCh

	// Check if sent and received event counters OK
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

func createClientAndSubscribe(eb *eventbroker.EventBroker, eventCounterCh chan int) (broker.Subscription, broker.Subscription) {
	clientNHCh := make(chan *utils.BlockResponse)
	subsNH := eb.SubscribeNewHeads(clientNHCh)

	clientLogCh := make(chan []*utils.LogResponse)
	subsLog := eb.SubscribeLogs(utils.LogSubscriptionOptions{}, clientLogCh)

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

func getNumberOfSubscriptions(eb *eventbroker.EventBroker) (int, int) {
	eb.DebugInfo <- -1
	numSubsNH := <-eb.DebugInfo
	eb.DebugInfo <- -2
	numSubsLog := <-eb.DebugInfo

	return numSubsNH, numSubsLog
}

func GenerateLogResponse() *utils.LogResponse {
	return &utils.LogResponse{
		Address: randomAddress(),
		Topics:  []utils.Bytea{randomBytea()},
	}
}

func randomAddress() utils.Address {
	return utils.Address{Address: common.BigToAddress(big.NewInt(0).SetBytes(randomBytes(20)))}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func randomBytea() utils.Bytea {
	return utils.Bytea(randomBytes(10))
}
