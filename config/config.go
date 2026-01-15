package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config 应用配置
type Config struct {
	Server      ServerConfig      `yaml:"server" json:"server"`
	AI          AIConfig          `yaml:"ai" json:"ai"`
	Logging     LoggingConfig     `yaml:"logging" json:"logging"`
	Monitoring  MonitoringConfig  `yaml:"monitoring" json:"monitoring"`
	Development DevelopmentConfig `yaml:"development" json:"development"`
	Production  ProductionConfig  `yaml:"production" json:"production"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string `yaml:"port" json:"port"`
	ReadTimeout  int    `yaml:"read_timeout" json:"read_timeout"`
	WriteTimeout int    `yaml:"write_timeout" json:"write_timeout"`
	Host         string `yaml:"host" json:"host"`
	TLSEnabled   bool   `yaml:"tls_enabled" json:"tls_enabled"`
	TLSCertFile  string `yaml:"tls_cert_file" json:"tls_cert_file"`
	TLSKeyFile   string `yaml:"tls_key_file" json:"tls_key_file"`
}

// AIConfig AI配置
type AIConfig struct {
	Mode            string               `yaml:"mode" json:"mode"`
	MaxAnalysisTime int                  `yaml:"max_analysis_time" json:"max_analysis_time"`
	CacheEnabled    bool                 `yaml:"cache_enabled" json:"cache_enabled"`
	Cache           CacheConfig          `yaml:"cache" json:"cache"`
	CircuitBreaker  CircuitBreakerConfig `yaml:"circuit_breaker" json:"circuit_breaker"`
	Concurrency     ConcurrencyConfig    `yaml:"concurrency" json:"concurrency"`
	RateLimit       RateLimitConfig      `yaml:"rate_limit" json:"rate_limit"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	TTL        int `yaml:"ttl" json:"ttl"`
	MaxEntries int `yaml:"max_entries" json:"max_entries"`
}

// CircuitBreakerConfig 熔断器配置
type CircuitBreakerConfig struct {
	MaxFailures int `yaml:"max_failures" json:"max_failures"`
	Timeout     int `yaml:"timeout" json:"timeout"`
}

// ConcurrencyConfig 并发控制配置
type ConcurrencyConfig struct {
	MaxConcurrent int `yaml:"max_concurrent" json:"max_concurrent"`
}

