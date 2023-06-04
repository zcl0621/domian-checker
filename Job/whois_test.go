package Job

import (
	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"golang.org/x/net/proxy"
	"net"
	"testing"
	"time"
)

func TestGetWhois(t *testing.T) {
	dialer, err := proxy.SOCKS5("tcp", "proxy-server.scraperapi.com:8001", &proxy.Auth{
		User:     "scraperapi",
		Password: "5b1f5899a2d4ff1da016aa9bc6f69be0",
	}, &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	// 创建 WHOIS 客户端
	client := whois.NewClient()

	client = client.SetDialer(dialer)
	client = client.SetTimeout(60 * time.Second)

	// 查询 WHOIS 信息
	result, err := client.Whois("facebook.com")
	if err == nil {
		parseResult, e := whoisparser.Parse(result)
		if e == nil {
			t.Log(parseResult.Domain)
		}
	}
	panic(err)
}

func TestJob_DoNsLookUpt(t *testing.T) {
	ns, code, err := dolookup("facebook.com", "8.8.8.8")
	if err != nil {
		panic(err)
	}
	t.Log(ns)
	t.Log(code)
}
