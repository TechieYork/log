# DarkMetrix Log

A fast network logging system powered by Apache Kafka.



## Architecture
![image](https://github.com/DarkMetrix/log/blob/master/doc/archetecture.png)


## Quick start

#### Preperation

An Apache Kafka need to be installed previously and create a topic.

```shell
./kafka-topics.sh --zookeeper localhost:2181 --create --replication-factor 1 --partitions 10 --topic net_log
```

See [DarkMetrix/log/agent](https://github.com/DarkMetrix/log/blob/master/agent/README.md) or [DarkMetrix/log/server](https://github.com/DarkMetrix/log/blob/master/server/README.md) to get more information about Kafka setting.



#### Install from source

```shell
git clone git@github.com:DarkMetrix/log.git 
```

#### Build agent & run (You will need a root user)

```shell
$cd DarkMetrix/log/agent/src
$go build -o ../bin/dm_log_agent
$../admin/start.sh
```

#### Build server & run

```shell
$cd DarkMetrix/log/server/src
$go build -o ../bin/dm_log_server
$../admin/start.sh
```



## Configuration

All configuration file is in json.

#### agent

See [DarkMetrix/log/agent](https://github.com/DarkMetrix/log/blob/master/agent/README.md) to get more information.

#### server

See [DarkMetrix/log/server](https://github.com/DarkMetrix/log/blob/master/server/README.md) to get more information.



## Lisense

#### DarkMetrix Log

MIT license

#### Dependencies

* github.com/cihub/seelog [BSD License](https://github.com/cihub/seelog/blob/master/LICENSE.txt)
* github.com/golang/protobuf/proto [BSD License](https://github.com/golang/protobuf/blob/master/LICENSE)
* github.com/spf13/viper [MIT License](https://github.com/spf13/viper/blob/master/LICENSE)
* github.com/Shopify/sarama [MIT License](https://github.com/Shopify/sarama/blob/master/LICENSE)
