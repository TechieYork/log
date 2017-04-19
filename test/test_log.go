package main

import (
	"fmt"
	"sync"
	"time"
	"flag"

	net_log "github.com/DarkMetrix/log/client/golang"
)

func HandleUdp(w *sync.WaitGroup, number_per_thread int, message string, level uint32) {
	begin_time := time.Now().Unix()

	for index := 0; index < number_per_thread; index++ {
		err := net_log.SendNetLog(level, message)

		if err != nil {
			fmt.Println("Sendto failed! error:" + err.Error())
		}

		fmt.Printf("index:%d\r\n", index)
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

	//Init net log
	err := net_log.InitNetLog("test_project", "test_service")

	if err != nil {
		fmt.Println("Init net log failed! error:" + err.Error())
		return
	}

	var w sync.WaitGroup
	w.Add(*threads)

	for i := 0; i < *threads; i++{
		go HandleUdp(&w, *number_per_thread, *message, uint32(*level))
	}

	w.Wait()

	err = net_log.UninitNetLog()

	if err != nil {
		fmt.Println("Uninit net log failed! error:" + err.Error())
		return
	}

	fmt.Println("unix socket test client end")
}

