package main

import (
	"fmt"
	"huiyuanguanli/config"
	"huiyuanguanli/controllers"
	"huiyuanguanli/models"
	"huiyuanguanli/routes"
	"huiyuanguanli/services"
	"huiyuanguanli/utils"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	if err := config.LoadConfig("config/config.yaml"); err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	cfg := config.GetConfig()

	// 初始化Redis连接
	redisAddr := fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port)
	if err := utils.InitRedis(redisAddr, cfg.Redis.Password, cfg.Redis.DB); err != nil {
		log.Fatalf("初始化Redis连接失败: %v", err)
	}

	// 初始化邮件客户端
	emailConfig := &utils.EmailConfig{
		Host:     cfg.Email.Host,
		Port:     cfg.Email.Port,
		Username: cfg.Email.Username,
		Password: cfg.Email.Password,
		From:     cfg.Email.From,
	}
	emailClient := utils.NewEmailClient(emailConfig)

	// 初始化短信客户端
	smsConfig := &utils.SMSConfig{
		SecretID:   cfg.TencentCloud.SecretID,
		SecretKey:  cfg.TencentCloud.SecretKey,
		AppID:      cfg.TencentCloud.AppID,
		SignName:   cfg.TencentCloud.SignName,
		TemplateID: cfg.TencentCloud.TemplateID,
	}
	smsClient, err := utils.NewSMSClient(smsConfig)
	if err != nil {
		log.Fatalf("初始化短信客户端失败: %v", err)
	}

	// 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/huiyuanguanli?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// 自动迁移数据库表
	err = db.AutoMigrate(
		&models.User{},
		&models.Member{},
		&models.Store{},
		&models.StoreStaff{},
		&models.StoreDevice{},
		&models.StoreDeviceCheck{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 创建 Gin 引擎
	r := gin.Default()

	// 初始化服务和控制器
	userService := services.NewUserService(db)
	userController := controllers.NewUserController(userService)

	authService := services.NewAuthService(db)
	authController := controllers.NewAuthController(authService)

	memberService := services.NewMemberService(db)
	memberController := controllers.NewMemberController(memberService)

	storeService := services.NewStoreService(db)
	storeController := controllers.NewStoreController(storeService)

	// 注册路由
	routes.SetupRoutes(r, userController, authController, memberController, storeController)

	// 启动服务器
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}

	log.Println("服务启动成功!")
}
