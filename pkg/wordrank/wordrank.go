package wordrank

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

const (
	BAIDU_API       = "https://www.baidu.com/s?wd=%s"
	GOOGLE_API      = "https://www.google.com/search?q=%s"
	UA              = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"
	ACCEPT          = "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9"
	ACCEPT_LANGUAGE = "zh-CN,zh;q=0.9,en;q=0.8"
	ACCEPT_ENCODING = "gzip, deflate"
)

var (
	FilePath         string
	Proxy            string
	PoolSize         int
	Timeout          time.Duration
	RetryCount       int
	RetryMaxWaitTime time.Duration
)

func client() *resty.Client {
	return resty.New().
		SetLogger(logrus.StandardLogger()).
		SetProxy(Proxy).
		SetTimeout(Timeout).
		SetRetryCount(RetryCount).
		SetRetryMaxWaitTime(RetryMaxWaitTime).
		SetHeader("User-Agent", UA).
		SetHeader("Accept", ACCEPT).
		SetHeader("Accept-Language", ACCEPT_LANGUAGE).
		SetHeader("Accept-Encoding", ACCEPT_ENCODING)

}
