package config

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Config struct {
	DBUser     string
	DBPass     string
	DBHost     string
	DBPort     string
	DBName     string
	RedisAddr  string
	ServerPort string
}

func LoadConfig() *Config {
	return &Config{
		DBUser:     "root",
		DBPass:     "123456",
		DBHost:     "192.168.235.129",
		DBPort:     "3306",
		DBName:     "lbs",
		RedisAddr:  "192.168.235.129:6379",
		ServerPort: "6061",
	}
}

func InitDB() {
	cfg := LoadConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db

	// 自动迁移模型
	//db.AutoMigrate(&models.User{}, &models.Clothing{}, &models.DailySalary{}, &models.MonthlySalary{})
}
