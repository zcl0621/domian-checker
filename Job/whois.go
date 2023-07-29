package Job

import (
	"dns-check/logger"
	"dns-check/whoisparser"
	"fmt"
	"github.com/likexian/whois"
)

func checkWhois(j *Job) *whoisparser.Domain {
	defer func() {
		if err := recover(); err != nil {
			j.Err = err.(error).Error()
			logger.Logger("checkWhois", logger.ERROR, nil, err.(error).Error())
		}
	}()
	if j.Domain == "" {
		return nil
	}
	var count int
	for {
		var client *whois.Client
		if count < 2 {
			client = newWhois()
		} else {
			client = newProxyWhois()
		}
		result, err := client.Whois(j.Domain)
		if err != nil {
			logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s result %s", j.Domain, err.Error()))
		}
		if result != "" {
			logger.Logger("checkWhois", logger.INFO, nil, fmt.Sprintf("domain %s result %s", j.Domain, result))
			parseResult, e := whoisparser.Parse(result)
			if e == nil {
				return parseResult.Domain
			} else {
				logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s result %s 格式化错误", j.Domain, result))
			}
		} else {
			logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s 未获取到返回值", j.Domain))
		}

		count++
		if count >= 3 {
			break
		}
	}
	return nil
}
