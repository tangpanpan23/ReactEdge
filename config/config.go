package config

import (
	"os"
	"strconv"
)

// Config 应用配置
type Config struct {
	Server ServerConfig `json:"server"`
	AI     AIConfig     `json:"ai"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `json:"port"`
	ReadTimeout  int    `json:"read_timeout"`
	WriteTimeout int    `json:"write_timeout"`
}

// AIConfig AI配置
type AIConfig struct {
	MaxAnalysisTime int `json:"max_analysis_time"` // 分析最大时间(秒)
	CacheEnabled    bool `json:"cache_enabled"`    // 是否启用缓存
}

// Load 加载配置
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getEnvAsInt("SERVER_READ_TIMEOUT", 30),
			WriteTimeout: getEnvAsInt("SERVER_WRITE_TIMEOUT", 30),
		},
		AI: AIConfig{
			MaxAnalysisTime: getEnvAsInt("AI_MAX_ANALYSIS_TIME", 10),
			CacheEnabled:    getEnvAsBool("AI_CACHE_ENABLED", true),
		},
	}

	return config
}

// getEnv 获取环境变量
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量作为整数
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvAsBool 获取环境变量作为布尔值
func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
