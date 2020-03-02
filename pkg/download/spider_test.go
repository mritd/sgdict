package download

import "testing"

func TestQueryMainCategory(t *testing.T) {
	data, err := queryMainCategory()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}

func TestQueryDictAddr(t *testing.T) {
	data, err := queryDictAddr("https://pinyin.sogou.com/dict/cate/index/31")
	if err != nil {
		t.Fatal(err)
	}
	for k, v := range data {
		t.Logf("%s: %s", k, v)
	}
}

func TestDownloadDict(t *testing.T) {
	err := downloadDict("/tmp/spdict", map[string]map[string]string{
		"城市信息大全": {
			"中国高等院校（大学）大全【官方推荐】": "http://download.pinyin.sogou.com/dict/download_cell.php?id=20647&name=%E4%B8%AD%E5%9B%BD%E9%AB%98%E7%AD%89%E9%99%A2%E6%A0%A1%EF%BC%88%E5%A4%A7%E5%AD%A6%EF%BC%89%E5%A4%A7%E5%85%A8%E3%80%90%E5%AE%98%E6%96%B9%E6%8E%A8%E8%8D%90%E3%80%91",
			"政府机关团体机构大全【官方推荐】":   "http://download.pinyin.sogou.com/dict/download_cell.php?id=22421&name=%E6%94%BF%E5%BA%9C%E6%9C%BA%E5%85%B3%E5%9B%A2%E4%BD%93%E6%9C%BA%E6%9E%84%E5%A4%A7%E5%85%A8%E3%80%90%E5%AE%98%E6%96%B9%E6%8E%A8%E8%8D%90%E3%80%91",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
