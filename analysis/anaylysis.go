package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"strings"
	"github.com/mgutz/str"
	"net/url"
	"crypto/md5"
	"encoding/hex"
)

const (
	HANDLE_DIG   = ` /dig?`
	HANDLE_MOVIE = "/movie/"
	HANDLE_LIST  = "/list/"
	HANDLE_HTML  = ".html"
)

var (
	logFilePath = flag.String("logFilePath", "/data/nginx/logs/dig.log", "log file path")
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
	log.Infof("Params: logFilePath=%s， routineNum=%d", *logFilePath, params.routineNum)

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
	fd, err := os.Open(params.logFilePath)
	if err != nil {
		log.Warningf("readFileLineByLine cannot open file:%s", params.logFilePath)
	}
	defer fd.Close()

	count := 0
	bufferRead := bufio.NewReader(fd)
	for {
		line, err := bufferRead.ReadString('\n')
		logChannel <- line
		log.Infof("readFileLineByLine line:%d", count)
		count++
		if count%(1000*params.routineNum) == 0 {
			log.Infof("readFileLineByLine line:%d", count)
		}
		if err != nil {
			if err == io.EOF {
				time.Sleep(3 * time.Second)
				log.Infof("readFileLineByLine wait, readline:%d", count)
			} else {
				log.Warningf("readFileLineByLine read error, %v", err)
			}
		}
	}
}

func logConsumer(logChannel chan string, pvChannel, uvChannel chan urlData) error {
	for logStr := range logChannel {
		//切割日志字符串，扣出打点数据
		data := cutLogFetchData(logStr)
		log.Infof("logConsumer, data, refer:%s, url:%s", data.refer, data.url)
		//uid
		//说明课程中模拟生成uid， md5(refer+ua)
		hasher := md5.New()
		hasher.Write([]byte(data.refer + data.ua))
		uid := hex.EncodeToString(hasher.Sum(nil))

		//很多解析的功能可以放到这里...
		uData := urlData{
			data: data,
			uid:  uid,
		}
		log.Infoln(uData)
		pvChannel <- uData
		uvChannel <- uData
	}
	return nil
}

func cutLogFetchData(logStr string) digData {
	logStr = strings.TrimSpace(logStr)
	pos1 := str.IndexOf(logStr, HANDLE_DIG, 0)
	if pos1 == -1 {
		return digData{}
	}
	pos2 := str.IndexOf(logStr, " HTTP/", pos1)
	d := str.Substr(logStr, pos1, pos2-pos1)
	urlInfo, err := url.Parse("http://localhost?" + d)
	if err != nil {
		return digData{}
	}
	data := urlInfo.Query()
	return digData{
		time:  data.Get("time"),
		url:   data.Get("url"),
		refer: data.Get("refer"),
		ua:    data.Get("ua"),
	}

}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {

}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock) {

}

func dataStorage(storageChannel chan storageBlock) {

}
