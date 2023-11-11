package Job

import (
	"dns-check/model"
	"fmt"
	"github.com/lixiangzhong/dnsutil"
	"strings"
	"testing"
)

func TestDig(t *testing.T) {
	lastNS := ""
	var dig dnsutil.Dig
	rsps, err := dig.Trace("google.site")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, rsp := range rsps {
		thisNs := ""
		//if rsp.Msg.Authoritative {
		//	for _, answer := range rsp.Msg.Answer {
		//		fmt.Println(answer)
		//	}
		//}
		for _, ns := range rsp.Msg.Ns {
			y := strings.Split(ns.String(), "\t")
			if len(y) > 1 {
				thisNs += y[len(y)-1] + ","
			}
		}
		if thisNs != "" {
			lastNS = thisNs
		}
		//fmt.Println("\tReceived from", rsp.Server, rsp.ServerIP)
	}
	fmt.Println(lastNS)
}

func TestDomain(t *testing.T) {
	var dm model.Domain
	tj := &Job{Domain: "dropbox.co"}
	tj.DoNsLookUp(&dm)
	fmt.Println("dm.Checked", dm.Checked, "dm.NameServers", dm.NameServers)
	tj.DoWhois(&dm, false)
	fmt.Printf("%v", dm)
}

func TestWhois(t *testing.T) {
	tj := &Job{Domain: "facebook.to"}
	whoisD, _ := checkWhois(tj, false)
	if whoisD != nil {
		fmt.Printf("%v", whoisD)
	}
}
