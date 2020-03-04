package wordrank

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

func queryBaiduRank(word string) (int, error) {

	cli := client()
	resp, err := cli.R().Get(fmt.Sprintf(BAIDU_API, word))
	if err != nil {
		return 0, fmt.Errorf("[BaiduWordRank] query word [%s] rank failed: %s", word, err)
	}

	rex, _ := regexp.Compile("百度为您找到相关结果约(.*)个")
	res := rex.FindStringSubmatch(string(resp.Body()))
	if len(res) != 2 {
		return 0, fmt.Errorf("[BaiduWordRank] get word rank failed")
	}
	rank, err := strconv.Atoi(strings.ReplaceAll(res[1], ",", ""))
	if err != nil {
		return 0, fmt.Errorf("[BaiduWordRank] rank [%s] format failed: %s", res[1], err)
	}

	return rank, nil
}

func BaiduWorkRank() {
	info, err := os.Stat(BaseDir)
	if err != nil {
		logrus.Fatal(err)
	}
	if !info.IsDir() {
		logrus.Fatalf("%s is not a dir", BaseDir)
	}

	var count int
	_ = filepath.Walk(BaseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			count++
		}
		return nil
	})

	var wg sync.WaitGroup
	wg.Add(count)

	pool, err := ants.NewPool(100, ants.WithPreAlloc(true))
	if err != nil {
		logrus.Fatal(err)
	}
	defer pool.Release()

	_ = filepath.Walk(BaseDir, func(path string, info os.FileInfo, err error) error {

		// skip dir
		if info.IsDir() {
			return nil
		}

		err = pool.Submit(func() {
			defer wg.Done()

			src, err := os.Open(path)
			if err != nil {
				logrus.Errorf("processing file [%s] failed: %s", path, err)
				return
			}
			defer func() { _ = src.Close() }()

			dst, err := os.OpenFile(path+".rank", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
			if err != nil {
				logrus.Errorf("processing file [%s] failed: %s", path, err)
				return
			}
			defer func() { _ = dst.Close() }()

			br := bufio.NewReader(src)
			for {
				s, err := br.ReadString('\n')
				if err != nil {
					break
				}
				ss := strings.Split(s, "\t")
				if len(ss) != 3 {
					break
				}
				rank, err := queryBaiduRank(ss[0])
				if err != nil {
					logrus.Error(err)
				} else {
					ss[2] = strconv.Itoa(rank)
				}

				_, err = fmt.Fprintln(dst, strings.Join(ss, "\t"))
				if err != nil {
					logrus.Error(err)
				}
			}

			logrus.Infof("[BaiduWorkRank] file %s processed", path)

		})

		return nil
	})
}
