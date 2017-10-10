package asyncLog

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	logFromWeb   = "web"
	logFromCli   = "cli"
	logChanCap   = 10000
	logWriterNum = 2
	logTimeOut   = 100 * time.Millisecond
	logFormat    = "[%s] [%s] [%s] [%s] [%s]\n"
)

//log level
const (
	//Info info
	Info = 1

	//Debug debug
	Debug = 2

	//Notice notice
	Notice = 3

	//Trace trace
	Trace = 4

	//Warning warning
	Warning = 5

	//Error error
	Error = 6
)

var logLeveName = map[int]string{
	1: "INFO",
	2: "DEBUG",
	3: "NOTICE",
	4: "TRACE",
	5: "WARNING",
	6: "ERROR",
}

var logChan = make(chan *logItem, logChanCap)

var initMutex = new(sync.Mutex)

var initStatus = false

type logItem struct {
	level    int
	protocol string
	logFile  string
	logMsg   string
	request  *http.Request
}

func init() {
	initMutex.Lock()
	defer initMutex.Unlock()
	if initStatus == false {
		initStatus = true
		for i := 0; i < logWriterNum; i++ {
			go createWriter(logChan)
		}
	}
}

func createWriter(lc chan *logItem) {
	defer func() {
		if err := recover(); err != nil {
			go createWriter(lc)
		}
	}()

	for {
		li := <-lc
		timeOutChan := make(chan int, 1)
		go func() {
			time.Sleep(logTimeOut)
			timeOutChan <- 1
		}()
		select {
		case <-timeOutChan:
			go createWriter(lc)
			runtime.Goexit()
		case <-do(li):
		}
	}
}

//write log to file
func do(li *logItem) chan int {
	file, _ := os.OpenFile(li.logFile, os.O_CREATE|os.O_APPEND, 0777)
	hostname, _ := os.Hostname()
	reqStr := ""
	if li.protocol == logFromWeb && li.request != nil {
		reqStr = li.request.RequestURI + " " + li.request.RemoteAddr
	}
	msg := fmt.Sprintf(logFormat, logLeveName[li.level], time.Now().Format("2006-01-02 15:04:05"), hostname, reqStr, strings.TrimSpace(li.logMsg))
	file.WriteString(msg)
	file.Close()
	c := make(chan int, 1)
	c <- 1
	return c
}

//WebLog log
func WebLog(req *http.Request, file string, log string, level int) {
	item := new(logItem)
	item.protocol = logFromWeb
	item.request = req
	item.logFile = file
	item.logMsg = log
	item.level = level
	logChan <- item
}

//CLILog log
func CLILog(file string, log string, level int) {
	item := new(logItem)
	item.protocol = logFromCli
	item.logFile = file
	item.logMsg = log
	item.level = level
	logChan <- item
}
