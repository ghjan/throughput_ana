package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"time"

	"crypto/md5"
	"encoding/hex"
	"net/url"
	"strings"
	"github.com/mediocregopher/radix.v2/redis"
	"github.com/mgutz/str"
	"github.com/sirupsen/logrus"
	"strconv"
	"fmt"
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
	data  digData
	uid   string
	unode urlNode
}

type urlNode struct {
	unType string //nrlnode type  详情页 或者首页或者列表页
	unRid  int    //urlnode request id  Resource ID 资源ID
	unUrl  string //当前页面的url
	unTime string //当前页面访问时间
}

type storageBlock struct {
	counterType  string
	storageModel string
	unode        urlNode
}

var log = logrus.New()
var redisCli redis.Client

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
	redisCli, err := dial()
	if err != nil {
		log.Fatalln("Redis connect failed!")
	} else {
		defer redisCli.Close()
	}
}

func dial() (*redis.Client, error) {
	client, err := redis.DialTimeout("tcp", "127.0.0.1:6379", 10*time.Second)
	return client, err
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
		//uid
		//说明课程中模拟生成uid， md5(refer+ua)
		hasher := md5.New()
		hasher.Write([]byte(data.refer + data.ua))
		uid := hex.EncodeToString(hasher.Sum(nil))

		//很多解析的功能可以放到这里...
		uData := urlData{
			data:  data,
			uid:   uid,
			unode: formatUrl(data.url, data.time),
		}
		pvChannel <- uData
		uvChannel <- uData
	}
	return nil
}

func formatUrl(url, t string) urlNode {
	//一定从量大的着手 详情页 movie
	pos1 := str.IndexOf(url, HANDLE_MOVIE, 0)
	if pos1 != -1 {
		pos1 += len(HANDLE_MOVIE)
		pos2 := str.IndexOf(url, HANDLE_HTML, pos1)
		idStr := str.Substr(url, pos1, pos2-pos1)
		unID, _ := strconv.Atoi(idStr)
		return urlNode{
			"movie", unID, url, t}
	} else {
		pos1 = str.IndexOf(url, HANDLE_LIST, 0)
		if pos1 != -1 {
			pos1 += len(HANDLE_LIST)
			pos2 := str.IndexOf(url, HANDLE_HTML, pos1)
			idStr := str.Substr(url, pos1, pos2-pos1)
			unID, _ := strconv.Atoi(idStr)
			return urlNode{
				"list", unID, url, t}
		} else { //首页
			return urlNode{
				"home", 0, url, t}
		}
	}
}

func cutLogFetchData(logStr string) digData {
	logStr = strings.TrimSpace(logStr)
	pos1 := str.IndexOf(logStr, HANDLE_DIG, 0)
	if pos1 == -1 {
		return digData{}
	}
	pos1 += len(HANDLE_DIG)
	pos2 := str.IndexOf(logStr, " HTTP/", pos1)
	d := str.Substr(logStr, pos1, pos2-pos1)
	urlInfo, err := url.Parse("http://localhost/?" + d)
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

//pv 访问多少次
func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {
	for data := range pvChannel {
		sItem := storageBlock{"pv", "ZINCRBY", data.unode}
		storageChannel <- sItem
	}
}

//user 需要去重
func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock) {
	for data := range uvChannel {
		//HyperLoglog redis
		hyperLogLogKey := "uv_hpll_" + getTime(data.data.time, "day")
		fmt.Printf("hyperLogLogKey:%s, uid:%s\n", hyperLogLogKey, data.uid)
		ret, err := redisCli.Cmd("PFADD", hyperLogLogKey, data.uid, "EX", 86400).Int()

		if err != nil {
			log.Warningln("uvCounter check redis hyperloglog failed, ", err)
		}
		if ret != 1 { //说明已经存在了
			continue
		}
		sItem := storageBlock{"uv", "ZINCRBY", data.unode}
		storageChannel <- sItem
	}
}

func getTime(logTime, timeType string) string {
	var item string
	switch timeType {
	case "day":
		item = "2006-01-02"
		break
	case "hour":
		item = "2006-01-02 15"
		break
	case "min":
		item = "2006-01-02 15:04"
		break
	}
	t, _ := time.Parse(item, time.Now().Format(item))
	return strconv.FormatInt(t.Unix(), 10)
}

func dataStorage(storageChannel chan storageBlock) {

}
