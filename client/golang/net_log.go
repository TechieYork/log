package golang

import (
	"strings"
	"errors"
	"syscall"

	"github.com/golang/protobuf/proto"

	"github.com/DarkMetrix/log/proto"
)

//Global variables setting
var (
	netLogPath = "/var/tmp/net_log.sock"
	netLogAddr = syscall.SockaddrUnix{Name:netLogPath}
	netLogSocket = 0
	netLogSendBufferLen = 4 * 1024 * 1024
	netLogMaxLen = netLogSendBufferLen - 10240
	netLogIsInit = false
	netLogProject = ""
	netLogService = ""

	returnTag = "[@RETURN]"
	newlineTag = "[@NEWLINE]"
)

//Init net log using project name and service name
func InitNetLog (project string, service string) error {
	var err error

	//Check is init or not
	if netLogIsInit {
		return nil
	}

	//Check params
	if len(project) > 128 || len(service) > 128 || len(project) == 0 || len(service) == 0 {
		return errors.New("Project or service params invalid!")
	}

	//Init socket and address
	netLogSocket, err = syscall.Socket(syscall.AF_UNIX, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)

	if err != nil {
		return err
	}

	netLogProject = project
	netLogService = service

	netLogIsInit = true

	return nil
}

//Uninit net log
func UninitNetLog () error {
	if !netLogIsInit {
		return nil
	}

	err := syscall.Close(netLogSocket)

	if err != nil {
		return err
	}

	netLogSocket = 0

	netLogIsInit = false
	return nil
}

//Send net log
//level include INFO, WARNING, ERROR and FATAL(INFO=0, WARNING=1, ERROR=2, FATAL=3)
//log must less than 4M - 10K
func SendNetLog (level uint32, log string) error {
	if !netLogIsInit {
		return errors.New("Net log not init!")
	}

	if len(log) > netLogMaxLen {
		return errors.New("Log size too large!")
	}

	//buffer := make([]byte, 1024)
	var logPackage log_proto.LogPackage
	logPackage.Project = proto.String(netLogProject)
	logPackage.Service = proto.String(netLogService)
	logPackage.Level = proto.Uint32(level)
	logPackage.Log = []byte(*(proto.String(log)))

	data, err := proto.Marshal(&logPackage)

	if err != nil {
		return err
	}

	err = syscall.Sendto(netLogSocket, data,0, &netLogAddr)

	if err != nil {
		return err
	}

	return nil
}

//Escape log
func EscapeLog (log string) string {
	escapeLog := strings.Replace(log, "\r", returnTag, 0)
	escapeLog = strings.Replace(escapeLog, "\n", newlineTag, 0)

	return escapeLog
}

//Unescape log
func UnescapeLog (log string) string {
	unescapeLog := strings.Replace(log, returnTag, "\r", 0)
	unescapeLog = strings.Replace(unescapeLog, newlineTag, "\n", 0)

	return unescapeLog
}
