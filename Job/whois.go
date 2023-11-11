package Job

import (
	"errors"
	"github.com/likexian/whois"
	"github.com/likexian/whois-parser"
	"strings"
)

func checkWhois(j *Job, useProxy bool) (*whoisparser.Domain, error) {
	defer func() {
		if err := recover(); err != nil {
			j.Err = err.(error).Error()
			//logger.Logger("checkWhois", logger.ERROR, nil, err.(error).Error())
		}
	}()
	if j.Domain == "" {
		return nil, nil
	}
	var outErr error
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
		result, err := client.Whois(j.Domain)
		if err != nil {
			outErr = err
			//logger.Logger("checkWhois", logger.ERROR, nil, fmt.Sprintf("domain %s result %s", j.Domain, err.Error()))
		}
		if result != "" {
			if strings.Contains(result, "Tonic whoisd") {
				var dms []string
				lines := strings.Split(result, "\n")
				for i := 1; i < len(lines); i++ {
					if strings.Contains(lines[i], ".") {
						d := strings.Split(lines[i], " ")
						dms = append(dms, d[len(d)-1])
					}
				}
				if len(dms) == 0 {
					return nil, errors.New("no domain")
				}
				var result whoisparser.Domain
				result.NameServers = dms
				result.Status = []string{"active"}
				return &result, nil
			} else {
				results := strings.Split(result, "source:       IANA")
				//logger.Logger("checkWhois", logger.INFO, nil, fmt.Sprintf("domain %s result %s", j.Domain, result))
				parseResult, e := whoisparser.Parse(results[len(results)-1])
				if e == nil {
					return parseResult.Domain, nil
				}
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
	return nil, outErr
}
