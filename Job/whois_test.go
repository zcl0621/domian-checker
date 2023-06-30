package Job

import (
	"fmt"
	"testing"
)

func TestGetWhois(t *testing.T) {
	j := &Job{
		Domain:     "baidu.com",
		NS:         nil,
		TryTime:    0,
		Err:        "",
		RecordType: "",
		JobId:      0,
		JobModel:   "",
	}
	x := checkWhois(j)
	fmt.Printf("%v\n", &j)
	fmt.Printf("%v\n", x)
}

func TestJob_DoNsLookUpt(t *testing.T) {
	ns, code, err := dolookup("facebook.com", "8.8.8.8")
	if err != nil {
		panic(err)
	}
	t.Log(ns)
	t.Log(code)
}
