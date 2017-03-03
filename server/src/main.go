package main

import (
	"os"
	"os/signal"
	"time"
	"errors"

	"github.com/DarkMetrix/log/server/src/config"
	"github.com/DarkMetrix/log/server/src/queue"
	"github.com/DarkMetrix/log/server/src/processor"

	log "github.com/cihub/seelog"
	"github.com/DarkMetrix/log/server/src/sinker"
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

	log.Infof("    sinker.log_path: %s", globalConfig.Sinker.LogPath)
	log.Infof("    sinker.log_file_max_size: %d Bytes", globalConfig.Sinker.LogFileMaxSize)
	log.Infof("    sinker.log_flush_duration: %d Seconds", globalConfig.Sinker.LogFlushDuration)
	log.Infof("    sinker.log_flush_duration: %d Seconds", globalConfig.Sinker.LogQueueSize)

	log.Infof("    kafka.broker: %s", globalConfig.Kafka.Broker)
	log.Infof("    kafka.topic: %s", globalConfig.Kafka.Topic)
	log.Infof("    kafka.partitions: %v", globalConfig.Kafka.Partitions)

	return globalConfig
}

//Init log queue
func InitLogQueue(config *config.Config) *queue.LogQueue {
	log.Info("Initialize log queue ...")

	logQueue := queue.NewLogQueue(config.Sinker.LogQueueSize)

	if logQueue == nil {
		panic(errors.New("Initialize log queue failed! error:logQueue == nil"))
	}

	return logQueue
}

//Init processor
func InitProcessor(config *config.Config, logQueue *queue.LogQueue) *processor.KafkaProcessor {
	log.Info("Initialize processor ...")

	logProcessor := processor.NewKafkaProcessor(config.Kafka.Broker, config.Kafka.Topic, config.Kafka.Partitions, logQueue)

	err := logProcessor.Run()

	if err != nil {
		panic(err)
	}

	return logProcessor
}

//Init sinker
func InitSinker(config *config.Config, logQueue *queue.LogQueue) *sinker.Sinker {
	log.Info("Initialize sinker ...")

	logSinker := sinker.NewSinker(config.Sinker.LogPath, config.Sinker.LogFileMaxSize, config.Sinker.LogFlushDuration, logQueue)

	err := logSinker.Run()

	if err != nil {
		panic(err)
	}

	return logSinker
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

	log.Info(time.Now().String(), " Log server starting ... ")

	//Initialize the configuration from "../conf/config.json"
	config := InitConfig("../conf/config.json")

	log.Info(time.Now().String(), " Log server started!")

	//Initialize log queue
	logQueue := InitLogQueue(config)

	//Initialize processor
	logProcessor := InitProcessor(config, logQueue)
	defer logProcessor.Close()

	//Initialize sinker
	logSinker := InitSinker(config, logQueue)
	defer logSinker.Close()

	//Deal with signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	signalOccur := <-signalChannel

	log.Info("Signal occured, signal:", signalOccur.String())
}