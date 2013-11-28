package gotldextractor

import (
	"fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	tldextract := TLDExtractor{}
	_, err := tldextract.Build()
	if err != nil {
		t.Error("Error building tree!")
	}
}

func TestParseHosts(t *testing.T) {
	tldextract := TLDExtractor{}
	_, err := tldextract.Build()
	if err != nil {
		t.Error("Error building tree!")
	}
	HOSTNAMES := []string{
		"meh.thisshouldwork.ck",
		"meh.www.ck",
		"www.guy.kawasaki.jp",
		"www.city.kawasaki.jp",
		"www.cnn.com",
		"www.bbc.co.uk",
		"www.sina.com.cn",
		"weibo.sina.com.cn",
		"news.ycombinator.com",
		"www.github.com",
		"www.github.com:443",
		"www.facebook.com:8080"}
	for _, hostname := range HOSTNAMES {
		r, err := tldextract.ParseHost(hostname)
		fmt.Println(r)
		if err != nil {
			t.Error(r)
		}
	}
}
