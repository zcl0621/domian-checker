package Job

import (
	"dns-check/config"
	"dns-check/database"
	"dns-check/logger"
	"dns-check/model"
	"github.com/emirpasic/gods/lists/arraylist"
	"gorm.io/gorm"
	"time"
)

var MainJob = arraylist.New()
var totalCount int64

var AllJob = make(chan *Job, 65535)
var DoneJob = make(chan uint, 65535)
var AddJobChan = make(chan []*Job, config.ProcessCount*5)

func init() {
	go GetJob()
	go addJob()
}

func GetJob() {
	for {
		jobId := <-DoneJob
		if jobId != 0 {
			finishJob(jobId)
		}
		func() {
			defer func() {
				if err := recover(); err != nil {
					logger.Logger("get job", logger.ERROR, nil, err.(error).Error())
				}
			}()
			j, _ := MainJob.Get(0)
			if j != nil {
				AllJob <- j.(*Job)
			}
			MainJob.Remove(0)
			totalCount--
		}()
	}
}

func finishJob(jobId uint) {
	db := database.GetInstance()
	db.Model(&model.Job{}).
		Where("id = ?", jobId).
		Updates(map[string]interface{}{
			"finish_numb": gorm.Expr("finish_numb + ?", 1),
			"status": gorm.Expr(`
			CASE
				WHEN finish_numb + 1 = domain_numb THEN 4
				ELSE status
			END
		`),
		})
}

func GetCount() int64 {
	return totalCount
}

func AddJob(jobs []*Job) {
	AddJobChan <- jobs
}

func addJob() {
	for {
		jobs := <-AddJobChan
		if len(jobs) == 0 {
			return
		}
		for i := range jobs {
			MainJob.Add(jobs[i])
			totalCount++
		}
		MainJob.Sort(jobComparator)
	}
}

func jobComparator(a, b interface{}) int {
	now := time.Now().UnixNano()
	if now%2 == 0 {
		return -1
	}
	if now%2 != 0 {
		return 1
	}
	return 0
}
