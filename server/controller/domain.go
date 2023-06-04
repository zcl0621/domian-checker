package controller

import (
	"dns-check/database"
	"dns-check/model"
	"github.com/gin-gonic/gin"
)

type domainListRequest struct {
	JobId    int `json:"job_id" form:"job_id" binding:"required"`
	Page     int `json:"page" form:"page" binding:"required"`
	PageSize int `json:"page_size" form:"page_size" binding:"required"`
}

type domainListResponse struct {
	Count int64          `json:"count"`
	Data  []model.Domain `json:"data"`
}

func ListDomain(c *gin.Context) {
	var request domainListRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	sDB := db.Model(&model.Domain{})
	var domains []model.Domain
	var count int64
	sDB.Where("job_id = ?", request.JobId).Count(&count)
	sDB.Order("id desc").Offset((request.Page - 1) * request.PageSize).Limit(request.PageSize).Scan(&domains)
	c.JSON(200, &domainListResponse{Count: count, Data: domains})
}
