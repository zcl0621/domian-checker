package Job

import (
	"dns-check/logger"
	"fmt"
	"github.com/likexian/whois-parser"
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
	for {
		brow := OpenBrowser()
		page := OpenPage(brow, "https://whois.iana.org")
		SendSearch(page, j.Domain)
		result := GetResult(page)
		if result != "" {
			logger.Logger("checkWhois", logger.INFO, nil, fmt.Sprintf("domain %s result %s", j.Domain, result))
			page.Close()
			parseResult, e := whoisparser.Parse(result)
			if e == nil {
				return parseResult.Domain
			} else {
				logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s result %s 格式化错误", j.Domain, result))
			}
		} else {
			logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s 未获取到返回值", j.Domain))
		}
		page.Close()
		time.Sleep(time.Second * 3)
		count++
		if count >= 3 {
			break
		}
	}
	return nil
}
