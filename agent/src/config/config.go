package config

import (
	"github.com/spf13/viper"
)

const Version = "0.0.1"

//Collector
type CollectorInfo struct {
	UnixDomainSocket string `mapstructure:"unix_domain_socket" json:"unix_domain_socket"`
	LogQueueSize uint32 `mapstructure:"log_queue_size" json:"log_queue_size"`
}

//Kafka
type KafkaInfo struct {
	Broker []string `mapstructure:"broker" json:"broker"`
	Topic string `mapstructure:"topic" json:"topic"`
	CompressCodec string `mapstructure:"compress_codec" json:"compress_codec"`
}

//Config sturcture
type Config struct {
	Collector CollectorInfo `mapstructure:"collector" json:"collector"`
	Kafka KafkaInfo `mapstructure:"kafka" json:"kafka"`
}

//Global config
var globalConfig *Config

//New Config
func NewConfig() *Config {
	return &Config{
		Collector:CollectorInfo{UnixDomainSocket:"/var/tmp/net_log.sock", LogQueueSize:20000},
		Kafka:KafkaInfo{Broker:make([]string, 0), Topic:"net_log", CompressCodec:"none"},
	}
}

//Get singleton config
func GetConfig() *Config {
	if globalConfig == nil {
		globalConfig = NewConfig()
	}

	return globalConfig
}

//Init config from json file
func (config *Config) Init (path string) error {
	//Set viper setting
	viper.SetConfigType("json")
	viper.SetConfigFile(path)
	viper.AddConfigPath("../conf/")

	//Read in config
	err := viper.ReadInConfig()

	if err != nil {
		return err
	}

	//Unmarshal config
	err = viper.Unmarshal(config)

	if err != nil {
		return err
	}

	return nil
}
