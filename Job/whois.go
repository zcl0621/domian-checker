package Job

import (
	"github.com/likexian/whois-parser"
	"time"
)

func checkWhois(j *Job) *whoisparser.Domain {
	if j.Domain == "" {
		return nil
	}
	var count int
	for {
		// 查询 WHOIS 信息
		result, err := whoisClient.Whois(j.Domain)
		if err == nil {
			parseResult, e := whoisparser.Parse(result)
			if e == nil {
				return parseResult.Domain
			}
		}
		time.Sleep(time.Second * 3)
		count++
		if count >= 3 {
			break
		}
	}

	return nil
}
