package download

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
)

const (
	Host             = "https://pinyin.sogou.com"
	MainCategoryPage = "https://pinyin.sogou.com/dict/"
	UA               = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.122 Safari/537.36"
)

func client() *resty.Client {
	return resty.New().
		SetLogger(logrus.StandardLogger()).
		SetTimeout(3*time.Second).
		SetRetryCount(2).
		SetRetryMaxWaitTime(3*time.Second).
		SetHeader("User-Agent", UA)
}

func queryMainCategory() (map[string]string, error) {
	data := make(map[string]string)

	resp, err := client().R().Get(MainCategoryPage)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}

	doc.Find("#dict_category_show > .dict_category_list > .dict_category_list_title > a").Each(func(i int, selection *goquery.Selection) {
		href, ok := selection.Attr("href")
		if !ok {
			return
		}
		hs := strings.Split(href, "?")
		if len(hs) == 0 {
			return
		}
		name := selection.Text()
		data[name] = Host + hs[0]
	})
	return data, nil
}

func queryDictAddr(addr string) (map[string]string, error) {
	data := make(map[string]string)

	cli := client()

	resp, err := cli.R().Get(addr)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, err
	}

	pageSize := 1
	doc.Find("#dict_page_list > ul > li > span > a").Each(func(i int, selection *goquery.Selection) {
		page := selection.Text()
		i, err := strconv.Atoi(page)
		if err != nil {
			return
		}
		if i > pageSize {
			pageSize = i
		}
	})

	logrus.Infof("[QueryDictAddr] %s page size: %d", addr, pageSize)

	var wg sync.WaitGroup
	wg.Add(pageSize)
	resCh := make(chan [2]string)

	for i := 1; i < pageSize+1; i++ {
		pageNum := i
		go func() {
			defer wg.Done()
			pageAddr := fmt.Sprintf("%s/default/%d", addr, pageNum)
			logrus.Debugf("[QueryDictAddr] request addr: %s", pageAddr)
			resp, err := cli.R().Get(pageAddr)
			if err != nil {
				logrus.Errorf("[QueryDictAddr] request page [%s] error: %s", pageAddr, err)
				return
			}
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
			if err != nil {
				logrus.Errorf("[QueryDictAddr] parse page [%s] error: %s", pageAddr, err)
				return
			}
			doc.Find("#dict_detail_list > div").Each(func(i int, selection *goquery.Selection) {
				name := selection.Find("div.dict_detail_title_block > div > a").Text()
				href, ok := selection.Find("div.dict_detail_show > div.dict_dl_btn > a").Attr("href")
				if !ok {
					return
				}
				resCh <- [2]string{name, href}
			})
		}()
	}

	go func() {
		for {
			select {
			case res, ok := <-resCh:
				if ok {
					data[res[0]] = res[1]
				} else {
					break
				}
			}
		}
	}()

	wg.Wait()
	close(resCh)
	return data, nil
}

func downloadDict(baseDir string, data map[string]map[string]string) error {

	var catWg sync.WaitGroup
	catWg.Add(len(data))

	for d, a := range data {
		addrs := a
		categoryDir := filepath.Join(baseDir, d)
		err := mkdir(strings.Replace(categoryDir, " ", "", -1))
		if err != nil {
			return err
		}
		go func() {
			defer catWg.Done()
			for n, a := range addrs {
				resp, err := client().R().Get(a)
				if err != nil {
					logrus.Errorf("download dict [%s] failed: %s", n, err)
					return
				}
				logrus.Infof("download dict [%s]", n)
				savePath := filepath.Join(categoryDir, n+".scel")
				err = ioutil.WriteFile(savePath, resp.Body(), 0644)
				if err != nil {
					logrus.Errorf("save dict [%s] failed: %s", savePath, err)
					return
				}
			}
		}()

	}
	catWg.Wait()
	return nil
}

func mkdir(dir string) error {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		if !info.IsDir() {
			return errors.New(fmt.Sprintf("[%s] already exist, but is not a dir", dir))
		}
	}
	return nil
}

func Run(baseDir string) {
	downMap := make(map[string]map[string]string)
	categories, err := queryMainCategory()
	if err != nil {
		logrus.Fatal(err)
	}

	for name, addr := range categories {
		data, err := queryDictAddr(addr)
		if err != nil {
			logrus.Error(err)
			continue
		}
		downMap[name] = data
	}

	err = downloadDict(baseDir, downMap)
	if err != nil {
		logrus.Fatal(err)
	}
}
