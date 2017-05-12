package collector

import (
	"os"
	"net"
	"time"
	"fmt"
	"errors"

	log_proto "github.com/DarkMetrix/log/proto"
	"github.com/DarkMetrix/log/agent/src/queue"

	log "github.com/cihub/seelog"
	"github.com/golang/protobuf/proto"
)

//Log collector to collect log from unix domain socket
type Collector struct {
	address string                      //Unix domain socket address
	conn net.Conn                       //Unix domain socket connection

	localIp string                      //Local ip

	logQueue *queue.LogQueue            //Log queue to buffer logs to send later
}

//New collector function
func NewCollector(address string, logQueue *queue.LogQueue) *Collector {
	return &Collector{
		address: address,
		conn: nil,
		logQueue: logQueue,
	}
}

//Get local machine ip
func (collector *Collector) getLocalIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

    if err != nil {
		log.Warn("Get local ip failed! err:" + err.Error())
	    return "", err
    }

    for _, addr := range addrs {
        if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
	            return ipnet.IP.String(), nil
            }
        }
    }

	return "", nil
}

//Run to open unix domain socket and collect
func (collector *Collector) Run() error {
	//Get local machine ip
	var err error

	collector.localIp, err = collector.getLocalIp()

	if err != nil {
		return err
	} else if collector.localIp == ""{
		return errors.New("Ip address that got it is empty!")
	}

	//Remove if unix domain socket exist
	_, err = os.Stat(collector.address)

	if err == nil {
		err = os.Remove(collector.address)

		if err != nil {
			log.Warn("Remove unix domain socket file failed! err:" + err.Error())
			return err
		}
	}

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

	//Set read buffer
	err = conn.SetReadBuffer(16 * 1024 * 1024)

	if nil != err {
		log.Warn("Set read buffer err:" + err.Error())
		return err
	}

    collector.conn = conn

	//Change unix domain socket file mode
	err = os.Chmod(collector.address, 0777)

	if err != nil {
		return err
	}

	//Change unix domain socket file to nobody
	err = os.Chown(collector.address, 99, 99)

	if err != nil {
		return err
	}

    //Begin recv local log
    go collector.Collect()

	return nil
}

//Close unix domain socket and remove socket file
func (collector *Collector) Close() error {
	//Close connection
	collector.conn.Close()

	//Remove unix domain socket file
	err := os.Remove(collector.address)

	if err != nil {
		return err
	}

	return nil
}

//Collect go routine function
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

				if err != nil {
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

		logPackage.Log = []byte(fmt.Sprintf("[%s]%s", collector.localIp, string(logPackage.GetLog())))

		if err != nil {
			log.Warn("Collector couldn't unmarshal received buffer! err:" + err.Error())
			continue
		}

		err = collector.logQueue.Push(&logPackage)

		if err != nil {
			log.Warn("Collector push log to queue failed! err:" + err.Error())
			continue
		}
	}

	return nil
}