// RateLimitConfig 请求限流配置
type RateLimitConfig struct {
	RequestsPerHour int `yaml:"requests_per_hour" json:"requests_per_hour"`
	BurstLimit      int `yaml:"burst_limit" json:"burst_limit"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level       string            `yaml:"level" json:"level"`
	Format      string            `yaml:"format" json:"format"`
	FileEnabled bool              `yaml:"file_enabled" json:"file_enabled"`
	FilePath    string            `yaml:"file_path" json:"file_path"`
	Rotation    LogRotationConfig `yaml:"rotation" json:"rotation"`
}

// LogRotationConfig 日志轮转配置
type LogRotationConfig struct {
	MaxSize    int `yaml:"max_size" json:"max_size"`
	MaxBackups int `yaml:"max_backups" json:"max_backups"`
	MaxAge     int `yaml:"max_age" json:"max_age"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	Enabled         bool          `yaml:"enabled" json:"enabled"`
	Port            string        `yaml:"port" json:"port"`
	HealthCheckPath string        `yaml:"health_check_path" json:"health_check_path"`
	Metrics         MetricsConfig `yaml:"metrics" json:"metrics"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	PrometheusEnabled bool   `yaml:"prometheus_enabled" json:"prometheus_enabled"`
	Path              string `yaml:"path" json:"path"`
}

// DevelopmentConfig 开发环境配置
type DevelopmentConfig struct {
	Debug          bool     `yaml:"debug" json:"debug"`
	CORSEnabled    bool     `yaml:"cors_enabled" json:"cors_enabled"`
	CORSOrigins    []string `yaml:"cors_origins" json:"cors_origins"`
	RequestLogging bool     `yaml:"request_logging" json:"request_logging"`
	SQLLogging     bool     `yaml:"sql_logging" json:"sql_logging"`
}

// ProductionConfig 生产环境配置
type ProductionConfig struct {
	GzipEnabled            bool `yaml:"gzip_enabled" json:"gzip_enabled"`
	RateLimitingEnabled    bool `yaml:"rate_limiting_enabled" json:"rate_limiting_enabled"`
	StaticCacheTTL         int  `yaml:"static_cache_ttl" json:"static_cache_ttl"`
	SecurityHeadersEnabled bool `yaml:"security_headers_enabled" json:"security_headers_enabled"`
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

// Load 加载配置
func Load() (*Config, error) {
	// 首先尝试加载主配置文件
	config, err := LoadFromFile("config/app.yaml")
	if err != nil {
		// 如果主配置文件不存在或加载失败，尝试加载示例配置文件
		fmt.Printf("⚠️ 主配置文件加载失败，尝试加载示例配置: %v\n", err)
		config, err = LoadFromFile("config/app.yaml.example")
		if err != nil {
			fmt.Printf("⚠️ 示例配置文件也不存在，使用默认配置: %v\n", err)
			return GetDefaultConfig(), nil
		}
		fmt.Println("✅ 从示例配置文件加载成功，请复制并修改为实际配置")
	}
	return config, nil
}

// LoadFromFile 从指定文件加载配置
func LoadFromFile(configPath string) (*Config, error) {
	// 默认配置
	config := GetDefaultConfig()

	// 尝试从文件加载
	if file, err := os.Open(configPath); err == nil {
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, fmt.Errorf("解析配置文件失败: %w", err)
		}
		fmt.Printf("✅ 从文件加载配置: %s\n", configPath)
	} else {
		fmt.Printf("⚠️ 配置文件不存在，使用默认配置: %s\n", configPath)
	}

	// 环境变量覆盖
	overrideFromEnv(config)

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	return config, nil
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         "6000",
			ReadTimeout:  30,
			WriteTimeout: 30,
			Host:         "0.0.0.0",
			TLSEnabled:   false,
		},
		AI: AIConfig{
			Mode:            "internal",
			MaxAnalysisTime: 60,
			CacheEnabled:    true,
			Cache: CacheConfig{
				TTL:        3600,
				MaxEntries: 1000,
			},
			CircuitBreaker: CircuitBreakerConfig{
				MaxFailures: 5,
				Timeout:     60,
			},
			Concurrency: ConcurrencyConfig{
				MaxConcurrent: 10,
			},
			RateLimit: RateLimitConfig{
				RequestsPerHour: 1000,
				BurstLimit:      100,
			},
		},
		Logging: LoggingConfig{
			Level:       "info",
			Format:      "text",
			FileEnabled: false,
			FilePath:    "logs/reactedge.log",
			Rotation: LogRotationConfig{
				MaxSize:    100,
				MaxBackups: 5,
				MaxAge:     30,
			},
		},
		Monitoring: MonitoringConfig{
			Enabled:         true,
			Port:            "6060",
			HealthCheckPath: "/health",
			Metrics: MetricsConfig{
				PrometheusEnabled: true,
				Path:              "/metrics",
			},
		},
		Development: DevelopmentConfig{
			Debug:          false,
			CORSEnabled:    true,
			CORSOrigins:    []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:6000"},
			RequestLogging: true,
			SQLLogging:     false,
		},
		Production: ProductionConfig{
			GzipEnabled:            true,
			RateLimitingEnabled:    true,
			StaticCacheTTL:         86400,
			SecurityHeadersEnabled: true,
		},
	}
}

// overrideFromEnv 环境变量覆盖配置
func overrideFromEnv(config *Config) {
	// 服务器配置
	if port := os.Getenv("SERVER_PORT"); port != "" {
		config.Server.Port = port
	}
	if host := os.Getenv("SERVER_HOST"); host != "" {
		config.Server.Host = host
	}
	if readTimeout := getEnvAsInt("SERVER_READ_TIMEOUT", 0); readTimeout > 0 {
		config.Server.ReadTimeout = readTimeout
	}
	if writeTimeout := getEnvAsInt("SERVER_WRITE_TIMEOUT", 0); writeTimeout > 0 {
		config.Server.WriteTimeout = writeTimeout
	}

	// AI配置
	if aiMode := os.Getenv("AI_MODE"); aiMode != "" {
		config.AI.Mode = aiMode
	}
	if maxAnalysisTime := getEnvAsInt("AI_MAX_ANALYSIS_TIME", 0); maxAnalysisTime > 0 {
		config.AI.MaxAnalysisTime = maxAnalysisTime
	}
	if cacheEnabled := os.Getenv("AI_CACHE_ENABLED"); cacheEnabled != "" {
		config.AI.CacheEnabled = getEnvAsBool("AI_CACHE_ENABLED", true)
	}

	// 日志配置
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.Logging.Level = logLevel
	}
	if logFormat := os.Getenv("LOG_FORMAT"); logFormat != "" {
		config.Logging.Format = logFormat
	}

	// 开发环境配置
	if debug := os.Getenv("DEBUG"); debug != "" {
		config.Development.Debug = getEnvAsBool("DEBUG", false)
	}
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	// 验证端口范围 (6000-6999)
	port, err := strconv.Atoi(config.Server.Port)
	if err != nil {
		return fmt.Errorf("服务器端口必须是数字: %s", config.Server.Port)
	}
	if port < 6000 || port > 6999 {
		return fmt.Errorf("服务器端口必须在6000-6999范围内: %d", port)
	}

	// 验证AI模式
	if config.AI.Mode != "internal" && config.AI.Mode != "external" {
		return fmt.Errorf("AI模式必须是 'internal' 或 'external': %s", config.AI.Mode)
	}

	// 验证超时配置
	if config.Server.ReadTimeout < 1 {
		config.Server.ReadTimeout = 30
	}
	if config.Server.WriteTimeout < 1 {
		config.Server.WriteTimeout = 30
	}

	return nil
}
