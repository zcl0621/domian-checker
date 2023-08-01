package Job

import (
	"dns-check/database"
	"dns-check/logger"
	"dns-check/model"
	"fmt"
	"gorm.io/gorm"
)

func HandlerJob() {
	go func() {
		for {
			j := <-AllJob
			do(j)
		}
	}()
}

func do(j *Job) {
	if j == nil {
		return
	}
	if j.Domain == "" {
		return
	}
	logger.Logger("job do", logger.INFO, nil, fmt.Sprintf("job %v domain %v", j.JobId, j.Domain))
	func(j *Job) {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger("job", logger.ERROR, nil, err.(error).Error())
			}
		}()
		db := database.GetInstance()
		var jm model.Job
		db.Where("id = ?", j.JobId).First(&jm)
		if jm.Status != 2 {
			return
		}
		var dm model.Domain
		db.Where(&model.Domain{Domain: j.Domain, JobId: j.JobId}).First(&dm)
		dm.Domain = j.Domain
		dm.JobId = j.JobId
		if j.JobModel == "DNS" {
			j.DoNsLookUp(&dm)
		} else if j.JobModel == "Whois" {
			j.DoWhois(&dm)
		} else {
			j.DoNsLookUp(&dm)
			if dm.Checked == "false" || dm.NameServers == "" {
				j.DoWhois(&dm)
			}
		}
		//logger.Logger("job switch", logger.INFO, nil, fmt.Sprintf("job %v domain %v", j.JobId, j.Domain))
		db.Save(&dm)
		//logger.Logger("job save", logger.INFO, nil, fmt.Sprintf("job %v domain %v", j.JobId, j.Domain))
		db.Model(&model.Job{}).
			Where("id = ?", j.JobId).
			Updates(map[string]interface{}{
				"finish_numb": gorm.Expr("finish_numb + ?", 1),
				"status": gorm.Expr(`
			CASE
				WHEN finish_numb + 1 = domain_numb THEN 4
				ELSE status
			END
		`),
			})
	}(j)
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
	//logger.Logger("DoWhois", logger.INFO, nil, fmt.Sprintf("%v", j))
	whoisD := checkWhois(j)
	if whoisD == nil {
		dm.WhoisStatus = "no-domain"
		dm.WhoisNameServers = "no-nameServer"
		return
	}
	//logger.Logger("DoWhois checkWhois Status", logger.INFO, nil, fmt.Sprintf("job %v whoisD %v", j, whoisD))
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
	if whoisD.NameServers == nil {
		dm.WhoisNameServers = "no-nameServer"
	} else if len(whoisD.NameServers) == 0 {
		dm.WhoisNameServers = "no-nameServer"
	} else {
		nameServer := ""
		for i := range whoisD.NameServers {
			nameServer += whoisD.NameServers[i] + ","
		}
		dm.WhoisNameServers = nameServer
	}
	if whoisD.CreatedDateInTime == nil {
		dm.WhoisCreatedDate = "no-date"
	} else if whoisD.CreatedDateInTime.IsZero() {
		dm.WhoisCreatedDate = "no-date"
	} else {
		dm.WhoisCreatedDate = whoisD.CreatedDateInTime.Format("2006-01-02 15:04:05")
	}
	if whoisD.ExpirationDateInTime == nil {
		dm.WhoisExpirationDate = "no-date"
	} else if whoisD.ExpirationDateInTime.IsZero() {
		dm.WhoisExpirationDate = "no-date"
	} else {
		dm.WhoisExpirationDate = whoisD.ExpirationDateInTime.Format("2006-01-02 15:04:05")
	}
}
