package Job

import (
	"dns-check/database"
	"dns-check/logger"
	"dns-check/model"
	"fmt"
	"time"
)

func HandlerJob() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				DoneJob <- 0
			case j := <-AllJob:
				do(j)
				DoneJob <- j.JobId
				j = nil
			}
		}
	}()
}

func do(j *Job) {
	if j == nil {
		return
	}
	logger.Logger("job do", logger.INFO, nil, fmt.Sprintf("job %v domain %v", j.JobId, j.Domain))
	var dm model.Domain
	dm.Domain = j.Domain
	dm.JobId = j.JobId
	if j.Domain != "" {
		defer func() {
			if err := recover(); err != nil {
				logger.Logger("job", logger.ERROR, nil, err.(error).Error())
			}
		}()
		if j.JobModel == "DNS" {
			j.DoNsLookUp(&dm)
		} else if j.JobModel == "Whois" {
			j.DoWhois(&dm, true)
		} else if j.JobModel == "WhoisNoProxy" {
			j.DoWhois(&dm, false)
		}
		//else {
		//	j.DoNsLookUp(&dm)
		//	if dm.RCode == "999" {
		//		j.DoWhois(&dm)
		//	}
		//}
	}
	db := database.GetInstance()
	db.Create(&dm)
}

func (j *Job) DoNsLookUp(dm *model.Domain) {
	ns, status, err := lookupNS(j)
	if err != nil {
		dm.Checked = "error"
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
func (j *Job) DoWhois(dm *model.Domain, useProxy bool) {
	//logger.Logger("DoWhois", logger.INFO, nil, fmt.Sprintf("%v", j))
	whoisD, err := checkWhois(j, useProxy)
	if err != nil {
		dm.Checked = "error"
		dm.WhoisStatus = "no-domain"
		dm.WhoisNameServers = "no-nameServer"
		return
	}
	if whoisD == nil {
		dm.Checked = "false"
		dm.WhoisStatus = "no-domain"
		dm.WhoisNameServers = "no-nameServer"
		return
	}
	//logger.Logger("DoWhois checkWhois Status", logger.INFO, nil, fmt.Sprintf("job %v whoisD %v", j, whoisD))
	if whoisD.Status == nil {
		dm.Checked = "false"
		dm.WhoisStatus = "no-domain"
	} else if len(whoisD.Status) == 0 {
		dm.Checked = "false"
		dm.WhoisStatus = "no-domain"
	} else {
		dm.Checked = "true"
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
	if whoisD.CreatedDate == "" {
		dm.WhoisCreatedDate = "no-date"
	} else {
		if whoisD.CreatedDateInTime == nil {
			dm.WhoisCreatedDate = whoisD.CreatedDate
		} else if whoisD.CreatedDateInTime.IsZero() {
			dm.WhoisCreatedDate = whoisD.CreatedDate
		} else {
			dm.WhoisCreatedDate = whoisD.CreatedDateInTime.Format("2006-01-02 15:04:05")
		}

	}
	if whoisD.ExpirationDate == "" {
		dm.WhoisExpirationDate = "no-date"
	} else {
		if whoisD.ExpirationDateInTime == nil {
			dm.WhoisExpirationDate = whoisD.ExpirationDate
		} else if whoisD.ExpirationDateInTime.IsZero() {
			dm.WhoisExpirationDate = whoisD.ExpirationDate
		} else {
			dm.WhoisExpirationDate = whoisD.ExpirationDateInTime.Format("2006-01-02 15:04:05")
		}

	}
}
