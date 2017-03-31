package main

import (
	"flag"
	"errors"
	"os"
	"os/signal"
	"time"
	"strconv"
	"strings"

	"github.com/DarkMetrix/log/server/src/config"
	"github.com/DarkMetrix/log/server/src/processor"
	"github.com/DarkMetrix/log/server/src/queue"

	"github.com/DarkMetrix/log/server/src/sinker"
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
	log_path := flag.String("log_path", "", "The path to save all the log files, default:/tmp/dark_metrix/log/")
	log_file_max_size := flag.Int64("log_file_max_size", 0, "The max size of one log file, default:500M(524288000), minimum 1M")
	log_flush_duration := flag.Int("log_flush_duration", 0, "The duration to call flush() to flush the buffer to disk, default:3 seconds, minimum 1 second")
	log_queue_size := flag.Int("log_queue_size", 0, "The size of log transform queue to buffer logs for each file, default:10000, minimum 1000")
	broker := flag.String("broker", "", "The kafka broker list, eg:localhost:9092,localhost:9093, default is empty")
	topic := flag.String("topic", "", "The kafka topic which the logs will be sent to, default:net_log")
	partition := flag.String("partition", "", "The partition list to consume the logs, default is empty")

	flag.Parse()

	if len(*log_path) != 0 {
		globalConfig.Sinker.LogPath = *log_path
	}

	if *log_file_max_size != 0 {
		globalConfig.Sinker.LogFileMaxSize = *log_file_max_size
	}

	if *log_flush_duration != 0 {
		globalConfig.Sinker.LogFlushDuration = uint32(*log_flush_duration)
	}

	if *log_queue_size != 0 {
		globalConfig.Sinker.LogQueueSize = uint32(*log_queue_size)
	}

	if len(*broker) != 0 {
		brokers := strings.Split(*broker, ",")
		globalConfig.Kafka.Broker = brokers
	}

	if len(*topic) != 0 {
		globalConfig.Kafka.Topic = *topic
	}

	if len(*partition) != 0 {
		partitions := strings.Split(*partition, ",")
		globalConfig.Kafka.Partitions = []int32{}

		for _, value := range partitions {
			number, err := strconv.Atoi(value)

			if err != nil {
				panic(err)
			}

			globalConfig.Kafka.Partitions = append(globalConfig.Kafka.Partitions, int32(number))
		}
	}

	log.Info("Config:")

	log.Infof("    sinker.log_path: %s", globalConfig.Sinker.LogPath)
	log.Infof("    sinker.log_file_max_size: %d Bytes", globalConfig.Sinker.LogFileMaxSize)
	log.Infof("    sinker.log_flush_duration: %d Seconds", globalConfig.Sinker.LogFlushDuration)
	log.Infof("    sinker.log_queue_size: %d", globalConfig.Sinker.LogQueueSize)

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
	}()

	//Initialize log using configuration from "../conf/log.config"
	InitLog("../conf/log.config")

	log.Info(time.Now().String(), " Log server starting ... ")
	log.Info("Version: " + config.Version)

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
