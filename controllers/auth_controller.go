package controllers

import (
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"huiyuanguanli/middleware"
	"huiyuanguanli/models"
	"huiyuanguanli/pkg/auth"
	"huiyuanguanli/pkg/response"
)

// AuthController 认证控制器
type AuthController struct {
	DB             *gorm.DB
	AuthMiddleware *middleware.AuthMiddleware
}

// NewAuthController 创建认证控制器
func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		DB:             db,
		AuthMiddleware: middleware.NewAuthMiddleware(),
	}
}

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	CaptchaCode string `json:"captcha_code"`
	CaptchaID   string `json:"captcha_id"`
}

// LoginResponse 登录响应数据
type LoginResponse struct {
	Token        string      `json:"token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	UserInfo     interface{} `json:"user_info"`
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ParamError(ctx, "请求参数错误")
		return
	}

	// TODO: 验证码校验
	if req.CaptchaCode != "" && req.CaptchaID != "" {
		// 验证验证码
		// 省略验证码校验逻辑
	}

	// 查询用户
	var user models.User
	if err := c.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			response.Unauthorized(ctx, "用户名或密码错误")
			return
		}
		response.ServerError(ctx, "系统错误")
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		response.Forbidden(ctx, "账号已被禁用或锁定")
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		// 记录登录失败
		c.DB.Model(&user).UpdateColumn("login_failed_count", gorm.Expr("login_failed_count + 1"))

		// 检查失败次数
		if user.LoginFailedCount+1 >= 5 {
			// 锁定账号
			c.DB.Model(&user).Updates(map[string]interface{}{
				"status": 2, // 锁定状态
			})
			response.Forbidden(ctx, "登录失败次数过多，账号已被锁定")
			return
		}

		response.Unauthorized(ctx, "用户名或密码错误")
		return
	}

	// 生成令牌
	accessToken, refreshToken, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.ServerError(ctx, "令牌生成失败")
		return
	}

	// 更新登录信息
	now := time.Now()
	loginIP := ctx.ClientIP()
	c.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_time":    now,
		"last_login_ip":      loginIP,
		"login_failed_count": 0,
	})

	// 记录登录日志
	loginLog := models.LoginLog{
		UserID:      &user.ID,
		Username:    user.Username,
		LoginType:   "PASSWORD",
		LoginStatus: "SUCCESS",
		IP:          loginIP,
		Device:      ctx.GetHeader("User-Agent"),
		LoginTime:   now,
	}
	c.DB.Create(&loginLog)

	// 组装用户信息
	var roles []string
	var permissions []string
	for _, role := range user.Roles {
		roles = append(roles, role.Code)
		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Code)
		}
	}

	userInfo := gin.H{
		"id":          user.ID,
		"username":    user.Username,
		"real_name":   user.RealName,
		"avatar":      user.Avatar,
		"roles":       roles,
		"permissions": permissions,
	}

	// 返回结果
	response.Success(ctx, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400, // 24小时
		UserInfo:     userInfo,
	})
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshToken 刷新访问令牌
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ParamError(ctx, "请求参数错误")
		return
	}

	// 刷新令牌
	newToken, err := auth.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(ctx, "无效的刷新令牌")
		return
	}

	// 解析刷新令牌
	claims, err := auth.ParseToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(ctx, "无效的刷新令牌")
		return
	}

	// 生成新的刷新令牌
	_, newRefreshToken, err := auth.GenerateToken(claims.UserID, claims.Username)
	if err != nil {
		response.ServerError(ctx, "令牌生成失败")
		return
	}

	// 返回结果
	response.Success(ctx, RefreshTokenResponse{
		Token:        newToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    86400, // 24小时
	})
}

// GetCurrentUser 获取当前用户信息
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "请登录后再操作")
		return
	}

	// 查询用户信息
	var user models.User
	if err := c.DB.Preload("Roles").Preload("Roles.Permissions").Where("id = ?", userID).First(&user).Error; err != nil {
		response.ServerError(ctx, "获取用户信息失败")
		return
	}

	// 组装权限信息
	var roles []map[string]interface{}
	var permissions []string
	for _, role := range user.Roles {
		roles = append(roles, map[string]interface{}{
			"id":          role.ID,
			"name":        role.Name,
			"code":        role.Code,
			"description": role.Description,
		})

		for _, perm := range role.Permissions {
			permissions = append(permissions, perm.Code)
		}
	}

	// 返回用户信息
	response.Success(ctx, gin.H{
		"id":              user.ID,
		"username":        user.Username,
		"real_name":       user.RealName,
		"avatar":          user.Avatar,
		"email":           user.Email,
		"phone":           user.Phone,
		"gender":          user.Gender,
		"status":          user.Status,
		"last_login_time": user.LastLoginTime,
		"last_login_ip":   user.LastLoginIP,
		"roles":           roles,
		"permissions":     permissions,
	})
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=20"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=NewPassword"`
}

// ChangePassword 修改密码
func (c *AuthController) ChangePassword(ctx *gin.Context) {
	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.ParamError(ctx, "请求参数错误")
		return
	}

	// 从上下文中获取用户ID
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Unauthorized(ctx, "请登录后再操作")
		return
	}

	// 查询用户信息
	var user models.User
	if err := c.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.ServerError(ctx, "获取用户信息失败")
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		response.Unauthorized(ctx, "原密码错误")
		return
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		response.ServerError(ctx, "密码加密失败")
		return
	}

	// 更新密码
	if err := c.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		response.ServerError(ctx, "密码修改失败")
		return
	}

	// 返回结果
	response.SuccessWithMessage(ctx, "密码修改成功", nil)
}

// Logout 用户退出登录
func (c *AuthController) Logout(ctx *gin.Context) {
	// JWT无状态，客户端只需清除本地存储的令牌即可
	// 这里只做日志记录
	userID, exists := ctx.Get("userID")
	if exists {
		// 记录登出日志
		loginLog := models.LoginLog{
			UserID:      userID.(*uint),
			Username:    ctx.GetString("username"),
			LoginType:   "LOGOUT",
			LoginStatus: "SUCCESS",
			IP:          ctx.ClientIP(),
			Device:      ctx.GetHeader("User-Agent"),
			LoginTime:   time.Now(),
		}
		c.DB.Create(&loginLog)
	}

	response.SuccessWithMessage(ctx, "已成功退出登录", nil)
}
