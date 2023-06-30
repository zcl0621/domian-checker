package Job

import (
	"github.com/miekg/dns"
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
	"8.8.8.8",
	"1.1.1.1",
	"8.26.56.26",
	"9.9.9.9",
	"208.67.222.222",
	"76.76.19.19",
	"176.103.130.130",
	"64.6.64.6",
	"185.225.168.168",
	"216.87.84.211",
	"77.88.8.8",
	"84.200.69.80",
	"209.244.0.3",
}

var dnsClient = &dns.Client{
	Net:          "tcp",
	Timeout:      5 * time.Second,
	DialTimeout:  5 * time.Second,
	ReadTimeout:  5 * time.Second,
	WriteTimeout: 5 * time.Second,
}