package Job

import (
	"github.com/likexian/whois"
	"github.com/likexian/whois-parser"
	"strings"
)

func checkWhois(j *Job, useProxy bool) *whoisparser.Domain {
	defer func() {
		if err := recover(); err != nil {
			j.Err = err.(error).Error()
			//logger.Logger("checkWhois", logger.ERROR, nil, err.(error).Error())
		}
	}()
	if j.Domain == "" {
		return nil
	}
	var count int
	for {
		var client *whois.Client
		if !useProxy {
			client = newWhois()
		} else {
			if count < 2 {
				client = newWhois()
			} else {
				client = newProxyWhois()
			}
		}
		result, err := client.Whois(j.Domain, "whois.iana.org")
		if err != nil {
			//logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s result %s", j.Domain, err.Error()))
		}
		if result != "" {
			results := strings.Split(result, "source:       IANA")
			//logger.Logger("checkWhois", logger.INFO, nil, fmt.Sprintf("domain %s result %s", j.Domain, result))
			parseResult, e := whoisparser.Parse(results[len(results)-1])
			if e == nil {
				return parseResult.Domain
			}
			//else {
			//	logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s result %s 格式化错误", j.Domain, result))
			//}
		}
		//else {
		//	logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s 未获取到返回值", j.Domain))
		//}

		count++
		if count >= 3 {
			break
		}
	}
	return nil
}
