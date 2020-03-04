package wordrank

import "testing"

func TestBaiduWordRank(t *testing.T) {
	rank, err := queryBaiduRank("mritd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(rank)
}
