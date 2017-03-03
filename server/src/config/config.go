package config

import (
	"github.com/spf13/viper"
)

const Version = "0.0.1"

//Sinker
type SinkerInfo struct {
	LogPath string `mapstructure:"log_path" json:"log_path"`
	LogFileMaxSize int64 `mapstructure:"log_file_max_size" json:"log_file_max_size"`
	LogFlushDuration uint32 `mapstructure:"log_flush_duration" json:"log_flush_duration"`
	LogQueueSize uint32 `mapstructure:"log_queue_size" json:"log_queue_size"`
}

//Kafka
type KafkaInfo struct {
	Broker []string `mapstructure:"broker" json:"broker"`
	Topic string `mapstructure:"topic" json:"topic"`
	Partitions []int32 `mapstructure:"partitions" json:"partitions"`
}

//Config sturcture
type Config struct {
	Sinker SinkerInfo `mapstructure:"sinker" json:"sinker"`
	Kafka KafkaInfo `mapstructure:"kafka" json:"kafka"`
}

//Global config
var globalConfig *Config

//New Config
func NewConfig() *Config {
	return &Config{
		Sinker:SinkerInfo{LogPath:"/tmp/dark_metrix/log", LogFileMaxSize:524288000, LogFlushDuration:3, LogQueueSize:5000},
		Kafka:KafkaInfo{Broker:make([]string, 0), Topic:"dark_metrix_log", Partitions:make([]int32, 0)},
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
