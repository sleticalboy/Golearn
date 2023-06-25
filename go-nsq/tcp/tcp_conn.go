package tcp

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"time"
)

var (
	url     = "127.0.0.1:4150"
	topic   = "nsq-test-1"
	channel = "nsq-test-1"
)

func tcpConnRun() {
	go startConsumer()
	startProducer()
}

func startConsumer() {
	cfg := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, cfg)
	if err != nil {
		panic(fmt.Errorf("nsq.NewConsumer() failed! %e\n", err))
	}
	// 设置消息处理函数
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		msg := string(message.Body)
		fmt.Printf("receive new message: %s\n", msg)
		if msg == "quit" {
			consumer.Stop()
			fmt.Printf("Consumer stopped!\n")
			return nil
		}
		return nil
	}))
	// 连接到 nsqd
	if err := consumer.ConnectToNSQD(url); err != nil {
		fmt.Printf("ConnectToNSQD() to %s failed %e\n", url, err)
	}
	fmt.Printf("Consumer started!\n")
	<-consumer.StopChan
}

func startProducer() {
	cfg := nsq.NewConfig()
	producer, err := nsq.NewProducer(url, cfg)
	if err != nil {
		panic(fmt.Errorf("nsq.NewProducer() failed for: %s %e\n", url, err))
	}
	// 发布消息
	var msg string
	i := 1
	for {
		if i > 30 {
			msg = "quit"
		} else {
			msg = fmt.Sprintf("seq of message is: %d", i)
		}
		if err := producer.Publish(topic, []byte(msg)); err != nil {
			fmt.Printf("Publish() to %s failed %e", "nsq-test-1\n", err)
		}
		fmt.Printf("published msg: %s\n", msg)
		time.Sleep(1 * time.Second)
		i += 1
		if i > 30 {
			break
		}
	}
	producer.Stop()
	fmt.Printf("Producer stopped!\n")
}
