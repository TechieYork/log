package collector

import (
	"os"
	"net"
	"time"

	log_proto "github.com/DarkMetrix/log/proto"
	"github.com/DarkMetrix/log/agent/src/queue"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
)

type Collector struct {
	address string
	conn net.Conn

	logQueue *queue.LogQueue
}

func NewCollector(address string, logQueue *queue.LogQueue) *Collector {
	return &Collector{
		address: address,
		conn: nil,
		logQueue: logQueue,
	}
}

func (collector *Collector) Run() error {
	//Initial local unix domain socket to recv log
	unixAddr, err := net.ResolveUnixAddr("unixgram", collector.address)

	if nil != err {
		log.Warn("Resolve unix socket addr err:" + err.Error())
		return err
	}

    //Listen unix domain socket and get connection
	conn, err := net.ListenUnixgram("unixgram", unixAddr)

	if nil != err {
		log.Warn("Listen unix socket addr err:" + err.Error())
		return err
	}

    collector.conn = conn

    //Begin recv local log
    go collector.Collect()

	return nil
}

//Close function
func (collector *Collector) Close() error {
	//Close conn
	collector.conn.Close()

	//Remove unix domain socket file
	err := os.Remove(collector.address)

	if err != nil {
		return err
	}

	return nil
}

func (collector *Collector) Collect() error {
	buffer := make([]byte, 4*1024*1024)

	//Loop to recv message
	for {
		len, err := collector.conn.Read(buffer)

		//If error occured close the connection and reconnect
		if err != nil {
			for {
				err := collector.Close()

				if err != nil {
					log.Warn("Collector close failed! err:" + err.Error())
					time.Sleep(5 * time.Second)

					continue
				}

				err = collector.Run()

				if nil != err {
					log.Warn("Collector restart failed! err:" + err.Error())
				    time.Sleep(5 * time.Second)

					continue
				}

				break
			}
		}

		//Put the log into queue
		var logPackage log_proto.LogPackage

		err = proto.Unmarshal(buffer[0:len], &logPackage)

		if err != nil {
			log.Warn("Collector couldn't unmarshal received buffer! err:" + err.Error())
			continue
		}

		collector.logQueue.Push(&logPackage)

		log.Info("log received:" + logPackage.String())
	}

	return nil
}
