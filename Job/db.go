package Job

import (
	"github.com/likexian/whois"
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
