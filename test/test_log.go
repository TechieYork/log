package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	log_proto "github.com/DarkMetrix/log/proto"
	"github.com/golang/protobuf/proto"
	"flag"
)

func HandleUdp(w *sync.WaitGroup, unix_domain_socket string, number_per_thread int, message string, level int, project string, service string) {
	conn, err := net.Dial("unixgram", unix_domain_socket)

	if nil != err {
		fmt.Println("Dail failed! error:" + err.Error())
		w.Done()
		return
	}

	defer conn.Close()

    begin_time := time.Now().Unix()

	//buffer := make([]byte, 1024)
    var log_package log_proto.LogPackage
    log_package.Project = proto.String(project)
	log_package.Service = proto.String(service)
	log_package.Level = proto.Uint32(uint32(level))
    //log_package.Log = []byte(*(proto.String("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789")))
	log_package.Log = []byte(*(proto.String(message)))

    data, err := proto.Marshal(&log_package)

	fmt.Print(log_package.String())

    if nil != err {
        fmt.Println("Marshal failed! err:" + err.Error())
    }

	for index := 0; index < number_per_thread; index++ {
		conn.Write(data)

        if index % 1000 == 0 {
            fmt.Printf("index:%d\r\n", index)
        }
	}

    end_time := time.Now().Unix()

    fmt.Printf("Test done -> using time:%d\r\n", end_time - begin_time)

	w.Done()
}

//Parse flag
var threads *int
var number_per_thread *int
var message *string
var unix_domain_socket *string
var project *string
var service *string
var level *int

func main() {
	fmt.Println("unix socket test client begin")

	threads = flag.Int("threads", 1, "Number of threads to send logs, default:1")
	number_per_thread = flag.Int("number_per_threads", 1, "Number of log to send each thread, default:1")
	message = flag.String("message", "test log for dm_log_agent", "Log content to send, default:\"test log for dm_log_agent\"")
	unix_domain_socket = flag.String("unix_domain_socket", "/var/tmp/net_log.sock", "The unix domain socket path, default:/var/tmp/net_log.sock")
	project = flag.String("project", "test_project", "The log project, default:test_project")
	service = flag.String("service", "test_service", "The log service, default:test_service")
	level = flag.Int("level", 1, "The log level, default:1 (1, 2, 3, 4 available)")

	flag.Parse()

	var w sync.WaitGroup
	w.Add(*threads)

	for i := 0; i < *threads; i++{
		go HandleUdp(&w, *unix_domain_socket, *number_per_thread, *message, *level, *project, *service)
	}

	w.Wait()
	fmt.Println("unix socket test client end")
}

