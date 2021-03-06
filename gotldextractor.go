package gotldextractor

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	MOZILLA_EFFECTIVE_TLD_NAMES_PATH = "http://mxr.mozilla.org/mozilla-central/source/netwerk/dns/effective_tld_names.dat?raw=1"
)

// TLDResult is the structure that stores the Subdomain, Domain, TLD makeup of a Hostname

type TLDResult struct {
	Subdomain string
	Domain    string
	TLD       string

	// rules that were used to match this TLDResult
	Rules []string
}

func (tldresult *TLDResult) GetHostname() string {
	if len(tldresult.Domain) == 0 {
		return tldresult.TLD
	}
	if len(tldresult.Subdomain) == 0 {
		return tldresult.Domain + "." + tldresult.TLD
	}
	return tldresult.Subdomain + "." + tldresult.Domain + "." + tldresult.TLD
}

func (tldresult *TLDResult) GetDomainTLD() string {
	if len(tldresult.Domain) == 0 {
		return tldresult.TLD
	}
	return tldresult.Domain + "." + tldresult.TLD
}

type TLDExtractor struct {
	RootNode *TLDExtractorNode
}

type TLDExtractorNode struct {
	Character  string              // character of current node
	ChildNodes []*TLDExtractorNode // reference to child nodes

	/*
		is current node an End?
		Note: nodes can be an End and contain ChildNodes, consider the case when
		 .ac.co.uk
		 .co.uk
	*/
	IsEnd       bool
	HasAsterisk bool
	HasNot      bool

	Count int
	Depth int

	Rules []string // list of rules that end with current Node
}

func (tldextractor *TLDExtractor) Build() (bool, error) {
	tldextractor.RootNode = &TLDExtractorNode{}
	tldextractor.RootNode.Depth = 0
	scanner := bufio.NewScanner(bytes.NewReader(effective_tld_names_dat()))
	for scanner.Scan() {
		tldextractor.AddTLD(scanner.Text())
	}
	return true, nil
}

