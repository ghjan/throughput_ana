package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type resource struct {
	url    string
	target string
	start  int
	end    int
}

func ruleResource() []resource {
	//主页
	r1 := resource{
		url:    "http://localhost/",
		target: "",
		start:  0,
		end:    0,
	}
	//列表页
	r2 := resource{
		url:    "http://localhost/list/{$id}.html",
		target: "{$id}",
		start:  1,
		end:    21,
	}
	//详情页
	r3 := resource{
		url:    "http://localhost/movie/{$id}.html",
		target: "{$id}",
		start:  1,
		end:    12924,
	}
	res := []resource{r1, r2, r3}
	return res
}

func buildUrl(res []resource) []string {
	var list []string
	for _, r := range res {
		if len(r.target) == 0 {
			list = append(list, r.url)
		} else {
			for i := r.start; i <= r.end; i++ {
				urlString := strings.Replace(r.url, r.target, strconv.Itoa(i), -1)
				list = append(list, urlString)
			}
		}
	}
	return list
}

func main() {
	total := flag.Int("total", 100, "how many rows created")
	filePath := flag.String("filePath", "/application/nginx/logs/dig.log", "dig log file path")
	flag.Parse()

	// 需要构造出真实的网站url集合
	res := ruleResource()
	list := buildUrl(res)
	// 按照要求， 生成 $total行日志内容， 源自于的这个集合
	logStr := ""
	for i := 0; i <= *total; i++ {
		currentUrl := list[randInt(0, len(list)-1)]
		referUrl := list[randInt(0, len(list)-1)]
		ua := UAs[randInt(0, len(UAs)-1)]
		logStr = logStr + makeLog(currentUrl, referUrl, ua) + "\n"
	}
	fd, _ := os.OpenFile(*filePath, os.O_RDWR|os.O_APPEND, 0644)
	fd.Write([]byte(logStr))
	defer fd.Close()
	fmt.Println("done.\n")
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

/*
localhost - - [02/May/2018:14:31:11 +0000]
"OPTIONS /dig?
time=2018%2F5%2F2+%E4%B8%8B%E5%8D%8810%3A30%3A51
&url=http%3A%2F%2Flocalhost%2Fmovie%2F12672.html
&refer=http%3A%2F%2Flocalhost%2Flist%2F4.html
&ua=Mozilla%2F5.0+(Windows+NT+6.1%3B+Win64%3B+x64)+AppleWebKit%2F537.36+(KHTML%2C+like+Gecko)+Chrome%2F66.0.3359.117+Safari%2F537.36 HTTP/1.1" 200 43 "-" "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.117 Safari/537.36"
"-"
*/
func makeLog(currentUrl, referUrl, ua string) string {
	u := url.Values{}
	time := time.Now()
	time1 := time.Format("02/Jan/2006:15:04:05 +0000")
	time2 := time.Format("2006-01-02 15:04:05 +0000")
	u.Set("time", time2)
	u.Set("url", currentUrl)
	u.Set("refer", referUrl)
	u.Set("ua", ua)
	paramsStr := u.Encode()
	logTemplate := `localhost - - [{$time}] "OPTIONS /dig?{$paramsStr} HTTP/1.1" 200 43 "-" {$ua} "-"`
	log := strings.Replace(logTemplate, "{$time}", time1, -1)
	log = strings.Replace(log, "{$paramsStr}", paramsStr, -1)
	//log = strings.Replace(log, "{$currentUrl}", currentUrl, -1)
	//log = strings.Replace(log, "{$referUrl}", referUrl, -1)
	log = strings.Replace(log, "{$ua}", ua, -1)
	return log
}

var UAs = []string{
	//浏览器User-Agent的详细信息
	//PC端：
	//safari 5.1 – MAC
	`User-Agent:Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50`,
	//safari 5.1 – Windows
	`User-Agent:Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50`,
	//IE 9.0
	`User-Agent:Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0;`,
	//IE 8.0
	`User-Agent:Mozilla/4.0 (compatible; MSIE 8.0; Windows NT 6.0; Trident/4.0)`,
	//IE 7.0
	`User-Agent:Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 6.0)`,
	//IE 6.0
	`User-Agent: Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)`,
	//Firefox 4.0.1 – MAC
	`User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.6; rv:2.0.1) Gecko/20100101 Firefox/4.0.1`,
	//Firefox 4.0.1 – Windows
	`User-Agent:Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1`,
	//Opera 11.11 – MAC
	`User-Agent:Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11`,
	//Opera 11.11 – Windows
	`User-Agent:Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11`,
	//Chrome 17.0 – MAC
	`User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11`,
	//傲游（Maxthon）
	`User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)`,
	//腾讯TT
	`User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; TencentTraveler 4.0)`,
	//世界之窗（The World） 2.x
	`User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)`,
	//世界之窗（The World） 3.x
	`User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)`,
	//搜狗浏览器 1.x
	`User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SE 2.X MetaSr 1.0; SE 2.X MetaSr 1.0; .NET CLR 2.0.50727; SE 2.X MetaSr 1.0)`,
	//360浏览器
	`User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)`,
	//Avant
	`User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Avant Browser)`,
	//Green Browser
	`User-Agent: Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1)`,
	//移动设备端：
	//safari iOS 4.33 – iPhone
	`User-Agent:Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5`,
	//safari iOS 4.33 – iPod Touch
	`User-Agent:Mozilla/5.0 (iPod; U; CPU iPhone OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5`,
	//safari iOS 4.33 – iPad
	`User-Agent:Mozilla/5.0 (iPad; U; CPU OS 4_3_3 like Mac OS X; en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5`,
	//Android N1
	`User-Agent: Mozilla/5.0 (Linux; U; Android 2.3.7; en-us; Nexus One Build/FRF91) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1`,
	//Android QQ浏览器 For android
	`User-Agent: MQQBrowser/26 Mozilla/5.0 (Linux; U; Android 2.3.7; zh-cn; MB200 Build/GRJ22; CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1`,
	//Android Opera Mobile
	`User-Agent: Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/build-1107180945; U; en-GB) Presto/2.8.149 Version/11.10`,
	//Android Pad Moto Xoom
	`User-Agent: Mozilla/5.0 (Linux; U; Android 3.0; en-us; Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13`,
	//BlackBerry
	`User-Agent: Mozilla/5.0 (BlackBerry; U; BlackBerry 9800; en) AppleWebKit/534.1+ (KHTML, like Gecko) Version/6.0.0.337 Mobile Safari/534.1+`,
	//WebOS HP Touchpad
	`User-Agent: Mozilla/5.0 (hp-tablet; Linux; hpwOS/3.0.0; U; en-US) AppleWebKit/534.6 (KHTML, like Gecko) wOSBrowser/233.70 Safari/534.6 TouchPad/1.0`,
	//Nokia N97
	`User-Agent: Mozilla/5.0 (SymbianOS/9.4; Series60/5.0 NokiaN97-1/20.0.019; Profile/MIDP-2.1 Configuration/CLDC-1.1) AppleWebKit/525 (KHTML, like Gecko) BrowserNG/7.1.18124`,
	//Windows Phone Mango
	`User-Agent: Mozilla/5.0 (compatible; MSIE 9.0; Windows Phone OS 7.5; Trident/5.0; IEMobile/9.0; HTC; Titan)`,
	//UC无
	`User-Agent: UCWEB7.0.2.37/28/999`,
	//UC标准
	`User-Agent: NOKIA5700/ UCWEB7.0.2.37/28/999`,
	//UCOpenwave
	`User-Agent: Openwave/ UCWEB7.0.2.37/28/999`,
	//UC Opera
	`User-Agent: Mozilla/4.0 (compatible; MSIE 6.0; ) Opera/UCWEB7.0.2.37/28/999`,
}
