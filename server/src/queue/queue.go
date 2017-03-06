package queue

import (
	"errors"
	"time"

	log_proto "github.com/DarkMetrix/log/proto"
)

//Log queue to buffer logs
type LogQueue struct {
	queueChannel chan *log_proto.LogPackage
}

//New log queue function
func NewLogQueue(bufferSize uint32) *LogQueue {
	return &LogQueue{
		queueChannel: make(chan *log_proto.LogPackage, bufferSize),
	}
}

//Push log to queue
func (queue *LogQueue) Push(item* log_proto.LogPackage) error {
	select {
	case queue.queueChannel <- item:
		return nil
	default:
		return errors.New("Channel full")
	}
}

//Pop log from queue
func (queue *LogQueue) Pop(ms time.Duration) (*log_proto.LogPackage, error) {
	select {
	case item, ok := <- queue.queueChannel:
		if !ok {
			return nil, errors.New("Channel closed!")
		}

		return item, nil
	case <- time.After(ms):
		return nil, errors.New("Channel pop timeout!")
	}
}

//Get log channel
func (queue *LogQueue) Chan() chan *log_proto.LogPackage {
	return queue.queueChannel
}
