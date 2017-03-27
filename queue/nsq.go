package queue

import (
	"bytes"
	"encoding/json"

	"github.com/maxence-charriere/wikirace"
	"github.com/nsqio/go-nsq"
)

type NsqEnqueuer struct {
	searchID string
	producer *nsq.Producer
}

func NewNsqEnqueuer(endpoint string) (enq *NsqEnqueuer, err error) {
	cfg := nsq.NewConfig()

	p, err := nsq.NewProducer(endpoint, cfg)
	if err != nil {
		return
	}

	enq = &NsqEnqueuer{
		searchID: "wikirace",
		producer: p,
	}
	return
}

func (enq *NsqEnqueuer) Enqueue(s wikirace.Search) error {
	var b bytes.Buffer

	enc := json.NewEncoder(&b)
	if err := enc.Encode(s); err != nil {
		return err
	}
	return enq.producer.Publish(enq.searchID, b.Bytes())
}

func (enq *NsqEnqueuer) Close() {
	enq.producer.Stop()
}

type NsqDequeuer struct {
	endpoint string
	consumer *nsq.Consumer
	stopChan chan interface{}
}

func NewNsqDequeuer(endpoint string) (deq *NsqDequeuer, err error) {
	cfg := nsq.NewConfig()

	c, err := nsq.NewConsumer("wikirace", "search", cfg)
	if err != nil {
		return
	}

	deq = &NsqDequeuer{
		endpoint: endpoint,
		consumer: c,
	}
	return
}

func (deq *NsqDequeuer) StartDequeue(h wikirace.DequeueHandler) error {
	deq.stopChan = make(chan interface{})

	deq.consumer.AddHandler(nsq.HandlerFunc(func(msg *nsq.Message) error {
		return h(msg.Body)
	}))

	if err := deq.consumer.ConnectToNSQD(deq.endpoint); err != nil {
		return err
	}
	<-deq.stopChan
	return nil
}

func (deq *NsqDequeuer) StopDequeue() {
	deq.stopChan <- true
}

func (deq *NsqDequeuer) Close() {
	deq.consumer.Stop()
}
