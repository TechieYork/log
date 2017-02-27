package processor

import (
	"time"
	"errors"

	"github.com/DarkMetrix/log/agent/src/queue"

	log "github.com/cihub/seelog"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type KafkaProcessor struct {
	address []string
	topic string
	codec string

	//producer sarama.AsyncProducer
	producer sarama.SyncProducer

	logQueue *queue.LogQueue
}

func NewKafkaProcessor(address []string, topic string, codec string, logQueue *queue.LogQueue) *KafkaProcessor {
	return &KafkaProcessor{
		address: address,
		topic: topic,
		codec: codec,
		producer: nil,
		logQueue: logQueue,
	}
}

func (processor *KafkaProcessor) Run() error {
	//Init kafka config
	config := sarama.NewConfig()

	config.ChannelBufferSize = 10000

	config.Producer.MaxMessageBytes = 4 * 1024 * 1024
	config.Producer.Partitioner = sarama.NewHashPartitioner
	config.Producer.Return.Successes = true

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
	//processor.producer, err = sarama.NewAsyncProducer(processor.address, config)
	processor.producer, err = sarama.NewSyncProducer(processor.address, config)

	if err != nil {
		log.Warn("New kafka producer failed! err:" + err.Error())
		return err
	}

	go processor.Process()

	return nil
}

func (processor *KafkaProcessor) Close() error {
	err := processor.producer.Close()

	if err != nil {
		log.Warn("Close kafka producer failed! err:" + err.Error())
		return err
	}

	return nil
}

func (processor *KafkaProcessor) Process() error {
	for {
		logPackage, err := processor.logQueue.Pop(time.Millisecond * 10)

		if err != nil {
			//log.Warn("Pop log from queue failed(empty)! err:" + err.Error())
			continue
		}

		data, err := proto.Marshal(logPackage)

		if err != nil {
			log.Warn("Marshal log failed! err:" + err.Error())
			continue
		}

		partition, offset, err := processor.producer.SendMessage(&sarama.ProducerMessage{Topic:processor.topic, Key:sarama.StringEncoder(logPackage.GetProject()), Value:sarama.ByteEncoder(data)})

		if err != nil {
			log.Warn("Send message failed! err:" + err.Error())
			continue
		}

		log.Infof("Send message success, partition:%d, offset:%d", partition, offset)

		/*
		select {
		case processor.producer.Input() <- &sarama.ProducerMessage{Topic:processor.topic, Key:sarama.StringEncoder(logPackage.GetProject()), Value:sarama.StringEncoder(data)}:
			log.Info("Send message success!")
		case err := <- processor.producer.Errors():
			log.Warn("Send message failed! err:" + err.Error())
		}
		*/
	}
}
