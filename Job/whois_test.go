package Job

import (
	"fmt"
	"github.com/likexian/whois"
	"golang.org/x/net/proxy"
	"testing"
	"time"
)

func TestKuaiDaili(t *testing.T) {
	//u, e := url.Parse("http://brd-customer-hl_3cf009f7-zone-data_center:9yvrj6jf2bqk@brd.superproxy.io:22225")
	//if e != nil {
	//	panic(e)
	//}
	//dialer, e := proxy.FromURL(u, proxy.Direct)

	dialer, e := proxy.SOCKS5("tcp", "127.0.0.1:24000", nil, proxy.Direct)
	if e != nil {
		panic(e)
	}
	// 创建 WHOIS 客户端
	client := whois.NewClient()

	client = client.SetDialer(dialer)
	client = client.SetTimeout(60 * time.Second)

	result, rerr := client.Whois("e3128df4fjwe62f11.io", "whois.iana.org")
	if rerr != nil {
		panic(rerr)
	}
	fmt.Printf("%v\n", result)
}
