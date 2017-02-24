package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	log_proto "github.com/DarkMetrix/log/proto"
	"github.com/golang/protobuf/proto"
)

func HandleUdp(w *sync.WaitGroup) {
	conn, err := net.Dial("unixgram", "/var/tmp/dark_metrix_log.sock")

	if nil != err {
		fmt.Println("Dail failed! error:" + err.Error())
		w.Done()
		return
	}

	defer conn.Close()

    begin_time := time.Now().Unix()

	//buffer := make([]byte, 1024)
    var log_package log_proto.LogPackage
    log_package.Project = proto.String("test_project_go")
	log_package.Service = proto.String("test_service_go")
	log_package.Level = proto.Uint32(1)
    //log_package.Log = proto.String("01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789")
	log_package.Log = []byte(*(proto.String("test log golang")))

    data, err := proto.Marshal(&log_package)

	fmt.Print(log_package.String())

    if nil != err {
        fmt.Println("Marshal failed! err:" + err.Error())
    }

	for index := 0; index < 3; index++ {
		conn.Write(data)
		//time.Sleep(10 * time.Millisecond)
		//fmt.Printf("Test the %d time and send -> %s\r\n", index, log_package.GetLog())
		//conn.Read(buffer)

        if index % 1000000 == 0 {
            fmt.Printf("index:%d\r\n", index)
        }
	}

    end_time := time.Now().Unix()

    fmt.Printf("Test done -> using time:%d\r\n", end_time - begin_time)

	w.Done()
}

func main() {
	fmt.Println("unix socket test client begin")

	var w sync.WaitGroup
	w.Add(1)

	go HandleUdp(&w)

	w.Wait()
	fmt.Println("unix socket test client end")
}

