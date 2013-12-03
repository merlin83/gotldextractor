A TLD Extractor for Go

# Motivation

I am writing this as my first library in Go to familiarize myself with the language.

This is written due to a need to extract a domain name and TLD from url.URL.Host

The TLDExtractor reads from a datafile, taken from http://mxr.mozilla.org/mozilla-central/source/netwerk/dns/effective_tld_names.dat?raw=1 and generates a trie structure (prefix-tree) based on the reversed rule.

E.g. if the rule is "!city.kawasaki.jp", it is added to the trie as

    p -> j -> . -> i -> -> k -> a -> s -> a -> w -> a -> k -> . -> y -> t -> i -> c(!)

After the trie is generated, when a search is performed, the Parse function reverses the hostname, walks the trie structure and generates the TLDResult.

# Sample output taken from go test

    Parse result for  meh.thisshouldwork.ck  subdomain:    domain:  meh  tld:  thisshouldwork.ck  rules:  [*.ck]
    Parse result for  meh.www.ck  subdomain:  meh  domain:  www  tld:  ck  rules:  [*.ck !www.ck]
    Parse result for  www.guy.kawasaki.jp  subdomain:    domain:  www  tld:  guy.kawasaki.jp  rules:  [jp *.kawasaki.jp]
    Parse result for  www.city.kawasaki.jp  subdomain:  www  domain:  city  tld:  kawasaki.jp  rules:  [jp *.kawasaki.jp !city.kawasaki.jp]
    Parse result for  www.cnn.com  subdomain:  www  domain:  cnn  tld:  com  rules:  [com]
    Parse result for  www.bbc.co.uk  subdomain:  www  domain:  bbc  tld:  co.uk  rules:  [*.uk]
    Parse result for  www.sina.com.cn  subdomain:  www  domain:  sina  tld:  com.cn  rules:  [cn com.cn]
    Parse result for  weibo.sina.com.cn  subdomain:  weibo  domain:  sina  tld:  com.cn  rules:  [cn com.cn]
    Parse result for  news.ycombinator.com  subdomain:  news  domain:  ycombinator  tld:  com  rules:  [com]
    Parse result for  www.github.com  subdomain:  www  domain:  github  tld:  com  rules:  [com]
    Parse result for  www.facebook.com  subdomain:  www  domain:  facebook  tld:  com  rules:  [com]


# Examples

    r, err := tldextract.ParseHost("www.guy.kawasaki.jp")
    -> r.Subdomain = "", r.Domain = "www", r.TLD = "guy.kawasaki.jp" // based on the rule *.kawasaki.jp

    r, err := tldextract.ParseHost("www.city.kawasaki.jp")
    -> r.Subdomain = "www", r.Domain = "city", r.TLD = "kawasaki.jp" // based on the rules *.kawasaki.jp and !city.kawasaki.jp
