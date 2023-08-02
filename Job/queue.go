package Job

import (
	"dns-check/config"
	"github.com/emirpasic/gods/lists/arraylist"
	"time"
)

var MainJob = arraylist.New()
var totalCount int64

var AllJob = make(chan *Job, config.ProcessCount*5)
var DoneJob = make(chan struct{}, config.ProcessCount*5)

func init() {
	go GetJob()
}

func GetJob() {
	for {
		<-DoneJob
		j, _ := MainJob.Get(0)
		if j != nil {
			AllJob <- j.(*Job)
		}
		MainJob.Remove(0)
		totalCount--
	}
	// every second get 100 jobs from MainJob add to AllJob
}

func GetCount() int64 {
	return totalCount
}

func AddJob(jobs []*Job) {
	if len(jobs) == 0 {
		return
	}
	for i := range jobs {
		MainJob.Add(jobs[i])
		totalCount++
	}
	MainJob.Sort(jobComparator)
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
