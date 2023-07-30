package Job

import (
	"dns-check/model"
	"dns-check/whoisparser"
	"encoding/json"
	"fmt"
	"github.com/likexian/whois"
	"github.com/miekg/dns"
	"testing"
	"time"
)

func TestKuaiDaili(t *testing.T) {
	//u, e := url.Parse("http://brd-customer-hl_3cf009f7-zone-data_center:9yvrj6jf2bqk@brd.superproxy.io:22225")
	//if e != nil {
	//	panic(e)
	//}
	//dialer, e := proxy.FromURL(u, proxy.Direct)

	//dialer, e := proxy.SOCKS5("tcp", "127.0.0.1:24000", nil, proxy.Direct)
	//if e != nil {
	//	panic(e)
	//}
	// 创建 WHOIS 客户端
	client := whois.NewClient()

	//client = client.SetDialer(dialer)
	client = client.SetTimeout(60 * time.Second)

	result, rerr := client.Whois("w3schools.co", "whois.iana.org")
	if rerr != nil {
		panic(rerr)
	}
	fmt.Printf("%v\n", result)
	parseResult, e := whoisparser.Parse(result)
	if e != nil {
		panic(e)
	}
	d, _ := json.Marshal(parseResult)
	fmt.Printf("%s\n", d)
}

func TestDomain(t *testing.T) {
	var dm model.Domain
	tj := &Job{Domain: "google.co"}
	tj.DoNsLookUp(&dm)
	fmt.Println("dm.Checked", dm.Checked, "dm.NameServers", dm.NameServers)
	tj.DoWhois(&dm)
	fmt.Printf("%v", dm)
}

func TestDNS(t *testing.T) {
	var dnsClient = &dns.Client{
		Net:          "tcp",
		Timeout:      5 * time.Second,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	m := dns.Msg{}
	m.SetQuestion(dns.Fqdn("youtube.co"), dns.TypeNS)
	r, _, err := dnsClient.Exchange(&m, "8.8.8.8:53")
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if r.Rcode != dns.RcodeSuccess {
		t.Errorf(err.Error())
	}
	var nameServers []string
	for _, ans := range r.Answer {
		if ns, ok := ans.(*dns.NS); ok {
			nameServers = append(nameServers, ns.Ns)
		}
	}
	t.Log(nameServers)
}
