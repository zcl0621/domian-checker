package Job

import (
	"fmt"
	"testing"
)

func TestJob(t *testing.T) {
	var jobs []*Job
	for i := 0; i < 3; i++ {
		jobs = append(jobs, &Job{
			Domain: fmt.Sprintf("test-%d", i),
		})
	}

	jobs = nil
	for i := 4; i < 7; i++ {
		jobs = append(jobs, &Job{
			Domain: fmt.Sprintf("test-%d", i),
		})
	}
	AddJob(jobs)
	fmt.Println(MainJob)
}

func TestAddJob(t *testing.T) {
	go func() {
		var jobs []*Job
		for i := 100000; i < 200000; i++ {
			jobs = append(jobs, &Job{
				Domain: fmt.Sprintf("test-%d", i),
			})
		}
		AddJob(jobs)
	}()

	go func() {
		var tjobs []*Job
		for i := 0; i < 100000; i++ {
			tjobs = append(tjobs, &Job{
				Domain: fmt.Sprintf("test-%d", i),
			})
		}
		AddJob(tjobs)
	}()
	go func() {
		var tjobs []*Job
		for i := 200000; i < 300000; i++ {
			tjobs = append(tjobs, &Job{
				Domain: fmt.Sprintf("test-%d", i),
			})
		}
		AddJob(tjobs)
	}()
	for {
		j := <-AllJob
		fmt.Println(j)
		DoneJob <- struct{}{}
	}

}
