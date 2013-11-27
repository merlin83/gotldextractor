package gotldextractor

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	_ "strconv"
	"strings"
)

// TLDResult is the structure that stores the Subdomain, Domain, TLD makeup of a Hostname

type TLDResult struct {
	Subdomain string
	Domain    string
	TLD       string
}

func (tldresult *TLDResult) GetHostname() string {
	return tldresult.Domain + "." + tldresult.Domain + "." + tldresult.TLD
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
}

func (tldextractor *TLDExtractor) Build() (bool, error) {
	return tldextractor.BuildFromDataFile("dat/effective_tld_names.dat")
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
		trimmed_text := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(trimmed_text, "//") {
			continue
		}
		if len(trimmed_text) == 0 {
			continue
		}
		tldextractor.AddTLD(trimmed_text)
	}
	if err := scanner.Err(); err != nil {
		return false, err
	}
	return true, nil
}

func (tldextractor *TLDExtractor) AddTLD(tld string) {
	if len(tld) == 0 {
		return
	}
	//fmt.Println("Adding: ", tld)
	use_tld := tld
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
			continue
		} else if current_char == "!" {
			current_node.HasNot = true
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
	current_node := tldextractor.RootNode
	lastIsEnd, lastDot := -1, -1
	hasAsterisk, hasNot := false, false
	for i := 0; i < len(host); i++ {
		// host[len(host)-1-i] is the effective character
		found := false
		for _, n := range current_node.ChildNodes {
			if n.Character == string(host[len(host)-1-i]) {
				found = true
				current_node = n
			}
			if current_node.IsEnd && current_node.Character == "." {
				lastIsEnd = len(host) - 1 - i
			}
			if current_node.Character == "." {
				lastDot = len(host) - 1 - i
			}
			if current_node.HasAsterisk {
				hasAsterisk = true
			}
			if current_node.HasNot {
				hasNot = true
			}
			if found {
				break
			}
		}
		if !found {
			break
		}
	}
	if lastIsEnd == -1 && !hasAsterisk && !hasNot {
		return TLDResult{"", "", host}, nil
	}
	if hasAsterisk {
		// if hasAsterisk, we can set lastIsEnd to the next (lower index) dot after lastDot
		hasAsterisk_index := strings.LastIndex(host[0:lastDot], ".")
		//fmt.Println("INSIDE hasAsterisk hasAsterisk_index:", hasAsterisk_index, " lastDot:", lastDot, " lastIsEnd:", lastIsEnd)
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
	tld := strings.TrimLeft(host[lastIsEnd+1:], ".")
	subdomain_domain := strings.TrimRight(host[0:lastIsEnd+1], ".")
	domain_index := strings.LastIndex(subdomain_domain, ".")
	subdomain, domain := "", ""
	if domain_index == -1 {
		domain = subdomain_domain
	} else {
		domain = subdomain_domain[domain_index+1:]
		subdomain = subdomain_domain[:domain_index]
	}
	r := TLDResult{subdomain, domain, tld}
	//fmt.Println("Subdomain: ", subdomain, " Domain: ", domain, " TLD: ", tld, " subdomain_domain: ", subdomain_domain)
	return r, nil
}
