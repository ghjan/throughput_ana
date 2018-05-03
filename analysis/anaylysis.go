package main

import (
	"flag"
	"time"
	"github.com/sirupsen/logrus"
	"os"
)

var (
	logFilePath = flag.String("logFilePath", "/application/nginx/logs/dig.log", "log file path")
	routineNum  = flag.Int("routineNum", 5, "consumer number by goroutine")
	l           = flag.String("l", "/tmp/log/analysis.log", "runtime log")
)

type cmdParams struct {
	logFilePath string
	routineNum  int
}
type digData struct {
	time  string
	url   string
	refer string
	ua    string
}
type urlData struct {
	data digData
	uid  string
}

type urlNode struct {
	unType string
}

type storageBlock struct {
	counterType  string
	storageModel string
	unode        urlNode
}

var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	//获取参数
	flag.Parse()
	params := cmdParams{*logFilePath, *routineNum}
	//打日志
	logFd, err := os.OpenFile(*l, os.O_CREATE|os.O_WRONLY, 0655)
	if err == nil {
		log.Out = logFd
		defer logFd.Close()
	}
	log.Infoln("Exec start.")
	log.Infoln("Params: logFilePath=%s， routineNum=%d", *logFilePath, params.routineNum)

	//初始化一些channel,用于数据传递
	var logChannel = make(chan string, params.routineNum*3)
	var pvChannel = make(chan urlData, params.routineNum)
	var uvChannel = make(chan urlData, params.routineNum)
	var storageChannel = make(chan storageBlock, params.routineNum)
	//日志消费者
	go readFileLineByLine(params, logChannel)
	//创建一组日志处理
	for i := 0; i < params.routineNum; i++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}
	//创建PV UV 统计器
	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel)
	//创建存储器
	go dataStorage(storageChannel)

	time.Sleep(1000 * time.Minute)
}

func readFileLineByLine(params cmdParams, logChannel chan string) {

}

func logConsumer(logChannel chan string, pvChannel, uvChannel chan urlData) {

}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {

}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock) {

}

func dataStorage(storageChannel chan storageBlock) {

}