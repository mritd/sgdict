package wordrank

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

func queryBaiduRank(word string) (int, error) {

	maxRetry := 3

retry:
	cli := client()
	resp, err := cli.R().Get(fmt.Sprintf(BAIDU_API, word))
	if err != nil {
		return 0, fmt.Errorf("[BaiduWordRank] query word [%s] rank failed: %s", word, err)
	}

	rex, _ := regexp.Compile("百度为您找到相关结果约(.*)个")
	res := rex.FindStringSubmatch(string(resp.Body()))
	if len(res) != 2 {
		if maxRetry > 0 {
			maxRetry--
			goto retry
		}
		return 0, fmt.Errorf("[BaiduWordRank] get word [%s] rank failed: \n%s", word, string(resp.Body()))
	}
	rank, err := strconv.Atoi(strings.ReplaceAll(res[1], ",", ""))
	if err != nil {
		return 0, fmt.Errorf("[BaiduWordRank] rank [%s] format failed: %s", res[1], err)
	}

	return rank, nil
}

func BaiduWorkRank() {
	info, err := os.Stat(FilePath)
	if err != nil {
		logrus.Fatal(err)
	}
	if info.IsDir() {
		logrus.Fatalf("%s is a dir", FilePath)
	}

	logrus.Infof("Pool size %d", PoolSize)

	pool, err := ants.NewPool(PoolSize, ants.WithPreAlloc(true))
	if err != nil {
		logrus.Fatal(err)
	}
	defer pool.Release()

	src, err := os.Open(FilePath)
	if err != nil {
		logrus.Errorf("open file [%s] failed: %s", FilePath, err)
		return
	}
	defer func() { _ = src.Close() }()

	//wordMap := make(map[string][]string, 1000000)
	wordCh := make(chan []string, 100)
	//var keys []string

	go func() {
		br := bufio.NewReader(src)
		for {
			s, err := br.ReadString('\n')
			if err != nil {
				break
			}
			s = strings.TrimSpace(s)
			_ = pool.Submit(func() {
				ss := strings.Split(s, "\t")
				if len(ss) != 3 {
					logrus.Errorf("dict format error: %s", s)
					return
				}

				rank, err := queryBaiduRank(ss[0])
				if err != nil {
					logrus.Error(err)
				} else {
					ss[2] = strconv.Itoa(rank)
				}
				wordCh <- ss
				logrus.Infof("processed %s", ss)
			})
		}
	}()

	outFile, err := os.OpenFile(FilePath+".rank", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() { _ = outFile.Close() }()
	for {
		select {
		case w := <-wordCh:
			s := strings.Join(w, "\t")
			_, err := fmt.Fprintln(outFile, strings.TrimSpace(s))
			if err != nil {
				logrus.Error(err)
			}
		case <-time.After(10 * time.Second):
			close(wordCh)
			goto done
		}
	}

done:

	//outFile, err := os.OpenFile(FilePath+".sort", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	//if err != nil {
	//	logrus.Fatal(err)
	//}
	//defer func() { _ = outFile.Close() }()
	//sort.Sort(pinyin.ByPinyin(keys))
	//
	//for _, k := range keys {
	//	s := strings.Join(wordMap[k], "\t")
	//	_, err := fmt.Fprintln(outFile, strings.TrimSpace(s))
	//	if err != nil {
	//		logrus.Error(err)
	//	}
	//}

}
