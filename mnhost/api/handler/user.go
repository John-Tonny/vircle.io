package handler

import (
	"context"
	"log"
	"net/http"

	userPb "vircle.io/mnhost/interface/out/user"

	"github.com/gin-gonic/gin"
)

type UserAPIHandler struct {
	userClient userPb.UserService
}

func GetUserHandler(userClient userPb.UserService) *UserAPIHandler {
	return &UserAPIHandler{
		userClient: userClient,
	}
}

func (s *UserAPIHandler) Login(c *gin.Context) {
	/*method := c.Request.Method

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
	c.Header("Access-Control-Allow-Credentials", "true")

	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}*/

	log.Printf("start login")
	user := userPb.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{"status": "err", "errmsg": err})
		//c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	resp, err := s.userClient.Auth(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "err", "errmsg": err})
		//c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"mobile":   user.Mobile,
		"password": user.Password,
		"token":    resp.Token,
		"ppp":      "bbb",
	})
}

func (s *UserAPIHandler) Sign(c *gin.Context) {
	user := userPb.User{}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "err", "errmsg": err})
		//c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, err := s.userClient.Create(context.Background(), &user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "err", "errmsg": err})
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "data": "bbb"})
}