func (tldextractor *TLDExtractor) BuildFromDataFile(path string) (bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	tldextractor.RootNode = &TLDExtractorNode{}
	tldextractor.RootNode.Depth = 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tldextractor.AddTLD(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (tldextractor *TLDExtractor) BuildFromURL(url string) (bool, error) {
	use_url := strings.TrimSpace(url)
	if len(use_url) == 0 {
		use_url = MOZILLA_EFFECTIVE_TLD_NAMES_PATH
	}
	resp, err := http.Get(use_url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	tldextractor.RootNode = &TLDExtractorNode{}
	tldextractor.RootNode.Depth = 0
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		tldextractor.AddTLD(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (tldextractor *TLDExtractor) AddTLD(tld string) {
	use_tld := strings.TrimSpace(tld)
	if strings.HasPrefix(use_tld, "//") {
		return
	}
	if len(use_tld) == 0 {
		return
	}
	//fmt.Println("Adding: ", tld)
	// if tld does not begin with a ".", we prepend it
	if !strings.HasPrefix(use_tld, "!") && !strings.HasPrefix(use_tld, "*") && !strings.HasPrefix(use_tld, ".") {
		use_tld = "." + use_tld
	}
	current_node := tldextractor.RootNode
	for i := 0; i < len(use_tld); i++ {
		// tld[len(use_tld)-1-i] is the effective character
		found := false
		current_char := string(use_tld[len(use_tld)-1-i])
		if current_char == "*" {
			current_node.HasAsterisk = true
			current_node.IsEnd = true
			current_node.Rules = append(current_node.Rules, tld)
			continue
		} else if current_char == "!" {
			current_node.HasNot = true
			current_node.IsEnd = true
			current_node.Rules = append(current_node.Rules, tld)
			continue
		}
		for _, n := range current_node.ChildNodes {
			if n.Character == current_char {
				found = true
				current_node = n
			}
			if found {
				current_node.Count = current_node.Count + 1
				break
			}
		}
		if !found {
			//fmt.Println(strings.Repeat(" ", current_node.Depth+1), "Creating a new node for ", string(use_tld[len(use_tld)-1-i]))
			node := TLDExtractorNode{}
			node.Character = current_char
			node.Count = 1
			node.Depth = current_node.Depth + 1
			current_node.ChildNodes = append(current_node.ChildNodes, &node)
			current_node = &node
		}
		if i == len(use_tld)-1 {
			current_node.IsEnd = true
			current_node.Rules = append(current_node.Rules, tld)
		}
		//fmt.Println("use_tld: ", use_tld, " i:", i, " len(use_tld): ", len(use_tld), " current_node (char): ", current_node.Character, " Itoa: ", string(current_node.Character), " current_node (Depth): ", current_node.Depth, " current_node (Count): ", current_node.Count)
	}
}

func (tldextractor *TLDExtractor) PrettyPrint() {
	pretty_print_traverse_node(tldextractor.RootNode, []string{})
}

func pretty_print_traverse_node(node *TLDExtractorNode, prefix []string) {
	if node.IsEnd {
		//fmt.Println(strings.Repeat("    ", node.Depth), "Name: ", strings.Join(prefix, "")+node.Character, " Count: ", node.Count, " Depth: ", node.Depth, " End: ", node.IsEnd)
		fmt.Println(strings.Join(prefix, "")+node.Character, " Count: ", node.Count, " Depth: ", node.Depth, " End: ", node.IsEnd)
	}
	for _, n := range node.ChildNodes {
		var tmp_prefix = make([]string, len(prefix))
		copy(tmp_prefix, prefix)
		pretty_print_traverse_node(n, append(tmp_prefix, node.Character))
	}
}

func (tldextractor *TLDExtractor) ParseURL(url *url.URL) (TLDResult, error) {
	return tldextractor.ParseHost(url.Host)
}

func (tldextractor *TLDExtractor) ParseHost(host string) (TLDResult, error) {
	//fmt.Println(host)
	// we remove port information from the tld if it is present i.e. www.facebook.com:443
	var err error
	use_host := host
	if strings.Contains(host, ":") {
		use_host, _, err = net.SplitHostPort(host)
		if err != nil {
			return TLDResult{}, err
		}
	}
	current_node := tldextractor.RootNode
	lastIsEnd, lastDot := -1, -1
	hasAsterisk, hasNot := false, false
	var rules []string

	for i := 0; i < len(use_host); i++ {
		// use_host[len(use_host)-1-i] is the effective character
		found := false
		for _, n := range current_node.ChildNodes {
			if n.Character == string(use_host[len(use_host)-1-i]) {
				found = true
				current_node = n
				if current_node.IsEnd {
					lastIsEnd = len(use_host) - 1 - i
				}
				if current_node.Character == "." {
					lastDot = len(use_host) - 1 - i
				}
				if current_node.HasAsterisk {
					hasAsterisk = true
				}
				if current_node.HasNot {
					hasNot = true
				}
			}
			if found {
				if len(current_node.Rules) > 0 {
					rules = StringConcat(rules, current_node.Rules)
				}
				break
			}
		}
		if !found {
			break
		}
	}
	if lastIsEnd == -1 && !hasAsterisk && !hasNot {
		return TLDResult{"", "", use_host, []string{}}, nil
	}
	if hasAsterisk {
		// if hasAsterisk, we can set lastIsEnd to the next (lower index) dot after lastDot
		use_last_position := -1
		if lastIsEnd > lastDot {
			use_last_position = lastIsEnd
		} else {
			use_last_position = lastDot
		}
		hasAsterisk_index := strings.LastIndex(use_host[0:use_last_position], ".")
		//fmt.Println(host, "INSIDE hasAsterisk use_host[0:use_last_position]: ", use_host[0:use_last_position], " hasAsterisk_index:", hasAsterisk_index, " lastDot:", lastDot, " lastIsEnd:", lastIsEnd, "rules:", rules)
		if hasAsterisk_index == -1 {
			lastIsEnd = lastDot
		} else {
			lastIsEnd = hasAsterisk_index
		}
	}
	if hasNot {
		// if hasNot, we can set lastIsEnd to lastDot
		//fmt.Println("INSIDE hasNot lastDot:", lastDot, " lastIsEnd:", lastIsEnd)
		lastIsEnd = lastDot
	}
	tld := use_host[lastIsEnd+1:]
	subdomain_domain := use_host[0:lastIsEnd]
	domain_index := strings.LastIndex(subdomain_domain, ".")
	subdomain, domain := "", ""
	if domain_index == -1 {
		domain = subdomain_domain
	} else {
		domain = subdomain_domain[domain_index+1:]
		subdomain = subdomain_domain[:domain_index]
	}
	r := TLDResult{subdomain, domain, tld, rules}
	//fmt.Println("Subdomain: ", subdomain, " Domain: ", domain, " TLD: ", tld, " subdomain_domain: ", subdomain_domain)
	return r, nil
}

func StringConcat(old1, old2 []string) []string {
	// taken from https://groups.google.com/forum/#!topic/golang-nuts/mRUD0KffSG4
	newslice := make([]string, len(old1)+len(old2))
	copy(newslice, old1)
	copy(newslice[len(old1):], old2)
	return newslice
}
