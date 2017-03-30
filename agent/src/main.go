package main

import (
	"errors"
	"time"
	"strings"
	"flag"

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

	//Parse flag
	unix_domain_socket := flag.String("unix_domain_socket", "", "The unix domain socket path, default:/var/tmp/net_log.sock")
	log_queue_size := flag.Int("log_queue_size", 0, "The size of log transform queue to buffer logs, default:20000")
	broker := flag.String("broker", "", "The kafka broker list, eg:localhost:9092,localhost:9093, default is empty")
	topic := flag.String("topic", "", "The kafka topic which the logs will be sent to, default:net_log")
	codec := flag.String("codec", "", "The compress codec will be used to send the logs, none, gzip, snappy, lz4 available, default:none")

	flag.Parse()

	if len(*unix_domain_socket) != 0 {
		globalConfig.Collector.UnixDomainSocket = *unix_domain_socket
	}

	if *log_queue_size != 0 {
		globalConfig.Collector.LogQueueSize = uint32(*log_queue_size)
	}

	if len(*broker) != 0 {
		brokers := strings.Split(*broker, ",")
		globalConfig.Kafka.Broker = brokers
	}

	if len(*topic) != 0 {
		globalConfig.Kafka.Topic = *topic
	}

	if len(*codec) != 0 {
		globalConfig.Kafka.CompressCodec = *codec
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
	log.Info("Version: " + config.Version)

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