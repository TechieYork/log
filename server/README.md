# DarkMetrix log server

## Introduction

DarkMetrix log server is an application that consume logs from Kafka and sink them to local file system as log files.

## Kafka requirements

* **Topic:** server need to consume logs from Kafka, so a log topic is needed. e.g.:"net_log"
* **Partitions:** In case the performance and space on one host isn't enough, you could create as many partitions as you need. So different host could consume different partitions.
* **Max message size:** you will need to set the max message size of Kafka, cause the max size of log message is 4M(Depending on the Linux kernel)

## Configuration

##### config.json

```json
{
    "sinker":
    {
        "log_path":"/tmp/dark_metrix/log/",
        "log_file_max_size":524288000,
        "log_flush_duration":3,
        "log_queue_size":10000
    },
    "kafka":
    {
        "broker":["localhost:9092"],
        "topic":"net_log",
        "partitions":[0, 1, 2, 3, 4, 5, 6, 7, 8, 9]
    }
}
```

- **sinker.log_path:** the path to save log files.
- **sinker.log_file_max_size:** the max size of one log file.
- **sinker.log_flush_duration:** the duration to flush log buffer to disk.
- **sinker.log_queue_size:** the size of log queue to buffer logs which would be sinked.
- **kafka.broker:** the Kafka broker list.
- **kafka.topic:** the Kafka topic which the log would be consumed.
- **kafka.partitions:** the partitions of kafka.topic which to be consumed.

##### log.config

See [cihub/seelog](https://github.com/cihub/seelog) to get more information.