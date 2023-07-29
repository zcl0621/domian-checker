package Job

import (
	"dns-check/logger"
	"dns-check/whoisparser"
	"fmt"
	"time"
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
	client := whoisClient
	for {
		result, err := client.Whois(j.Domain, "whois.iana.org")
		if err != nil {
			logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s result %s", j.Domain, err.Error()))
			time.Sleep(time.Second * 3)
			count++
			if count >= 2 {
				client = whoisProxyClient
			}
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
		if count >= 3 {
			break
		}
	}
	return nil
}
