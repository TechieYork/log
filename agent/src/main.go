package main

import (
	"errors"
	"time"

	"os"
	"os/signal"

	"github.com/DarkMetrix/log/agent/src/config"
	"github.com/DarkMetrix/log/agent/src/queue"
	"github.com/DarkMetrix/log/agent/src/collector"
	"github.com/DarkMetrix/log/agent/src/processor"

	log "github.com/cihub/seelog"
)

//Init log
func InitLog(path string) {
	logger, err := log.LoggerFromConfigAsFile(path)

	if err != nil {
		panic(err)
	}

	err = log.ReplaceLogger(logger)

	if err != nil {
		panic(err)
	}
}

//Init config
func InitConfig(path string) *config.Config {
	log.Info("Initialize log agent configuration from " + path + " ...")

	globalConfig := config.GetConfig()

	if globalConfig == nil {
		panic(errors.New("Get global config failed!"))
	}

	err := globalConfig.Init(path)

	if err != nil {
		panic(err)
	}

	log.Info("Config:")

	log.Infof("    collector.unix_domain_socket: %s", globalConfig.Collector.UnixDomainSocket)
	log.Infof("    collector.log_queue_size: %d", globalConfig.Collector.LogQueueSize)

	log.Infof("    kafka.broker: %s", globalConfig.Kafka.Broker)
	log.Infof("    kafka.topic: %s", globalConfig.Kafka.Topic)
	log.Infof("    kafka.compress_codec: %s", globalConfig.Kafka.CompressCodec)

	return globalConfig
}

//Init log queue
func InitLogQueue(config *config.Config) *queue.LogQueue {
	log.Info("Initialize log queue ...")

	logQueue := queue.NewLogQueue(config.Collector.LogQueueSize)

	if logQueue == nil {
		panic(errors.New("Initialize log queue failed! error:logQueue == nil"))
	}

	return logQueue
}

//Init collector
func InitCollector(config *config.Config, logQueue *queue.LogQueue) *collector.Collector {
	log.Info("Initialize collector ...")

	logCollector := collector.NewCollector(config.Collector.UnixDomainSocket, logQueue)

	err := logCollector.Run()

	if err != nil {
		panic(err)
	}

	return logCollector
}

//Init processor
func InitProcessor(config *config.Config, logQueue *queue.LogQueue) *processor.KafkaProcessor {
	log.Info("Initialize processor ...")

	logProcessor := processor.NewKafkaProcessor(config.Kafka.Broker, config.Kafka.Topic, config.Kafka.CompressCodec, logQueue)

	err := logProcessor.Run()

	if err != nil {
		panic(err)
	}

	return logProcessor
}

func main() {
	defer log.Flush()

	defer func() {
		err := recover()

		if err != nil {
			log.Critical("Got panic, err:", err)
		}
	} ()

	//Initialize log using configuration from "../conf/log.config"
	InitLog("../conf/log.config")

	log.Info(time.Now().String(), " Log agent starting ... ")

	//Initialize the configuration from "../conf/config.json"
	config := InitConfig("../conf/config.json")

	//Initialize log queue
	logQueue := InitLogQueue(config)

	//Initialize unix domain socket collector
	logCollector := InitCollector(config, logQueue)
	defer logCollector.Close()

	//Initialize processor
	logProcessor := InitProcessor(config, logQueue)
	defer logProcessor.Close()

	log.Info(time.Now().String(), " Log agent started!")

	//Deal with signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	signalOccur := <- signalChannel

	log.Info("Signal occured, signal:", signalOccur.String())
}