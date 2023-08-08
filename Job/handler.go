package Job

import (
	"dns-check/database"
	"dns-check/logger"
	"dns-check/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

func HandlerJob() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				DoneJob <- struct{}{}
			case j := <-AllJob:
				do(j)
				DoneJob <- struct{}{}
			}
		}
	}()
}

func do(j *Job) {
	if j == nil {
		return
	}
	logger.Logger("job do", logger.INFO, nil, fmt.Sprintf("job %v domain %v", j.JobId, j.Domain))
	db := database.GetInstance()
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
	if j.Domain == "" {
		return
	}
	var dm model.Domain
	dm.Domain = j.Domain
	dm.JobId = j.JobId
	defer func() {
		if err := recover(); err != nil {
			logger.Logger("job", logger.ERROR, nil, err.(error).Error())
		}
		db.Create(&dm)
	}()
	if j.JobModel == "DNS" {
		j.DoNsLookUp(&dm)
	} else if j.JobModel == "Whois" {
		j.DoWhois(&dm)
	} else {
		j.DoNsLookUp(&dm)
		if dm.RCode == "999" {
			j.DoWhois(&dm)
		}
	}
}

func (j *Job) DoNsLookUp(dm *model.Domain) {
	ns, status, err := lookupNS(j)
	if err != nil {
		dm.Checked = "false"
		dm.RCode = "999"
		return
	}
	if ns != nil {
		nss := ""
		for i := range *ns {
			nss += (*ns)[i] + ","
		}
		dm.NameServers = nss
	}
	if status == "free" {
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
