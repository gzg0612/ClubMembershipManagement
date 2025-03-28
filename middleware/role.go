package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从上下文中获取用户角色
		userRole, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未获取到用户角色"})
			ctx.Abort()
			return
		}

		// 检查用户角色是否匹配
		if userRole != role {
			ctx.JSON(http.StatusForbidden, gin.H{"code": 403, "message": "权限不足"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
