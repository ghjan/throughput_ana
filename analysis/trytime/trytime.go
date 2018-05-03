package main

import (
	"fmt"
	"strings"

	"time"
)

const (
	time1 = "2018-05-03+10:02:58++0000"
	time2 = "2018-05-03 10:02:58++0000"
)

func main() {

	//获取时间戳

	timestamp := time.Now().Unix()

	fmt.Println(timestamp)

	//格式化为字符串,tm为Time类型

	tm := time.Unix(timestamp, 0)

	fmt.Println(tm.Format("2006-01-02 03:04:05 PM"))

	fmt.Println(tm.Format("02/01/2006 15:04:05 PM"))

	//从字符串转为时间戳，第一个参数是格式，第二个是要转换的时间字符串
	fmt.Println("-------------")
	getNum("01/02/2006", "02/08/2015")
	getNum(getFormat("day"), time1)
	getNum(getFormat("hour"), time1)
	getNum(getFormat("min"), time1)
	getNum(getFormat("day"), time2)
	getNum(getFormat("hour"), time2)
	getNum(getFormat("min"), time2)
}

func getNum(format string, sTime string) {
	t := strings.Replace(sTime[:len(format)], "+", " ", -1)
	tm2, _ := time.Parse(format, t)
	fmt.Println(tm2.Unix())
}

func getFormat(sTimeType string) string {
	var item string
	switch sTimeType {
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
	return item
}
