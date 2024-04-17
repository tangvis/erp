package engine

import (
	jsonLib "encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		rawUserInfo := session.Get(UserInfoKey)
		if rawUserInfo == nil {
			// todo 错误返回
			c.JSON(http.StatusForbidden, gin.H{
				"message": "not authed",
			})
			c.Abort()
			return
		}
		var userInfo UserInfo
		_ = jsonLib.Unmarshal([]byte(rawUserInfo.(string)), &userInfo)
		c.Set(UserInfoKey, &userInfo)
	}
}
