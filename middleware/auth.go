package middleware

import (
	"net/http"
	"strings"

	"huiyuanguanli/utils"

	"github.com/gin-gonic/gin"

	"huiyuanguanli/pkg/response"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct{}

// NewAuthMiddleware 创建认证中间件实例
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

// AuthRequired 认证中间件
func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取token
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "未提供认证信息"})
			ctx.Abort()
			return
		}

		// 检查token格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "认证信息格式错误"})
			ctx.Abort()
			return
		}

		// 验证token
		token := parts[1]
		claims, err := utils.ParseToken(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"code": 401, "message": "认证信息无效"})
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文中
		ctx.Set("userID", claims.UserID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}

// CheckPermission 权限检查中间件
func (m *AuthMiddleware) CheckPermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取用户ID
		_, exists := ctx.Get("user_id")
		if !exists {
			response.Unauthorized(ctx, "请先登录")
			ctx.Abort()
			return
		}

		// TODO: 从数据库获取用户权限并检查
		// 这里需要实现具体的权限检查逻辑

		ctx.Next()
	}
}

// CheckRole 角色检查中间件
func (m *AuthMiddleware) CheckRole(role string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取用户ID
		_, exists := ctx.Get("user_id")
		if !exists {
			response.Unauthorized(ctx, "请先登录")
			ctx.Abort()
			return
		}

		// TODO: 从数据库获取用户角色并检查
		// 这里需要实现具体的角色检查逻辑

		ctx.Next()
	}
}

// Cors 跨域中间件
func (m *AuthMiddleware) Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		ctx.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		ctx.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type")
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
