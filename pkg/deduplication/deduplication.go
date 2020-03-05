package deduplication

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/mritd/sgdict/pkg/pinyin"

	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
)

var (
	BaseDir string
	OutFile string
)

func deduplication() {
	info, err := os.Stat(BaseDir)
	if err != nil {
		logrus.Fatal(err)
	}
	if !info.IsDir() {
		logrus.Fatalf("%s is not a dir", BaseDir)
	}

	var count int
	_ = filepath.Walk(BaseDir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(path, ".rime") {
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

	wordCh := make(chan []string, 1000)
	wordMap := make(map[string][]string, 10000)
	go func() {
		for {
			select {
			case w := <-wordCh:
				wordMap[w[0]] = w
			}
		}
	}()

	_ = filepath.Walk(BaseDir, func(path string, info os.FileInfo, err error) error {

		if info.IsDir() || !strings.HasSuffix(path, ".rime") {
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

			br := bufio.NewReader(src)
			for {
				s, err := br.ReadString('\n')
				if err != nil {
					break
				}
				ss := strings.Split(s, "\t")
				if len(ss) != 3 {
					logrus.Errorf("dict format error: [%s]", s)
					break
				}
				wordCh <- []string{ss[0], ss[1], ss[2]}
			}

			logrus.Infof("[Deduplication] file %s processed", path)

		})

		return nil
	})

	wg.Wait()

	outFile, err := os.OpenFile(OutFile, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() { _ = outFile.Close() }()

	var keys []string
	for k := range wordMap {
		keys = append(keys, k)
	}
	sort.Sort(pinyin.ByPinyin(keys))

	for _, k := range keys {
		s := strings.Join(wordMap[k], "\t")
		_, err := fmt.Fprintln(outFile, strings.TrimSpace(s))
		if err != nil {
			logrus.Error(err)
		}
	}

}
