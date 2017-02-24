package main

import (
	"os"
	"os/signal"
	"time"

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

/*
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
*/

func main() {
	defer log.Flush()

	//Initialize log using configuration from "../conf/monitor_agent_log.config"
	err := InitLog("../conf/log_server_log.config")

	if err != nil {
		log.Warnf("Read config failed! error:%s", err)
		return
	}

	log.Info(time.Now().String(), "Starting log_server ... ")

	//Initialize the configuration from "../conf/monitor_agent_config.json"
	/*log.Info("Initialize monitor_agent configuration from ../conf/monitor_agent_config.json ...")
	config, err := InitConfig("../conf/monitor_agent_config.json")

	if err != nil {
		log.Warnf("Initialize monitor_agent configuration failed! error:%s", err)
		return
	}

	log.Info("Initialize monitor_agent configuration successed! config:", config)
	*/

	//Deal with signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, os.Kill)

	signalOccur := <-signalChannel

	log.Info("Signal occured, signal:", signalOccur.String())
}