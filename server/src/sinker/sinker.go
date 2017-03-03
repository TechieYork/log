package sinker

import (
	"fmt"
	"os"
	"path"
	"bufio"
	"time"

	"github.com/DarkMetrix/log/server/src/queue"
	"github.com/DarkMetrix/log/proto"

	log "github.com/cihub/seelog"
)

//Sink worker
type SinkWorker struct {
	filePath string                     //File path
	fileName string                     //File name
	fileSize int64                      //Log file size currently
	fileMaxSize int64                   //Max size of log file in bytes
	flushDuration uint32                //Flush called duration in seconds

	logQueue *queue.LogQueue            //Log queue of specified service to buffer log which to sink later

	file *os.File                       //Log file
	writer *bufio.Writer                //Log buffer writer
}

//New sink worker function
func NewSinkWorker(filePath string, fileName string, fileMaxSize int64, flushDuration uint32) *SinkWorker {
	return &SinkWorker{
		filePath: filePath,
		fileName: fileName,
		fileSize: 0,
		fileMaxSize: fileMaxSize,
		flushDuration: flushDuration,
		logQueue: queue.NewLogQueue(5000),
		file: nil,
		writer: nil,
	}
}

//Get log queue
func (worker *SinkWorker) LogQueue() *queue.LogQueue {
	return worker.logQueue
}

//Init
func (worker *SinkWorker) initFile() error {
	//Open file
	var err error
    worker.file, err = os.OpenFile(path.Join(worker.filePath, worker.fileName + ".log." + time.Now().Format("20060102T150405")), os.O_RDWR | os.O_CREATE | os.O_APPEND, os.ModePerm)

    if nil != err {
	    log.Warnf("Worker open file failed! file name:%s, err:%s", worker.fileName, err.Error())
	    return err
    }

    fileInfo, err := worker.file.Stat()

    if nil != err {
	    log.Warnf("Worker get file stat failed! file name:%s, err:%s", worker.fileName, err.Error())
	    return err
    }

    worker.fileSize = fileInfo.Size()
    worker.writer = bufio.NewWriterSize(worker.file, 1 * 1024 * 1024)

	return nil
}

//Run to sink log to file
func (worker *SinkWorker) Run() error {
	//Check dir
	_, err := os.Stat(worker.filePath)

    if os.IsNotExist(err) {
        err = os.MkdirAll(worker.filePath, os.ModePerm)

        if nil != err {
	        log.Warnf("Create dir failed! path:%s, err:%s", worker.filePath, err.Error())
	        return err
        }

	    log.Info("Create dir success! path:" + worker.filePath)
    }

	//Init file
	err = worker.initFile()

	if err != nil {
		return err
	}

	go worker.process()

	return nil
}

//Close function
func (worker *SinkWorker) Close() error {
	if worker.writer != nil {
		err := worker.writer.Flush()

		if err != nil {
			log.Warnf("Worker flush failed! file name:%s, err:%s", worker.fileName, err)
		}

		worker.writer = nil
	}

	if worker.file != nil {
		err := worker.file.Close()

		if err != nil {
			log.Warnf("Worker file close failed! file name:%s, err:%s", worker.fileName, err)
		}

		worker.file = nil
    }

	return nil
}

//Generate log format
func (worker *SinkWorker) generateLog(logPackage *log_proto.LogPackage) string {
	var level string

	switch logPackage.GetLevel() {
	case 0:
		level = "INFO"
	case 1:
		level = "WARNING"
	case 2:
		level = "ERROR"
	case 3:
		level = "FATAL"
	default:
		level = "UNKNOWN"
	}
	
	return fmt.Sprintf("[%s][%s][%s][%s]%s\r\n", level, time.Now().Format("2006-01-02T15:04:05"), logPackage.GetProject(), logPackage.GetService(), string(logPackage.Log))
}

//Write log
func (worker *SinkWorker) writeLog(logContent string) error {
	begin := 0

	for {
		writeLen, err := worker.writer.Write([]byte(logContent[begin:]))

		begin += writeLen

		if begin >= len(logContent) {
			break
		}

		if err != nil {
			log.Warn("worker bufio write failed! err:" + err.Error())
			return err
		}
	}

	return nil
}

