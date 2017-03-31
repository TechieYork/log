package processor

import (
	"errors"
	"time"

	"github.com/DarkMetrix/log/agent/src/queue"

	"github.com/Shopify/sarama"
	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
)

//Kafka log processor
type KafkaProcessor struct {
	address []string //Broker list
	topic   string   //Log topic, default "net_log"
	codec   string   //Compression codec, "none", "gzip", "snappy" or "lz4"

	producer sarama.AsyncProducer //Async kafka producer

	logQueue *queue.LogQueue //Log queue to buffer log which process to send later
}

//New kafka processor function
func NewKafkaProcessor(address []string, topic string, codec string, logQueue *queue.LogQueue) *KafkaProcessor {
	return &KafkaProcessor{
		address:  address,
		topic:    topic,
		codec:    codec,
		producer: nil,
		logQueue: logQueue,
	}
}

//Run to process log
func (processor *KafkaProcessor) Run() error {
	//Init kafka config
	config := sarama.NewConfig()

	config.ChannelBufferSize = 10000

	config.Producer.MaxMessageBytes = 4 * 1024 * 1024
	config.Producer.Partitioner = sarama.NewHashPartitioner
	//config.Producer.Return.Successes = true

	var codec sarama.CompressionCodec

	switch processor.codec {
	case "none":
		codec = sarama.CompressionNone
	case "gzip":
		codec = sarama.CompressionGZIP
	case "snappy":
		codec = sarama.CompressionSnappy
	case "lz4":
		codec = sarama.CompressionLZ4
	default:
		return errors.New("Codec error, not one of 'none', 'gzip', 'snappy' or 'lz4'")
	}

	config.Producer.Compression = codec

	//Init kafka producer
	var err error
	processor.producer, err = sarama.NewAsyncProducer(processor.address, config)

	if err != nil {
		log.Warn("New kafka producer failed! err:" + err.Error())
		return err
	}

	go processor.Process()

	return nil
}

//Close kafka producer
func (processor *KafkaProcessor) Close() error {
	err := processor.producer.Close()

	if err != nil {
		log.Warn("Close kafka producer failed! err:" + err.Error())
		return err
	}

	return nil
}

//Process go routine function
func (processor *KafkaProcessor) Process() error {
	//Loop to pop log from log queue
	for {
		//Pop log
		logPackage, err := processor.logQueue.Pop(time.Millisecond * 10)

		if err != nil {
			//log.Warn("Pop log from queue failed(empty)! err:" + err.Error())
			continue
		}

		//Marshal log
		data, err := proto.Marshal(logPackage)

		if err != nil {
			log.Warn("Marshal log failed! err:" + err.Error())
			continue
		}

		select {
		case processor.producer.Input() <- &sarama.ProducerMessage{Topic: processor.topic, Key: sarama.StringEncoder(logPackage.GetProject()), Value: sarama.ByteEncoder(data)}:
			//log.Info("Send message success!")
		case err := <-processor.producer.Errors():
			log.Warn("Send message failed! err:" + err.Error())
		}
	}
}
