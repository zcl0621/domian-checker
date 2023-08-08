package Job

import (
	"dns-check/config"
	"dns-check/logger"
	"github.com/emirpasic/gods/lists/arraylist"
	"time"
)

var MainJob = arraylist.New()
var totalCount int64

var AllJob = make(chan *Job, 65535)
var DoneJob = make(chan struct{}, 65535)
var AddJobChan = make(chan []*Job, config.ProcessCount*5)

func init() {
	go GetJob()
	go addJob()
}

func GetJob() {
	for {
		<-DoneJob
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
