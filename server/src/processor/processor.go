package processor

import (
	"github.com/DarkMetrix/log/proto"
	"github.com/DarkMetrix/log/server/src/queue"

	log "github.com/cihub/seelog"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

//Kafka log processor
type KafkaProcessor struct {
	address []string                    //Broker list
	topic string                        //Log topic, default "dark_metrix_log"
	partitions []int32                  //Partitions to consume

	consumer sarama.Consumer            //Kafka consumer
	partitionConsumers []sarama.PartitionConsumer //Kafka partition consumers

	logQueue *queue.LogQueue            //Log queue to buffer log which to sink later
}

//New kafka processor function
func NewKafkaProcessor(address []string, topic string, partitions []int32, logQueue *queue.LogQueue) *KafkaProcessor {
	return &KafkaProcessor{
		address: address,
		topic: topic,
		partitions: partitions,
		consumer: nil,
		logQueue: logQueue,
	}
}

//Run to process log
func (processor *KafkaProcessor) Run() error {
	//Init kafka config
	config := sarama.NewConfig()

	config.ChannelBufferSize = 10000

	//Init kafka consumer
	var err error
	processor.consumer, err = sarama.NewConsumer(processor.address, config)

	if err != nil {
		log.Warn("New kafka consumer failed! err:" + err.Error())
		return err
	}

	//Init all partition consumers
	for _, partition := range processor.partitions {
		partitionConsumer, err := processor.consumer.ConsumePartition(processor.topic, partition, sarama.OffsetNewest)

		if err != nil {
			log.Warn("New kafka partition consumer failed! err:" + err.Error())
			return err
		}

		processor.partitionConsumers = append(processor.partitionConsumers, partitionConsumer)

		go processor.process(partitionConsumer)
	}

	return nil
}

//Close kafka consumer
func (processor *KafkaProcessor) Close() error {
	//Close all partition consumers
	for _, partitionConsumer := range processor.partitionConsumers {
		err := partitionConsumer.Close()

		if err != nil {
			log.Warn("Close kafka partition consumer failed! err:" + err.Error())
			return err
		}
	}

	//Close consumer
	err := processor.consumer.Close()

	if err != nil {
		log.Warn("Close kafka consumer failed! err:" + err.Error())
		return err
	}

	return nil
}

//Process go routine function
func (processor *KafkaProcessor) process(partitionConsumer sarama.PartitionConsumer) error {
	//Loop to consume log from kafka
	for {
		select {
		//Consume log
		case message := <- partitionConsumer.Messages():
			var logPackage log_proto.LogPackage

			//Unmashal log
			err := proto.Unmarshal(message.Value, &logPackage)

			if err != nil {
				log.Warn("Processor couldn't unmarshal received buffer! err:" + err.Error())
				continue
			}

			//Push to log queue
			err = processor.logQueue.Push(&logPackage)

			if err != nil {
				log.Warn("Processor push log to queue failed! err:" + err.Error())
				continue
			}
		}
	}
}
