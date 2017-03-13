# DarkMetrix log agent

## Introduction

DarkMetrix log agent is an application that collect logs from unix domain socket(Datagram) and produce them to Kafka.

## Kafka requirements

- **Topic:** server need to consume logs from Kafka, so a log topic is needed. e.g.:"net_log"
- **Partitions:** In case the performance and space on one host isn't enough, you could create as many partitions as you need. So different host could consume different partitions.
- **Max message size:** you will need to set the max message size of Kafka, cause the max size of log message is 4M(Depending on the Linux kernel)

## OS requirements

* **POSIX:** Since we are using unix domain socket to collect logs, the host machine system should support POSIX.
* **Kernel version:** The kernel version decides the max message size which unix domain socket(Datagram) could send or receive. The 2.6.32 and above support maximum 4M message size.And you may need to set some kernel params if needed.

## Configuration

##### config.json

```json
{
    "collector":
    {
        "unix_domain_socket":"/var/tmp/net_log.sock",
        "log_queue_size":20000
    },
    "kafka":
    {
        "broker":["localhost:9092"],
        "topic":"net_log",
        "compress_codec":"none"
    }
}
```

- **collector.unix_domain_socket:** the unix domain socket path.
- **collector.log_queue_size:** the size of log queue to buffer logs which would be sent.
- **kafka.broker:** the Kafka broker list.
- **kafka.topic:** the Kafka topic which the log would be produced.
- **kafka.compress_codec:** the compression setting, support "**gzip**", "**snappy**", "**lz4**" and "**none**".

##### log.config

See [cihub/seelog](https://github.com/cihub/seelog) to get more information.



