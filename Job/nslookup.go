package Job

import (
	"fmt"
	"github.com/miekg/dns"
)

func dolookup(domain string, dnsServer string) (*[]string, int, error) {
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn(domain), dns.TypeNS)
	r, _, err := dnsClient.Exchange(&m, dnsServer)
	if err != nil {
		return nil, 0, err
	}
	if r.Rcode != dns.RcodeSuccess {
		return nil, r.Rcode, nil
	}
	var nameServers []string
	for _, ans := range r.Answer {
		if ns, ok := ans.(*dns.NS); ok {
			nameServers = append(nameServers, ns.Ns)
		}
	}
	return &nameServers, 0, nil
}

func lookupNS(j *Job) (*[]string, string, error) {
	if j.Domain == "" {
		return nil, "", fmt.Errorf("domain is empty")
	}
	for {
		ns, rCode, err := dolookup(j.Domain, DNS[j.TryTime])
		if err != nil {
			j.TryTime++
			if j.TryTime >= len(DNS) {
				return nil, "", err
			}
			continue
		}
		if rCode != dns.RcodeSuccess {
			return nil, dns.RcodeToString[rCode], nil
		}
		return ns, dns.RcodeToString[rCode], nil
	}
}
