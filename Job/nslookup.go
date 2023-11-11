package Job

import (
	"fmt"
	"github.com/lixiangzhong/dnsutil"
	"strings"
)

//func dolookup(domain string, dnsServer string) (*[]string, int, error) {
//	m := dns.Msg{}
//	m.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
//	r, _, err := dnsClient.Exchange(&m, dnsServer)
//	if err != nil {
//		//logger.Logger("dolookup", logger.ERROR, nil, fmt.Sprintf("error: %s", err.Error()))
//		return nil, 0, err
//	}
//	if r.Rcode != dns.RcodeSuccess {
//		//logger.Logger("dolookup", logger.WARNING, nil, fmt.Sprintf("Rcode: %d", r.Rcode))
//		return nil, r.Rcode, nil
//	}
//	var nameServers []string
//	for _, ans := range r.Answer {
//		if ns, ok := ans.(*dns.NS); ok {
//			nameServers = append(nameServers, ns.Ns)
//		}
//	}
//	return &nameServers, 0, nil
//}

//func dolookup(domain string) (*[]string, string, error) {
//	var lookupNSError *net.DNSError
//	var nameServers []string
//	ns, err := net.LookupNS(domain)
//	if err != nil {
//		errors.As(err, &lookupNSError)
//	} else {
//		for i := range ns {
//			nameServers = append(nameServers, ns[i].Host)
//		}
//	}
//	if lookupNSError != nil {
//		if lookupNSError.IsNotFound {
//			return nil, "free", nil
//		} else {
//			return nil, "", err
//		}
//	}
//	return &nameServers, "taken", nil
//}

func dolookup(domain string) (*[]string, string, error) {
	var lastNS []string
	var dig dnsutil.Dig
	rsps, err := dig.Trace(domain)
	if err != nil {
		return nil, "", err
	}
	for _, rsp := range rsps {
		var thisNs []string
		for _, ns := range rsp.Msg.Ns {
			y := strings.Split(ns.String(), "\t")
			if len(y) > 1 {
				thisNs = append(thisNs, y[len(y)-1])
			}
		}
		if len(thisNs) > 0 {
			lastNS = thisNs
		}
	}
	if len(lastNS) == 0 {
		return nil, "free", nil
	}
	return &lastNS, "taken", nil
}

func lookupNS(j *Job) (*[]string, string, error) {
	if j.Domain == "" {
		return nil, "", fmt.Errorf("domain is empty")
	}
	for i := 0; i < 3; i++ {
		ns, status, err := dolookup(j.Domain)
		if err != nil {
			continue
		}
		return ns, status, nil
	}
	return nil, "", fmt.Errorf("cannot lookup domain")
}
