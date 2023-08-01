package Job

import (
	"sync"
	"time"
)

var MainJob *Job
var MainLock sync.Mutex
var totalCount int64

var AllJob = make(chan *Job, 100)

func init() {
	go GetJob()
}

func GetJob() {
	// every second get 100 jobs from MainJob add to AllJob
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		MainLock.Lock()
		if MainJob == nil {
			MainLock.Unlock()
			continue
		}
		for i := 0; i < 100; i++ {
			if MainJob == nil {
				break
			}
			AllJob <- MainJob
			MainJob = MainJob.NextJob
			totalCount--
		}
		MainLock.Unlock()
	}
}

func GetCount() int64 {
	return totalCount
}

func AddJob(jobs []*Job) {
	MainLock.Lock()
	defer MainLock.Unlock()
	if len(jobs) == 0 {
		return
	}
	if MainJob == nil {
		totalCount = 1
		MainJob = jobs[0]
		jobs = jobs[1:]
	}
	currentJob := MainJob
	for _, job := range jobs {
		if job == nil {
			continue
		}
		totalCount++
		if currentJob.NextJob == nil {
			currentJob.NextJob = job
			currentJob = job
		} else {
			nextJob := currentJob.NextJob
			currentJob.NextJob = job
			job.NextJob = nextJob
			currentJob = job
		}
	}
}
