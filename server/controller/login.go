package controller

import (
	"dns-check/database"
	"dns-check/model"
	"dns-check/server/middleware/adminJwt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	Id    uint   `json:"id"`
	Token string `json:"token"`
}

func Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "参数错误"})
		return
	}
	db := database.GetInstance()
	var u model.User
	db.Model(&model.User{}).Where(&model.User{Username: request.Username, Password: request.Password}).First(&u)
	if u.ID == 0 {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "用户名或密码错误"})
		return
	}
	token, err := adminJwt.GenerateToken(u.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &errorResponse{ErrorCode: "生成token失败"})
		return
	}
	c.JSON(http.StatusOK, &loginResponse{Id: u.ID, Token: token})
}