//Process go routine function
func (worker *SinkWorker) process() error {
	//Loop to pop log from queue
	for {
		select {
		//Get log
		case logPackage := <- worker.logQueue.Chan():

			//Generate log
			logContent := worker.generateLog(logPackage)

			//Check file size
			if int64(len(logContent)) + worker.fileSize > worker.fileMaxSize {
				worker.Close()

				for {
					err := worker.initFile()

					if err != nil {
						time.Sleep(5 * time.Second)
						continue
					}

					break
				}
			}

			//Check buffer is full or not
			if len(logContent) > worker.writer.Available() {
				worker.writer.Flush()
			}

			//Write log
			for {
				err := worker.writeLog(logContent)

				if err != nil {
					worker.Close()

					for {
						err := worker.initFile()

						if err != nil {
							time.Sleep(5 * time.Second)
							continue
						}

						break
					}
				}

				break
			}

			worker.fileSize += int64(len(logContent))
		//Flush after duration
		case <- time.After(time.Duration(worker.flushDuration) * time.Second):
			worker.writer.Flush()
		}
	}
}

//Log sinker
type Sinker struct {
	path string                         //Directory to save log files
	fileMaxSize int64                   //Max size of log file in bytes
	flushDuration uint32                //Flush called duration in seconds

	logQueue *queue.LogQueue            //Log queue to buffer log which to sink later

	sinkWorkers map[string]*SinkWorker  //Sink worker map
}

//New sinker function
func NewSinker(path string, fileMaxSize int64, flushDuration uint32, logQueue *queue.LogQueue) *Sinker {
	return &Sinker{
		path: path,
		fileMaxSize: fileMaxSize,
		flushDuration: flushDuration,
		logQueue: logQueue,
		sinkWorkers: make(map[string]*SinkWorker),
	}
}

//Run to sink
func (sinker *Sinker) Run() error {
	//Check dir
	_, err := os.Stat(sinker.path)

    if os.IsNotExist(err) {
        err = os.MkdirAll(sinker.path, os.ModePerm)

        if nil != err {
	        log.Warnf("Create dir failed! path:%s, err:%s", sinker.path, err.Error())
	        return err
        }

	    log.Info("Create dir success! path:" + sinker.path)
    }

	//Begin to sink
	go sinker.sink()

	return nil
}

//Close and flush file buffer
func (sinker *Sinker) Close() error {
	return nil
}

//Generate file key
func (sinker *Sinker) generateFileKey(logPackage *log_proto.LogPackage) string {
	return path.Join(sinker.path, logPackage.GetProject(), logPackage.GetService())
}

//Sink go routine function
func (sinker *Sinker) sink() error {
	//Loop to pop log from log queue and dispatch log to worker according to file name
	for {
		//Pop log
		logPackage, err := sinker.logQueue.Pop(time.Millisecond * 10)

		if err != nil {
			//log.Warn("Pop log from queue failed(empty)! err:" + err.Error())
			continue
		}

		//File key, used to identify in map eg:/data/log/project/service
		fileKey := sinker.generateFileKey(logPackage)

		worker, ok := sinker.sinkWorkers[fileKey]

		if !ok {
			log.Info("New sink worker, file key:" + fileKey)

			//Check dir
			dir := path.Join(sinker.path, logPackage.GetProject())
			_, err := os.Stat(dir)

			if os.IsNotExist(err) {
				err = os.MkdirAll(dir, os.ModePerm)

				if nil != err {
					log.Warnf("Create dir failed! path:%s, err:%s", sinker.path, err.Error())
					return err
				}

				log.Info("Create dir success! path:" + sinker.path)
			}

			//Init new sink worker
			worker = NewSinkWorker(path.Join(sinker.path, logPackage.GetProject()), logPackage.GetService(), sinker.fileMaxSize, sinker.flushDuration)

			err = worker.Run()

			if err != nil {
				log.Warn("Sinker create new worker failed! err:" + err.Error())
				continue
			}

			sinker.sinkWorkers[fileKey] = worker
		}

		err = worker.LogQueue().Push(logPackage)

		if err != nil {
			//log.Warn("Sinker push worker log queue failed! err:" + err.Error())
			continue
		}
	}
}
