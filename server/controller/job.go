package controller

import (
	"dns-check/Job"
	"dns-check/database"
	"dns-check/model"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type listRequest struct {
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"page_size" form:"page_size" binding:"required"`
}

type listResponse struct {
	Count int64         `json:"count"`
	Data  []oneResponse `json:"data"`
}
type oneResponse struct {
	Id         uint      `json:"id"`
	CreatedAt  time.Time `json:"created_at"`
	JobModel   string    `json:"job_model"` // dns whois mix
	DomainNumb int       `json:"domain_numb"`
	Status     int       `json:"status" gorm:"default:1"` //1:未开始 2:进行中 3:暂停 4:已完成
}

func ListJob(c *gin.Context) {
	var request listRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	sDB := db.Model(&model.Job{})
	var jobs []oneResponse
	var count int64
	sDB.Count(&count)
	sDB.Order("id desc").Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize).Scan(&jobs)
	c.JSON(http.StatusOK, &listResponse{Count: count, Data: jobs})
}

type addJobRequest struct {
	JobModel string   `json:"job_model"  binding:"required"`
	Domains  []string `json:"domains" binding:"required"`
}

type addJobResponse struct {
	Id       uint     `json:"id"`
	JobModel string   `json:"job_model"`
	Domains  []string `json:"domains"`
	Status   int      `json:"status"`
}

func AddJob(c *gin.Context) {
	if Job.GetCount() >= 100000 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务数量已达上限"})
		return
	}
	var request addJobRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	d, err := json.Marshal(request.Domains)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	job := model.Job{
		JobModel:   request.JobModel,
		DomainNumb: len(request.Domains),
		Domains:    fmt.Sprintf("%s", d),
		Status:     1,
	}
	querySet := db.Create(&job)
	if querySet.Error != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "新增失败"})
		return
	}
	c.JSON(http.StatusOK, &addJobResponse{Id: job.ID, JobModel: job.JobModel, Domains: request.Domains, Status: job.Status})
}

type jobInfoRequest struct {
	Id uint `json:"id"`
}

func StartJob(c *gin.Context) {
	var request jobInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	var job model.Job
	db.Model(&model.Job{}).Where("id = ?", request.Id).First(&job)
	if job.ID == 0 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务不存在"})
		return
	}
	if job.Status != 1 && job.Status != 3 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务状态错误"})
		return
	}
	var domains []string
	err := json.Unmarshal([]byte(job.Domains), &domains)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务数据错误"})
		return
	}
	go func(domains *[]string, jobId uint, jobModel string) {
		var jobs []*Job.Job
		for i := 0; i < len(*domains); i++ {
			jobs = append(jobs, &Job.Job{
				Domain:   (*domains)[i],
				JobId:    jobId,
				JobModel: jobModel,
			})
		}
		Job.AddJob(jobs)
	}(&domains, job.ID, job.JobModel)
	job.Status = 2
	db.Save(&job)
	c.JSON(http.StatusOK, &oneResponse{
		Id:         job.ID,
		CreatedAt:  job.CreatedAt,
		JobModel:   job.JobModel,
		DomainNumb: job.DomainNumb,
		Status:     job.Status,
	})

}

func PausedJob(c *gin.Context) {
	var request jobInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	var job model.Job
	db.Model(&model.Job{}).Where("id = ?", request.Id).First(&job)
	if job.ID == 0 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务不存在"})
		return
	}
	if job.Status != 2 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务状态错误"})
		return
	}
	job.Status = 3
	db.Save(&job)
	c.JSON(http.StatusOK, &oneResponse{
		Id:         job.ID,
		CreatedAt:  job.CreatedAt,
		JobModel:   job.JobModel,
		DomainNumb: job.DomainNumb,
		Status:     job.Status,
	})
}

func EndJob(c *gin.Context) {
	var request jobInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	var job model.Job
	db.Model(&model.Job{}).Where("id = ?", request.Id).First(&job)
	if job.ID == 0 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务不存在"})
		return
	}
	job.Status = 4
	db.Save(&job)
	c.JSON(http.StatusOK, &oneResponse{
		Id:         job.ID,
		CreatedAt:  job.CreatedAt,
		JobModel:   job.JobModel,
		DomainNumb: job.DomainNumb,
		Status:     job.Status,
	})
}

func DeleteJob(c *gin.Context) {
	var request jobInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	var job model.Job
	db.Model(&model.Job{}).Where("id = ?", request.Id).First(&job)
	if job.ID == 0 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务不存在"})
		return
	}
	db.Delete(&job)
	db.Model(&model.Domain{}).Where("job_id = ?", request.Id).Delete(&model.Domain{})
	c.JSON(http.StatusOK, "")
}

type processJobResponse struct {
	Process int `json:"process"`
}

func ProcessJob(c *gin.Context) {
	var request jobInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	var job model.Job
	db.Model(&model.Job{}).Where("id = ?", request.Id).First(&job)
	if job.ID == 0 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务不存在"})
		return
	}
	percent := int(float64(job.FinishNumb) / float64(job.DomainNumb) * 100)
	c.JSON(http.StatusOK, &processJobResponse{
		Process: percent,
	})
}

func ExportJob(c *gin.Context) {
	var request jobInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	var job model.Job
	db.Model(&model.Job{}).Where("id = ?", request.Id).First(&job)
	if job.ID == 0 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务不存在"})
		return
	}
	if job.Status != 4 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "任务未完成"})
		return
	}
	var domains []model.Domain
	db.Model(&model.Domain{}).Where("job_id = ?", request.Id).Find(&domains)
	csvData, err := convertToCsv(&domains, &job)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "导出失败"})
		return
	}

	// 将 CSV 数据返回给前端
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=job_%d.csv", job.ID))
	c.Data(http.StatusOK, "text/csv", csvData)
}

func convertToCsv(data *[]model.Domain, j *model.Job) ([]byte, error) {
	file, err := os.CreateTemp("", fmt.Sprintf("job_%d.csv", j.ID))
	if err != nil {
		return nil, err
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	writer := csv.NewWriter(file)

	e := exportMix(writer, data)
	if e != nil {
		return nil, e
	}
	writer.Flush()
	return ioutil.ReadFile(file.Name())
}

func exportMix(w *csv.Writer, data *[]model.Domain) error {
	if err := w.Write([]string{"Domain", "DNS", "Status", "Created Date", "Expires Date", "Domain Status"}); err != nil {
		return err
	}
	// 写入数据
	for _, domain := range *data {
		var ns string
		if domain.NameServers != "" && domain.WhoisNameServers != "" {
			ns = domain.NameServers + "|" + domain.WhoisNameServers
		} else if domain.WhoisNameServers != "" {
			ns = domain.WhoisNameServers
		} else if domain.NameServers != "" {
			ns = domain.NameServers
		}
		checked := ""
		if domain.Checked == "true" {
			checked = "taken"
		} else if domain.Checked == "false" {
			checked = "free"
		} else {
			checked = domain.Checked
		}
		if err := w.Write([]string{domain.Domain, ns, checked, domain.WhoisCreatedDate, domain.WhoisExpirationDate, domain.WhoisStatus}); err != nil {
			return err
		}
	}
	return nil
}
