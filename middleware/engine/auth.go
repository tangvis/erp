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
		user := userInfo(session.Get(common.UserInfoKey))
		if user == nil {
			json(c, nil, common.ErrAuth)
			c.Abort()
			return
		}
		c.Set(common.UserInfoKey, user)
	}
}

func userInfo(v any) *common.UserInfo {
	if v == nil {
		return nil
	}
	var user common.UserInfo
	_ = jsonLib.Unmarshal([]byte(v.(string)), &user)

	return &user
}
