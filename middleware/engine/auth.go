package engine

import (
	jsonLib "encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/tangvis/erp/common"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		rawUserInfo := session.Get(common.UserInfoKey)
		if rawUserInfo == nil {
			json(c, nil, common.ErrAuth)
			c.Abort()
			return
		}
		var userInfo common.UserInfo
		_ = jsonLib.Unmarshal([]byte(rawUserInfo.(string)), &userInfo)
		c.Set(common.UserInfoKey, &userInfo)
	}
}
