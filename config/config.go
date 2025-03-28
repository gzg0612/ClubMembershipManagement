package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config 配置结构体
type Config struct {
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redis"`

	TencentCloud struct {
		SecretID   string `yaml:"secret_id"`
		SecretKey  string `yaml:"secret_key"`
		AppID      string `yaml:"app_id"`
		SignName   string `yaml:"sign_name"`
		TemplateID string `yaml:"template_id"`
	} `yaml:"tencent_cloud"`

	Email struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		From     string `yaml:"from"`
	} `yaml:"email"`
}

var GlobalConfig Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &GlobalConfig)
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return &GlobalConfig
}
