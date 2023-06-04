package model

import "gorm.io/gorm"

type Domain struct {
	gorm.Model
	JobId               uint   `json:"job_id" gorm:"index"`
	Domain              string `json:"domain"`
	NameServers         string `json:"name_servers"`
	RCode               string `json:"r_code"`
	Checked             string `json:"checked"`
	WhoisStatus         string `json:"whois_status"`
	WhoisNameServers    string `json:"whois_name_servers"`
	WhoisCreatedDate    string `json:"whois_created_date"`
	WhoisExpirationDate string `json:"whois_expiration_date"`
}
