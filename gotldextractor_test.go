package gotldextractor

import (
	_ "fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	tldextract := TLDExtractor{}
	_, err := tldextract.Build()
	if err != nil {
		t.Error("Error building tree! - ", err)
	}
}

func TestFetchEffectiveTLDNames(t *testing.T) {
	tldextract := TLDExtractor{}
	_, err := tldextract.BuildFromURL("")
	if err != nil {
		t.Error("Error building tree! - ", err)
	}
}

func invalidParseResult(t *testing.T, r *TLDResult) {
	t.Error("invalid parse result for ", r.GetHostname(), " subdomain: ", r.Subdomain, " domain: ", r.Domain, " tld: ", r.TLD, " rules: ", r.Rules)
}

func TestParseHosts(t *testing.T) {
	tldextract := TLDExtractor{}
	_, err := tldextract.Build()
	if err != nil {
		t.Error("Error building tree! - ", err)
	}
	var r TLDResult
	r, err = tldextract.ParseHost("meh.thisshouldwork.ck")
	if r.Subdomain != "" && r.Domain != "meh" && r.TLD != "thisshouldwork.ck" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("meh.www.ck")
	if r.Subdomain != "meh" && r.Domain != "www" && r.TLD != "ck" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("www.guy.kawasaki.jp")
	if r.Subdomain != "" && r.Domain != "www" && r.TLD != "guy.kawasaki.jp" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("www.city.kawasaki.jp")
	if r.Subdomain != "www" && r.Domain != "city" && r.TLD != "kawasaki.jp" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("www.cnn.com")
	if r.Subdomain != "www" && r.Domain != "cnn" && r.TLD != "com" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("www.bbc.co.uk")
	if r.Subdomain != "www" && r.Domain != "bbc" && r.TLD != "co.uk" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("www.sina.com.cn")
	if r.Subdomain != "www" && r.Domain != "sina" && r.TLD != "com.cn" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("weibo.sina.com.cn:8080")
	if r.Subdomain != "weibo" && r.Domain != "sina" && r.TLD != "com.cn" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("news.ycombinator.com")
	if r.Subdomain != "news" && r.Domain != "ycombinator" && r.TLD != "com" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("www.github.com")
	if r.Subdomain != "www" && r.Domain != "github" && r.TLD != "com" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("www.github.com:443")
	if r.Subdomain != "www" && r.Domain != "github" && r.TLD != "com" {
		invalidParseResult(t, &r)
	}
	r, err = tldextract.ParseHost("www.facebook.com:8080")
	if r.Subdomain != "www" && r.Domain != "facebook" && r.TLD != "com" {
		invalidParseResult(t, &r)
	}
}
