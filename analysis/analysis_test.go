package main

import "testing"

const (
	logstr  = `localhost - - [02/May/2018:14:31:11 +0000] "OPTIONS /dig?time=2018%2F5%2F2+%E4%B8%8B%E5%8D%8810%3A30%3A51&url=http%3A%2F%2Flocalhost%2Fmovie%2F12672.html&refer=http%3A%2F%2Flocalhost%2Flist%2F4.html&ua=Mozilla%2F5.0+(Windows+NT+6.1%3B+Win64%3B+x64)+AppleWebKit%2F537.36+(KHTML%2C+like+Gecko)+Chrome%2F66.0.3359.117+Safari%2F537.36 HTTP/1.1" 200 43 "-" "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.117 Safari/537.36"`
	logstr2 = `localhost - - [03/May/2018:10:02:58 +0000] "OPTIONS /dig?refer=http%3A%2F%2Flocalhost%2Fmovie%2F10698.html&time=2018-05-03+10%3A02%3A58+%2B0000&ua=User-Agent%3A+Opera%2F9.80+%28Android+2.3.4%3B+Linux%3B+Opera+Mobi%2Fbuild-1107180945%3B+U%3B+en-GB%29+Presto%2F2.8.149+Version%2F11.10&url=http%3A%2F%2Flocalhost%2Fmovie%2F7165.html HTTP/1.1" 200 43 "-" User-Agent: Opera/9.80 (Android 2.3.4; Linux; Opera Mobi/build-1107180945; U; en-GB) Presto/2.8.149 Version/11.10 "-"`
)

func TestCutLogFetchData(t *testing.T) {
	digData := cutLogFetchData(logstr)
	t.Log(digData.refer)
	t.Log(digData.url)

	digData = cutLogFetchData(logstr2)
	t.Log(digData.refer)
	t.Log(digData.url)

}
