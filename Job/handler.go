package Job

import (
	"dns-check/database"
	"dns-check/logger"
	"dns-check/model"
	"dns-check/redisUtils"
	"dns-check/utils"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func HandlerJob() {
	go func() {
		for {
			currentJobByte, err := redisUtils.Get("current_job")
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			jobLen, err := redisUtils.LLen(fmt.Sprintf("job_%s", currentJobByte))
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			if jobLen == 0 {
				redisUtils.Del("current_job")
				db := database.GetInstance()
				jobIdStr := fmt.Sprintf("%s", currentJobByte)
				jobId := utils.ConvertAStringToInt(jobIdStr)
				db.Model(&model.Job{}).Where("id = ?", jobId).Update("status", 4)
				redisUtils.Del(fmt.Sprintf("job_%s", currentJobByte))
				continue
			}
			d, err := redisUtils.LRPop(fmt.Sprintf("job_%s", currentJobByte))
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			var j Job
			err = json.Unmarshal(d, &j)
			if err != nil {
				continue
			}
			if j.Domain == "" {
				continue
			}
			func(j *Job) {
				defer func() {
					if err := recover(); err != nil {
						logger.Logger("job", logger.ERROR, nil, err.(error).Error())
						redisUtils.LPush("job", d)
					}
				}()
				db := database.GetInstance()
				var dm model.Domain
				db.Where(&model.Domain{Domain: j.Domain, JobId: j.JobId}).First(&dm)
				dm.Domain = j.Domain
				dm.JobId = j.JobId
				switch j.JobModel {
				case "DNS":
					j.DoNsLookUp(&dm)
					break
				case "Whois":
					j.DoWhois(&dm)
					break
				default:
					j.DoNsLookUp(&dm)
					j.DoWhois(&dm)
					break
				}
				db.Save(&dm)
				db.Model(&model.Job{}).Where("id = ?", j.JobId).Update("finish_numb", gorm.Expr("finish_numb + ?", 1))
			}(&j)
		}
	}()
}

func (j *Job) DoNsLookUp(dm *model.Domain) {
	ns, rcode, err := lookupNS(j)
	if err != nil {
		dm.Checked = "false"
	}
	if ns != nil {
		nss := ""
		for i := range *ns {
			nss += (*ns)[i] + ","
		}
		dm.NameServers = nss
	}
	dm.RCode = rcode
	if rcode == "NXDOMAIN" || rcode == "REFUSED" || rcode == "SERVFAIL" {
		dm.Checked = "false"
	} else {
		dm.Checked = "true"
	}
}
func (j *Job) DoWhois(dm *model.Domain) {
	logger.Logger("DoWhois", logger.INFO, nil, fmt.Sprintf("%v", j))
	whoisD := checkWhois(j)
	if whoisD == nil {
		dm.WhoisStatus = "no-domain"
		dm.WhoisNameServers = "no-nameServer"
		return
	}
	logger.Logger("DoWhois checkWhois Status", logger.INFO, nil, fmt.Sprintf("job %v whoisD %v", j, whoisD))
	if whoisD.Status == nil {
		dm.WhoisStatus = "no-domain"
	} else if len(whoisD.Status) == 0 {
		dm.WhoisStatus = "no-domain"
	} else {
		status := ""
		for i := range whoisD.Status {
			status += whoisD.Status[i] + ","
		}
		dm.WhoisStatus = status
	}
	logger.Logger("DoWhois checkWhois NameServers", logger.INFO, nil, fmt.Sprintf("job %v whoisD %v", j, whoisD))
	if whoisD.NameServers == nil {
		dm.WhoisNameServers = "no-nameServer"
	} else if len(whoisD.NameServers) == 0 {
		dm.WhoisNameServers = "no-nameServer"
	} else {
		nameServer := ""
		for i := range whoisD.NameServers {
			nameServer += whoisD.NameServers[i]
		}
	}
	logger.Logger("DoWhois checkWhois Date", logger.INFO, nil, fmt.Sprintf("job %v whoisD %v", j, whoisD))
	dm.WhoisCreatedDate = whoisD.CreatedDateInTime.Format("2006-01-02 15:04:05")
	dm.WhoisExpirationDate = whoisD.ExpirationDateInTime.Format("2006-01-02 15:04:05")
}
