package model

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	JobModel   string `json:"job_model"` // DNS Whois 混合
	DomainNumb int    `json:"domain_numb"`
	FinishNumb int    `json:"finish_numb"`
	Domains    string `json:"domains"`
	Status     int    `json:"status" gorm:"default:1"` //1:未开始 2:进行中 3:暂停 4:已完成
}
