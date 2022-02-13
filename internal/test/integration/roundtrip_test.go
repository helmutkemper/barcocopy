//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/barcostreams/barco/internal/conf"
	. "github.com/barcostreams/barco/internal/test/integration"
	. "github.com/barcostreams/barco/internal/types"
	"github.com/klauspost/compress/zstd"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

const consumerContentType = "application/vnd.barco.consumermessage"

func TestData(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration test suite")
}

var _ = Describe("A 3 node cluster", func() {
	// Note that on macos you need to manually create the alias for the loopback addresses, for example
	// sudo ifconfig lo0 alias 127.0.0.2 up && sudo ifconfig lo0 alias 127.0.0.3 up

	Describe("Producing and consuming", func() {
		var b1 *TestBroker
		var b2 *TestBroker
		var b3 *TestBroker

		BeforeEach(func ()  {
			b1 = NewTestBroker(0)
			b2 = NewTestBroker(1)
			b3 = NewTestBroker(2)
		})

		AfterEach(func ()  {
			b1.Shutdown()
			b2.Shutdown()
			b3.Shutdown()
		})

		It("should work", func() {
			start := time.Now()
			b1.WaitForStart()
			b2.WaitForStart()
			b3.WaitForStart()

			log.Debug().Msgf("All brokers started successfully")

			b1.WaitOutput("Setting committed version 1 with leader 0 for range")
			log.Debug().Msgf("Waited for first broker")
			b2.WaitOutput("Setting committed version 1 with leader 1 for range")
			log.Debug().Msgf("Waited for second broker")
			b3.WaitOutput("Setting committed version 1 with leader 2 for range")
			log.Debug().Msgf("Waited for third broker")

			message := `{"hello": "world"}`

			// Test with HTTP/2
			client := NewTestClient(nil)
			resp := client.ProduceJson(0, "abc", message, "")
			expectResponseOk(resp)

			// Use different partition keys
			// expectResponseOk(client.ProduceJson(0, "abc", message, "123")) // B0
			expectResponseOk(client.ProduceJson(0, "abc", message, "567")) // Re-routed to B1
			expectResponseOk(client.ProduceJson(0, "abc", message, "234")) // Re-routed to B2

			client.RegisterAsConsumer(3, `{"id": "c1", "group": "g1", "topics": ["abc"]}`)
			log.Debug().Msgf("Registered as consumer")

			// Wait for the consumer to be considered
			time.Sleep(500 * time.Millisecond)

			resp = client.ConsumerPoll(0)
			defer resp.Body.Close()
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
			Expect(resp.Header.Get("Content-Type")).To(Equal(consumerContentType))
			var messageLength uint16
			binary.Read(resp.Body, conf.Endianness, &messageLength)
			Expect(messageLength).To(Equal(uint16(1)))
			item := unmarshalConsumerResponseItem(resp.Body)
			Expect(*item.topic).To(Equal(TopicDataId{
				Name:       "abc",
				Token:      -9223372036854775808,
				RangeIndex: 1,
				GenId:      1,
			}))

			Expect(item.records).To(HaveLen(1))
			Expect(item.records[0].timestamp.UnixMilli()).To(BeNumerically(">=", start.UnixMilli()))
			Expect(item.records[0].timestamp.UnixMilli()).To(BeNumerically("<=", time.Now().UnixMilli()))
			Expect(item.records[0].body).To(Equal(message))

			// Test with HTTP/1
			expectResponseOk(NewTestClient(&TestClientOptions{HttpVersion: 1}).ProduceJson(0, "abc", message, ""))

			client.Close()
		})
	})
})

func expectResponseOk(resp *http.Response) {
	defer resp.Body.Close()
	Expect(ReadBody(resp)).To(Equal("OK"))
	Expect(resp.StatusCode).To(Equal(200))
}

type consumerResponseItem struct {
	topic *TopicDataId
	records []record
}

type record struct {
	timestamp time.Time
	body string
}

func unmarshalConsumerResponseItem(r io.Reader) consumerResponseItem {
	item := consumerResponseItem{}
	item.topic = unmarshalTopicId(r)
	payloadLength := int32(0)
	binary.Read(r, conf.Endianness, &payloadLength)
	payload := make([]byte, payloadLength)
	n, err := r.Read(payload)
	Expect(err).NotTo(HaveOccurred())
	Expect(n).To(Equal(int(payloadLength)))

	payloadReader, err := zstd.NewReader(bytes.NewReader(payload))
	Expect(err).NotTo(HaveOccurred())
	uncompressed, err := io.ReadAll(payloadReader)
	Expect(err).NotTo(HaveOccurred())
	recordsReader := bytes.NewReader(uncompressed)
	item.records = make([]record, 0)
	// for recordsReader.Len()
	item.records = append(item.records, unmarshalRecord(recordsReader))

	return item
}

func unmarshalTopicId(r io.Reader) *TopicDataId {
	topic := TopicDataId{}
	topicLength := uint8(0)
	err := binary.Read(r, conf.Endianness, &topic.Token)
	Expect(err).NotTo(HaveOccurred())
	err = binary.Read(r, conf.Endianness, &topic.RangeIndex)
	Expect(err).NotTo(HaveOccurred())
	err = binary.Read(r, conf.Endianness, &topic.GenId)
	Expect(err).NotTo(HaveOccurred())
	err = binary.Read(r, conf.Endianness, &topicLength)
	Expect(err).NotTo(HaveOccurred())
	topicName := make([]byte, topicLength)
	_, err = r.Read(topicName)
	Expect(err).NotTo(HaveOccurred())
	topic.Name = string(topicName)

	return &topic
}

func unmarshalRecord(r io.Reader) record {
	length := uint32(0)
	timestamp := int64(0)
	result := record{}
	err := binary.Read(r, conf.Endianness, &timestamp)
	Expect(err).NotTo(HaveOccurred())
	result.timestamp = time.UnixMicro(timestamp)
	err = binary.Read(r, conf.Endianness, &length)
	body := make([]byte, length)
	n, err := r.Read(body)
	Expect(err).NotTo(HaveOccurred())
	Expect(n).To(Equal(int(length)))
	result.body = string(body)
	return result
}