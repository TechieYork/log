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

func InitLog(path string) error {
	logger, err := log.LoggerFromConfigAsFile(path)

	if err != nil {
		return err
	}

	err = log.ReplaceLogger(logger)

	if err != nil {
		return err
	}

	return nil
}

func InitConfig(path string) (*config.Config, error) {
	globalConfig := config.GetConfig()

	if globalConfig == nil {
		return nil, errors.New("Get global config failed!")
	}

	err := globalConfig.Init(path)

	if err != nil {
		return nil, err
	}

	return globalConfig, nil
}

func main() {
	defer log.Flush()

	//Initialize log using configuration from "../conf/log.config"
	err := InitLog("../conf/log.config")

	if err != nil {
		log.Warnf("Read config failed! error:%s", err)
		return
	}

	log.Info(time.Now().String(), " Starting log agent ... ")

	//Initialize the configuration from "../conf/config.json"
	log.Info("Initialize log agent configuration from ../conf/config.json ...")
	config, err := InitConfig("../conf/config.json")

	if err != nil {
		log.Warnf("Initialize log agent configuration failed! error:%s", err)
		return
	}

	log.Info("Initialize log agent configuration successed! config:", config)

	//Initialize log queue
	log.Info("Initialize log queue ...")
	logQueue := queue.NewLogQueue(config.Collector.LogQueueSize)

	if logQueue == nil {
		log.Warn("Initialize log queue failed! error:logQueue == nil")
		return
	}

	//Initialize unix domain socket collector
	log.Info("Initialize collector ...")
	logCollector := collector.NewCollector(config.Collector.UnixDomainSocket, logQueue)

	err = logCollector.Run()

	if err != nil {
		log.Warnf("Initialize log collector failed! error:%s", err)
		return
	}

	defer logCollector.Close()

	//Initialize processor
	log.Info("Initialize processor ...")
	logProcessor := processor.NewKafkaProcessor(config.Kafka.Broker, config.Kafka.Topic, config.Kafka.CompressCodec, logQueue)

	err = logProcessor.Run()

	if err != nil {
		log.Warnf("Initialize log processor failed! error:%s", err)
		return
	}

	defer logProcessor.Close()

	//Deal with signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	signalOccur := <- signalChannel

	log.Info("Signal occured, signal:", signalOccur.String())
}