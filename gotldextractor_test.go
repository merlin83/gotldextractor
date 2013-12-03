package gotldextractor

import (
	"fmt"
	"testing"
)

func TestBuild(t *testing.T) {
	tldextract := TLDExtractor{}
	_, err := tldextract.Build()
	if err != nil {
		t.Error("Error building tree! - ", err)
	}
}

func TestBuildFromDataFile(t *testing.T) {
	tldextract := TLDExtractor{}
	_, err := tldextract.BuildFromDataFile("dat/effective_tld_names.dat")
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

func printParseResult(r *TLDResult) {
	fmt.Println("Parse result for ", r.GetHostname(), " subdomain: ", r.Subdomain, " domain: ", r.Domain, " tld: ", r.TLD, " rules: ", r.Rules)
}

func test_hostname(t *testing.T, tldextract TLDExtractor, hostname, check_subdomain, check_domain, check_tld string) {
	r, err := tldextract.ParseHost(hostname)
	if err != nil {
		t.Error("test_hostname ParseHost (", hostname, ") error:", err)
	}
	if r.Subdomain != check_subdomain && r.Domain != check_domain && r.TLD != check_tld {
		invalidParseResult(t, &r)
	}
	printParseResult(&r)
}

func TestParseHosts(t *testing.T) {
	tldextract := TLDExtractor{}
	_, err := tldextract.Build()
	if err != nil {
		t.Error("Error building tree! - ", err)
	}
	test_hostname(t, tldextract, "meh.thisshouldwork.ck", "", "meh", "thisshouldwork.ck")
	test_hostname(t, tldextract, "meh.www.ck", "meh", "www", "ck")
	test_hostname(t, tldextract, "www.meh.ck", "", "www", "meh.ck")
	test_hostname(t, tldextract, "www.guy.kawasaki.jp", "", "www", "guy.kawasaki.jp")
	test_hostname(t, tldextract, "www.city.kawasaki.jp", "www", "city", "kawasaki.jp")
	test_hostname(t, tldextract, "www.cnn.com", "www", "cnn", "com")
	test_hostname(t, tldextract, "www.bbc.co.uk", "www", "bbc", "co.uk")
	test_hostname(t, tldextract, "www.sina.com.cn", "www", "sina", "com.cn")
	test_hostname(t, tldextract, "weibo.sina.com.cn:8080", "weibo", "sina", "com.cn")
	test_hostname(t, tldextract, "news.ycombinator.com", "news", "ycombinator", "com")
	test_hostname(t, tldextract, "www.github.com", "www", "github", "com")
	test_hostname(t, tldextract, "www.github.com:443", "www", "github", "com")
	test_hostname(t, tldextract, "www.facebook.com:8080", "www", "facebook", "com")
}
