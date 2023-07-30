package Job

import (
	"github.com/likexian/whois"
	"github.com/miekg/dns"
	"golang.org/x/net/proxy"
	"time"
)

type Job struct {
	Domain     string
	NS         []string
	TryTime    int
	Err        string
	RecordType string
	JobId      uint
	JobModel   string //cdns whois mix
}

var DNS = []string{
	"8.8.8.8:53",
	"1.1.1.1:53",
	"8.26.56.26:53",
	"9.9.9.9:53",
	"208.67.222.222:53",
	"76.76.19.19:53",
	"176.103.130.130:53",
	"64.6.64.6:53",
	"185.225.168.168:53",
	"216.87.84.211:53",
	"77.88.8.8:53",
	"84.200.69.80:53",
	"209.244.0.3:53",
}

var dnsClient = &dns.Client{
	Net:          "tcp",
	Timeout:      5 * time.Second,
	DialTimeout:  5 * time.Second,
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 5 * time.Second,
}

func newWhois() *whois.Client {
	return whois.NewClient()
}

func newProxyWhois() *whois.Client {
	var whoisProxyClient = whois.NewClient()
	dialer, e := proxy.SOCKS5("tcp", "proxy-manager:24000", nil, proxy.Direct)
	if e != nil {
		panic(e)
	}
	whoisProxyClient = whoisProxyClient.SetDialer(dialer)
	whoisProxyClient = whoisProxyClient.SetTimeout(60 * time.Second)
	return whoisProxyClient
}
