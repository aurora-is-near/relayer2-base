package events_test

import (
	"crypto/rand"
	"testing"
	"time"

	"github.com/aurora-is-near/relayer2-base/broker"
	"github.com/aurora-is-near/relayer2-base/rpc/node/events"
	"github.com/aurora-is-near/relayer2-base/types/common"
	"github.com/aurora-is-near/relayer2-base/types/event"
	"github.com/aurora-is-near/relayer2-base/types/primitives"
	"github.com/aurora-is-near/relayer2-base/types/request"
	"github.com/aurora-is-near/relayer2-base/types/response"

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

func TestBrokerReturnsCorrectEventsEmptyParams(t *testing.T) {
	eb := events.NewEventBroker()
	go eb.Start()

	rcvEventCounter := 0
	eventCounterCh := make(chan int)
	go func() {
		for range eventCounterCh {
			rcvEventCounter++
		}
	}()

	filterParams := request.LogSubscriptionOptions{}

	clientLogSub := createClientAndSubscribeLogs(eb, eventCounterCh, filterParams)

	time.Sleep(1 * time.Second)

	sentNumMsgCh := make(chan int)
	go func() {
		sentMsgCounter := 0
		allLogs := make(event.Logs, 0)

		for i := 0; i < 10; i++ {
			tmpLRS := make(event.Logs, 1)
			tmpLRS[0] = GenerateLogResponse()
			allLogs = append(allLogs, tmpLRS[0])
			sentMsgCounter++
		}

		eb.PublishLogs(allLogs)
		sentNumMsgCh <- sentMsgCounter
	}()

	sentEventCounter := <-sentNumMsgCh

	time.Sleep(1 * time.Second)

	assert.NotEqual(t, 0, sentEventCounter)
	assert.Equal(t, 10, rcvEventCounter, "Number of received events does not match expected value")

	for i := 0; i < numClients; i++ {
		eb.UnsubscribeFromLogs(clientLogSub)
	}

	time.Sleep(1 * time.Second)
}

func TestBrokerReturnsCorrectEventsWithAddressAndTopics(t *testing.T) {
	eb := events.NewEventBroker()
	go eb.Start()

	rcvEventCounter := 0
	eventCounterCh := make(chan int)
	go func() {
		for range eventCounterCh {
			rcvEventCounter++
		}
	}()

	contractAddresses := []common.Address{
		common.BytesToAddress(primitives.Data20FromHex("0x2f41af687164062f118297ca10751f4b55478ae1").Bytes()),
		common.BytesToAddress(primitives.Data20FromHex("0x03b666f3488a7992b2385b12df7f35156d7b29cd").Bytes()),
		common.BytesToAddress(primitives.Data20FromHex("0x20f8aefb5697b77e0bb835a8518be70775cda1b0").Bytes()),
		common.BytesToAddress(primitives.Data20FromHex("0x63da4db6ef4e7c62168ab03982399f9588fcd198").Bytes()),
		common.BytesToAddress(primitives.Data20FromHex("0x61c9e05d1cdb1b70856c7a2c53fa9c220830633c").Bytes()),
	}

	topics := request.Topics{
		{},
		{},
		{[]byte(`0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e`)},
	}

	filterParams := request.LogSubscriptionOptions{
		Address: contractAddresses,
		Topics:  topics,
	}

	clientLogSub := createClientAndSubscribeLogs(eb, eventCounterCh, filterParams)

	time.Sleep(1 * time.Second)

	sentNumMsgCh := make(chan int)
	go func() {
		sentMsgCounter := 0
		timeout := time.After(eventTimeoutSeconds * time.Second)

		allLogs := make(event.Logs, 0)

		for {
			tmpLRS := make(event.Logs, 1)
			tmpLRS[0] = GenerateLogResponse()
			allLogs = append(allLogs, tmpLRS[0])
			sentMsgCounter += 1

			time.Sleep(10 * time.Millisecond)
			select {
			case <-timeout:
				tmpLog1 := &response.Log{
					Removed:          false,
					LogIndex:         primitives.HexUint(5),
					TransactionIndex: primitives.HexUint(0),
					TransactionHash:  primitives.Data32FromHex("0x29d3cd070a26eb34cd1c8abb70cb1e966819a342bc03965a4cd662442f712615"),
					BlockHash:        primitives.Data32FromHex("0x0579fb6c14a212998fc0e3792c2994f5f0179d8f64aa6e9059edd1f69df05155"),
					BlockNumber:      primitives.HexUint(107219211),
					Address:          primitives.Data20FromHex("0x63da4db6ef4e7c62168ab03982399f9588fcd198"),
					Data:             primitives.VarDataFromHex("0x0000000000000000000000000000000000000000000b6afb14c2d46e19ffffc40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000305d9662647959"),
					Topics: []primitives.Data32{
						primitives.Data32FromHex("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"),
						primitives.Data32FromHex("0x0000000000000000000000002cb45edb4517d5947afde3beabf95a582506858b"),
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
					},
				}
				allLogs = append(allLogs, tmpLog1)
				sentMsgCounter += 1

				tmpLog2 := &response.Log{
					Removed:          false,
					LogIndex:         primitives.HexUint(5),
					TransactionIndex: primitives.HexUint(0),
					TransactionHash:  primitives.Data32FromHex("0x29d3cd070a26eb34cd1c8abb70cb1e966819a342bc03965a4cd662442f712615"),
					BlockHash:        primitives.Data32FromHex("0x0579fb6c14a212998fc0e3792c2994f5f0179d8f64aa6e9059edd1f69df05155"),
					BlockNumber:      primitives.HexUint(107219211),
					Address:          primitives.Data20FromHex("0x63da4db6ef4e7c62168ab03982399f9588fcd198"),
					Data:             primitives.VarDataFromHex("0x0000000000000000000000000000000000000000000b6afb14c2d46e19ffffc40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000305d9662647959"),
					Topics: []primitives.Data32{
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
						primitives.Data32FromHex("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"),
					},
				}

				allLogs = append(allLogs, tmpLog2)
				sentMsgCounter += 1

				eb.PublishLogs(allLogs)

				sentNumMsgCh <- sentMsgCounter
				return
			default:
			}
		}
	}()

	sentEventCounter := <-sentNumMsgCh
	time.Sleep(1 * time.Second)

	assert.NotEqual(t, 0, sentEventCounter)
	assert.Equal(t, 1, rcvEventCounter, "Number of received events does not match expected value")

	for i := 0; i < numClients; i++ {
		eb.UnsubscribeFromLogs(clientLogSub)
	}

	time.Sleep(1 * time.Second)
}

func TestBrokerReturnsCorrectEventsWithTopics(t *testing.T) {
	eb := events.NewEventBroker()
	go eb.Start()

	rcvEventCounter := 0
	eventCounterCh := make(chan int)
	go func() {
		for range eventCounterCh {
			rcvEventCounter++
		}
	}()

	topics := request.Topics{
		{[]byte(`0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e`)},
	}

	filterParams := request.LogSubscriptionOptions{
		Topics: topics,
	}

	clientLogSub := createClientAndSubscribeLogs(eb, eventCounterCh, filterParams)

	time.Sleep(1 * time.Second)

	sentNumMsgCh := make(chan int)
	go func() {
		sentMsgCounter := 0
		timeout := time.After(eventTimeoutSeconds * time.Second)

		allLogs := make(event.Logs, 0)

		for {
			tmpLRS := make(event.Logs, 1)
			tmpLRS[0] = GenerateLogResponse()
			allLogs = append(allLogs, tmpLRS[0])
			sentMsgCounter += 1

			time.Sleep(10 * time.Millisecond)
			select {
			case <-timeout:
				tmpLog1 := &response.Log{
					Removed:          false,
					LogIndex:         primitives.HexUint(5),
					TransactionIndex: primitives.HexUint(0),
					TransactionHash:  primitives.Data32FromHex("0x29d3cd070a26eb34cd1c8abb70cb1e966819a342bc03965a4cd662442f712615"),
					BlockHash:        primitives.Data32FromHex("0x0579fb6c14a212998fc0e3792c2994f5f0179d8f64aa6e9059edd1f69df05155"),
					BlockNumber:      primitives.HexUint(107219211),
					Address:          primitives.Data20FromHex("0x63da4db6ef4e7c62168ab03982399f9588fcd198"),
					Data:             primitives.VarDataFromHex("0x0000000000000000000000000000000000000000000b6afb14c2d46e19ffffc40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000305d9662647959"),
					Topics: []primitives.Data32{
						primitives.Data32FromHex("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"),
						primitives.Data32FromHex("0x0000000000000000000000002cb45edb4517d5947afde3beabf95a582506858b"),
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
					},
				}
				allLogs = append(allLogs, tmpLog1)
				sentMsgCounter += 1

				tmpLog2 := &response.Log{
					Removed:          false,
					LogIndex:         primitives.HexUint(5),
					TransactionIndex: primitives.HexUint(0),
					TransactionHash:  primitives.Data32FromHex("0x29d3cd070a26eb34cd1c8abb70cb1e966819a342bc03965a4cd662442f712615"),
					BlockHash:        primitives.Data32FromHex("0x0579fb6c14a212998fc0e3792c2994f5f0179d8f64aa6e9059edd1f69df05155"),
					BlockNumber:      primitives.HexUint(107219211),
					Address:          primitives.Data20FromHex("0x63da4db6ef4e7c62168ab03982399f9588fcd198"),
					Data:             primitives.VarDataFromHex("0x0000000000000000000000000000000000000000000b6afb14c2d46e19ffffc40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000305d9662647959"),
					Topics: []primitives.Data32{
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
						primitives.Data32FromHex("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"),
					},
				}

				allLogs = append(allLogs, tmpLog2)
				sentMsgCounter += 1

				eb.PublishLogs(allLogs)

				sentNumMsgCh <- sentMsgCounter
				return
			default:
			}
		}
	}()

	sentEventCounter := <-sentNumMsgCh
	time.Sleep(1 * time.Second)

	assert.NotEqual(t, 0, sentEventCounter)
	assert.Equal(t, 1, rcvEventCounter, "Number of received events does not match expected value")

	for i := 0; i < numClients; i++ {
		eb.UnsubscribeFromLogs(clientLogSub)
	}

	time.Sleep(1 * time.Second)
}

func TestBrokerReturnsCorrectEventsWithAddress(t *testing.T) {
	eb := events.NewEventBroker()
	go eb.Start()

	rcvEventCounter := 0
	eventCounterCh := make(chan int)
	go func() {
		for range eventCounterCh {
			rcvEventCounter++
		}
	}()

	contractAddresses := []common.Address{
		common.BytesToAddress(primitives.Data20FromHex("0x2f41af687164062f118297ca10751f4b55478ae1").Bytes()),
		common.BytesToAddress(primitives.Data20FromHex("0x03b666f3488a7992b2385b12df7f35156d7b29cd").Bytes()),
		common.BytesToAddress(primitives.Data20FromHex("0x20f8aefb5697b77e0bb835a8518be70775cda1b0").Bytes()),
		common.BytesToAddress(primitives.Data20FromHex("0x63da4db6ef4e7c62168ab03982399f9588fcd198").Bytes()),
		common.BytesToAddress(primitives.Data20FromHex("0x61c9e05d1cdb1b70856c7a2c53fa9c220830633c").Bytes()),
	}

	filterParams := request.LogSubscriptionOptions{
		Address: contractAddresses,
	}

	clientLogSub := createClientAndSubscribeLogs(eb, eventCounterCh, filterParams)

	time.Sleep(1 * time.Second)

	sentNumMsgCh := make(chan int)
	go func() {
		sentMsgCounter := 0
		timeout := time.After(eventTimeoutSeconds * time.Second)

		allLogs := make(event.Logs, 0)

		for {
			tmpLRS := make(event.Logs, 1)
			tmpLRS[0] = GenerateLogResponse()
			allLogs = append(allLogs, tmpLRS[0])
			sentMsgCounter += 1

			time.Sleep(10 * time.Millisecond)
			select {
			case <-timeout:
				tmpLog1 := &response.Log{
					Removed:          false,
					LogIndex:         primitives.HexUint(5),
					TransactionIndex: primitives.HexUint(0),
					TransactionHash:  primitives.Data32FromHex("0x29d3cd070a26eb34cd1c8abb70cb1e966819a342bc03965a4cd662442f712615"),
					BlockHash:        primitives.Data32FromHex("0x0579fb6c14a212998fc0e3792c2994f5f0179d8f64aa6e9059edd1f69df05155"),
					BlockNumber:      primitives.HexUint(107219211),
					Address:          primitives.Data20FromHex("0x20f8aefb5697b77e0bb835a8518be70775cda1b0"),
					Data:             primitives.VarDataFromHex("0x0000000000000000000000000000000000000000000b6afb14c2d46e19ffffc40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000305d9662647959"),
					Topics: []primitives.Data32{
						primitives.Data32FromHex("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"),
						primitives.Data32FromHex("0x0000000000000000000000002cb45edb4517d5947afde3beabf95a582506858b"),
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
					},
				}
				allLogs = append(allLogs, tmpLog1)
				sentMsgCounter += 1

				tmpLog2 := &response.Log{
					Removed:          false,
					LogIndex:         primitives.HexUint(5),
					TransactionIndex: primitives.HexUint(0),
					TransactionHash:  primitives.Data32FromHex("0x29d3cd070a26eb34cd1c8abb70cb1e966819a342bc03965a4cd662442f712615"),
					BlockHash:        primitives.Data32FromHex("0x0579fb6c14a212998fc0e3792c2994f5f0179d8f64aa6e9059edd1f69df05155"),
					BlockNumber:      primitives.HexUint(107219211),
					Address:          primitives.Data20FromHex("0x63da4db6ef4e7c62168ab03982399f9588fcd198"),
					Data:             primitives.VarDataFromHex("0x0000000000000000000000000000000000000000000b6afb14c2d46e19ffffc40000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000305d9662647959"),
					Topics: []primitives.Data32{
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
						primitives.Data32FromHex("0x0000000000000000000000005eec60f348cb1d661e4a5122cf4638c7db7a886e"),
						primitives.Data32FromHex("0xd78ad95fa46c994b6551d0da85fc275fe613ce37657fb8d5e3d130840159d822"),
					},
				}

				allLogs = append(allLogs, tmpLog2)
				sentMsgCounter += 1

				eb.PublishLogs(allLogs)

				sentNumMsgCh <- sentMsgCounter
				return
			default:
			}
		}
	}()

	sentEventCounter := <-sentNumMsgCh
	time.Sleep(1 * time.Second)

	assert.NotEqual(t, 0, sentEventCounter)
	assert.Equal(t, 2, rcvEventCounter, "Number of received events does not match expected value")

	for i := 0; i < numClients; i++ {
		eb.UnsubscribeFromLogs(clientLogSub)
	}

	time.Sleep(1 * time.Second)
}

func createClientAndSubscribeLogs(eb *events.EventBroker, eventCounterCh chan int, subOptions request.LogSubscriptionOptions) broker.Subscription {
	clientLogCh := make(chan event.Logs)
	subsLog := eb.SubscribeLogs(subOptions, clientLogCh)

	go func() {
		for logs := range clientLogCh {
			for range logs {
				eventCounterCh <- 1
			}
		}
	}()

	return subsLog
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
