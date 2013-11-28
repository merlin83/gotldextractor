A TLD Extractor for Go

# Motivation

I am writing this as my first library in Go to familiarize myself with the language.

This is written due to a need to extract a domain name and TLD from url.URL.Host

The TLDExtractor reads from a datafile, taken from http://mxr.mozilla.org/mozilla-central/source/netwerk/dns/effective_tld_names.dat?raw=1 and generates a trie structure (prefix-tree) based on the reversed rule.

E.g. if the rule is "!city.kawasaki.jp", it is added to the trie as

  p -> j -> . -> i -> -> k -> a -> s -> a -> w -> a -> k -> . -> y -> t -> i -> c(!)

After the trie is generated, when a search is performed, the Parse function reverses the hostname, walks the trie structure and generates the TLDResult.

# Example

   r, err := tldextract.ParseHost("www.guy.kawasaki.jp")

   -> r.Subdomain = "", r.Domain = "www", r.TLD = "guy.kawasaki.jp" // based on the rule *.kawasaki.jp

   r, err := tldextract.ParseHost("www.city.kawasaki.jp")

   -> r.Subdomain = "www", r.Domain = "city", r.TLD = "kawasaki.jp" // based on the rules *.kawasaki.jp and !city.kawasaki.jp
