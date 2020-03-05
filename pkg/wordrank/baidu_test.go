package wordrank

import (
	"testing"
	"time"
)

func TestBaiduWordRank(t *testing.T) {
	rank, err := queryBaiduRank("mritd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rank)
}

func TestBaiduWorkRank(t *testing.T) {
	FilePath = "/Users/natural/Desktop/natural.rime"
	Timeout = 10 * time.Second
	RetryCount = 5
	RetryMaxWaitTime = 10 * time.Second
	BaiduWorkRank()
}
