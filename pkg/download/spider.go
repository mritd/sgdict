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

	for i := 1; i < pageSize+1; i++ {
		pageAddr := fmt.Sprintf("%s/default/%d", addr, i)
		logrus.Infof("[QueryDictAddr] request addr: %s", pageAddr)
		resp, err := cli.R().Get(pageAddr)
		if err != nil {
			break
		}
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
		if err != nil {
			logrus.Warnf("request page %s error: %s", pageAddr, err)
			continue
		}
		doc.Find("#dict_detail_list > div").Each(func(i int, selection *goquery.Selection) {
			name := selection.Find("div.dict_detail_title_block > div > a").Text()
			href, ok := selection.Find("div.dict_detail_show > div.dict_dl_btn > a").Attr("href")
			if !ok {
				return
			}
			data[name] = href
		})
	}
	return data, nil
}

func downloadDict(baseDir string, data map[string]map[string]string) error {
	for d, addrs := range data {
		categoryDir := filepath.Join(baseDir, d)
		err := mkdir(categoryDir)
		if err != nil {
			return err
		}
		for n, l := range addrs {
			resp, err := client().R().Get(l)
			if err != nil {
				logrus.Errorf("failed to download dict %s: %s", n, err)
				continue
			}
			savePath := filepath.Join(categoryDir, n+".scel")
			err = ioutil.WriteFile(savePath, resp.Body(), 0644)
			if err != nil {
				logrus.Errorf("save dict [%s] failed: %s", savePath, err)
				continue
			}
		}
	}
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
