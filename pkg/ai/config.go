package ai

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config AI服务配置
type Config struct {
	// AI模式：internal(对内) 或 external(对外)
	AIMode string `json:"aiMode" yaml:"aiMode"`

	// 默认服务商
	DefaultProvider string `json:"defaultProvider" yaml:"defaultProvider"`

	// OpenAI兼容服务配置
	OpenAI OpenAIConfig `json:"openai" yaml:"openai"`

	// Claude配置
	Claude ClaudeConfig `json:"claude" yaml:"claude"`

	// 其他服务商配置可以在这里扩展
	Azure   AzureConfig   `json:"azure" yaml:"azure"`
	Baidu   BaiduConfig   `json:"baidu" yaml:"baidu"`
	TAL     TALConfig     `json:"tal" yaml:"tal"`
	Spark   SparkConfig   `json:"spark" yaml:"spark"`
}

// OpenAIConfig OpenAI兼容服务配置
type OpenAIConfig struct {
	APIKey     string  `json:"apiKey" yaml:"apiKey"`
	BaseURL    string  `json:"baseURL" yaml:"baseURL"`
	Timeout    int     `json:"timeout" yaml:"timeout"`
	MaxTokens  int     `json:"maxTokens" yaml:"maxTokens"`
	Temperature float32 `json:"temperature" yaml:"temperature"`

	// 模型映射
	Models ModelMapping `json:"models" yaml:"models"`
}

// ClaudeConfig Claude配置
type ClaudeConfig struct {
	APIKey     string  `json:"apiKey" yaml:"apiKey"`
	BaseURL    string  `json:"baseURL" yaml:"baseURL"`
	Timeout    int     `json:"timeout" yaml:"timeout"`
	MaxTokens  int     `json:"maxTokens" yaml:"maxTokens"`
	Temperature float32 `json:"temperature" yaml:"temperature"`

	// 模型映射
	Models ModelMapping `json:"models" yaml:"models"`
}

// AzureConfig Azure AI配置
type AzureConfig struct {
	APIKey      string `json:"apiKey" yaml:"apiKey"`
	Endpoint    string `json:"endpoint" yaml:"endpoint"`
	Deployment  string `json:"deployment" yaml:"deployment"`
	APIVersion  string `json:"apiVersion" yaml:"apiVersion"`
	Timeout     int    `json:"timeout" yaml:"timeout"`
	MaxTokens   int    `json:"maxTokens" yaml:"maxTokens"`
	Temperature float32 `json:"temperature" yaml:"temperature"`
}

// BaiduConfig 百度AI配置
type BaiduConfig struct {
	APIKey     string  `json:"apiKey" yaml:"apiKey"`
	SecretKey  string  `json:"secretKey" yaml:"secretKey"`
	Timeout    int     `json:"timeout" yaml:"timeout"`
	MaxTokens  int     `json:"maxTokens" yaml:"maxTokens"`
	Temperature float32 `json:"temperature" yaml:"temperature"`
}

// TALConfig TAL内部AI服务配置
type TALConfig struct {
	TAL_MLOPS_APP_ID  string  `json:"talMLOpsAppId" yaml:"talMLOpsAppId"`
	TAL_MLOPS_APP_KEY string  `json:"talMLOpsAppKey" yaml:"talMLOpsAppKey"`
	BaseURL           string  `json:"baseURL" yaml:"baseURL"`
	Timeout           int     `json:"timeout" yaml:"timeout"`
	MaxTokens         int     `json:"maxTokens" yaml:"maxTokens"`
	Temperature       float32 `json:"temperature" yaml:"temperature"`

	// 模型映射
	Models ModelMapping `json:"models" yaml:"models"`
}

// SparkConfig 星火AI配置
type SparkConfig struct {
	AppID      string  `json:"appId" yaml:"appId"`
	APIKey     string  `json:"apiKey" yaml:"apiKey"`
	APISecret  string  `json:"apiSecret" yaml:"apiSecret"`
	Model      string  `json:"model" yaml:"model"`
	Timeout    int     `json:"timeout" yaml:"timeout"`
	MaxTokens  int     `json:"maxTokens" yaml:"maxTokens"`
	Temperature float32 `json:"temperature" yaml:"temperature"`
}

// ModelMapping 模型映射配置
type ModelMapping struct {
	ImageAnalysis     string `json:"imageAnalysis" yaml:"imageAnalysis"`
	TextGeneration    string `json:"textGeneration" yaml:"textGeneration"`
	AdvancedReasoning string `json:"advancedReasoning" yaml:"advancedReasoning"`
	VoiceInteraction  string `json:"voiceInteraction" yaml:"voiceInteraction"`
	VideoAnalysis     string `json:"videoAnalysis" yaml:"videoAnalysis"`
	VideoGeneration   string `json:"videoGeneration" yaml:"videoGeneration"`
}

// ProviderType AI服务商类型
type ProviderType string

const (
	ProviderOpenAI ProviderType = "openai"
	ProviderClaude ProviderType = "claude"
	ProviderAzure  ProviderType = "azure"
	ProviderBaidu  ProviderType = "baidu"
	ProviderTAL    ProviderType = "tal"
	ProviderSpark  ProviderType = "spark"
)

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		AIMode:          "internal",              // 默认对内模式
		DefaultProvider: string(ProviderTAL),     // 默认使用TAL内部服务
		TAL: TALConfig{
			BaseURL:     "http://ai-service.tal.com/openai-compatible/v1",
			Timeout:     30,
			MaxTokens:   2000,
			Temperature: 0.7,
			Models: ModelMapping{
				ImageAnalysis:     "gpt-4o",           // GPT-4o支持多模态
				TextGeneration:    "deepseek-chat",    // Deepseek Chat通用文本生成
				AdvancedReasoning: "deepseek-reasoner", // Deepseek Reasoner推理能力强
				VoiceInteraction:  "gpt-4o",           // GPT-4o支持多模态
				VideoAnalysis:     "gpt-4o",           // GPT-4o支持多模态
				VideoGeneration:   "doubao-pro-128k",  // Doubao模型支持
			},
		},
		OpenAI: OpenAIConfig{
			BaseURL:     "https://api.openai.com/v1",
			Timeout:     30,
			MaxTokens:   2000,
			Temperature: 0.7,
			Models: ModelMapping{
				ImageAnalysis:     "gpt-4o",
				TextGeneration:    "gpt-4",
				AdvancedReasoning: "gpt-4",
				VoiceInteraction:  "gpt-4o",
			},
		},
		Claude: ClaudeConfig{
			BaseURL:     "https://api.anthropic.com",
			Timeout:     30,
			MaxTokens:   2000,
			Temperature: 0.7,
			Models: ModelMapping{
				ImageAnalysis:     "claude-3-opus-20240229",
				TextGeneration:    "claude-3-haiku-20240307",
				AdvancedReasoning: "claude-3-opus-20240229",
				VoiceInteraction:  "claude-3-opus-20240229",
			},
		},
	}
}

// LoadConfig 从文件加载配置
func LoadConfig(configPath string) (*Config, error) {
	// 默认配置
	config := DefaultConfig()

	// 尝试从文件加载
	if file, err := os.Open(configPath); err == nil {
		defer file.Close()

		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(config); err != nil {
			return nil, fmt.Errorf("解析AI配置文件失败: %w", err)
		}
		fmt.Printf("✅ 从文件加载AI配置: %s\n", configPath)
	} else {
		// 如果主配置文件不存在，尝试加载示例配置
		examplePath := configPath + ".example"
		if exampleFile, exampleErr := os.Open(examplePath); exampleErr == nil {
			defer exampleFile.Close()

			decoder := yaml.NewDecoder(exampleFile)
			if err := decoder.Decode(config); err != nil {
				return nil, fmt.Errorf("解析AI示例配置文件失败: %w", err)
			}
			fmt.Printf("✅ 从示例配置文件加载AI配置: %s，请复制并修改为实际配置\n", examplePath)
		} else {
			fmt.Printf("⚠️ AI配置文件不存在，使用默认配置: %s\n", configPath)
		}
	}

	// 从环境变量加载敏感信息（环境变量优先级更高）
	loadFromEnv(config)

	// 验证配置
	if err := validateAIConfig(config); err != nil {
		return nil, fmt.Errorf("AI配置验证失败: %w", err)
	}

	fmt.Printf("✅ AI配置加载成功，默认服务商: %s，可用服务商: %v\n", config.DefaultProvider, config.GetAvailableProviders())
	return config, nil
}

// validateAIConfig 验证AI配置
func validateAIConfig(config *Config) error {
	// 验证AI模式
	if config.AIMode != "internal" && config.AIMode != "external" {
		return fmt.Errorf("AI模式必须是 'internal' 或 'external': %s", config.AIMode)
	}

	// 验证默认服务商
	validProviders := []string{string(ProviderTAL), string(ProviderOpenAI), string(ProviderClaude), string(ProviderAzure), string(ProviderBaidu), string(ProviderSpark)}
	isValid := false
	for _, provider := range validProviders {
		if config.DefaultProvider == provider {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("无效的默认服务商: %s，可选值: %s", config.DefaultProvider, strings.Join(validProviders, ", "))
	}

	return nil
}

// loadFromEnv 从环境变量加载配置
func loadFromEnv(config *Config) {
	// TAL配置
	if talAppID := os.Getenv("TAL_MLOPS_APP_ID"); talAppID != "" {
		config.TAL.TAL_MLOPS_APP_ID = talAppID
	}
	if talAppKey := os.Getenv("TAL_MLOPS_APP_KEY"); talAppKey != "" {
		config.TAL.TAL_MLOPS_APP_KEY = talAppKey
	}

	// OpenAI配置
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey != "" {
		config.OpenAI.APIKey = openaiKey
	}

	// Claude配置
	if claudeKey := os.Getenv("ANTHROPIC_API_KEY"); claudeKey != "" {
		config.Claude.APIKey = claudeKey
	}

	// Azure配置
	if azureKey := os.Getenv("AZURE_OPENAI_API_KEY"); azureKey != "" {
		config.Azure.APIKey = azureKey
	}
	if azureEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT"); azureEndpoint != "" {
		config.Azure.Endpoint = azureEndpoint
	}

	// 百度配置
	if baiduKey := os.Getenv("BAIDU_API_KEY"); baiduKey != "" {
		config.Baidu.APIKey = baiduKey
	}
	if baiduSecret := os.Getenv("BAIDU_SECRET_KEY"); baiduSecret != "" {
		config.Baidu.SecretKey = baiduSecret
	}
}

// GetProviderConfig 获取指定服务商的配置
func (c *Config) GetProviderConfig(provider ProviderType) interface{} {
	switch provider {
	case ProviderOpenAI:
		return c.OpenAI
	case ProviderClaude:
		return c.Claude
	case ProviderAzure:
		return c.Azure
	case ProviderBaidu:
		return c.Baidu
	case ProviderTAL:
		return c.TAL
	default:
		return nil
	}
}

// ValidateConfig 验证配置有效性
func (c *Config) ValidateConfig() error {
	if c.DefaultProvider == "" {
		return fmt.Errorf("defaultProvider不能为空")
	}

	// 检查默认服务商是否有效
	validProviders := []string{string(ProviderOpenAI), string(ProviderClaude), string(ProviderAzure), string(ProviderBaidu), string(ProviderTAL)}
	isValid := false
	for _, provider := range validProviders {
		if c.DefaultProvider == provider {
			isValid = true
			break
		}
	}
	if !isValid {
		return fmt.Errorf("无效的defaultProvider: %s，可选值: %s", c.DefaultProvider, strings.Join(validProviders, ", "))
	}

	// 验证TAL配置
	if c.TAL.TAL_MLOPS_APP_ID == "" || c.TAL.TAL_MLOPS_APP_KEY == "" {
		fmt.Printf("⚠️ TAL MLOps配置不完整，将使用其他服务商\n")
	}

	// 验证OpenAI配置
	if c.OpenAI.APIKey == "" {
		fmt.Printf("⚠️ OpenAI API Key未配置\n")
	}

	// 验证Claude配置
	if c.Claude.APIKey == "" {
		fmt.Printf("⚠️ Claude API Key未配置\n")
	}

	return nil
}

// GetModelForTask 根据任务类型获取模型名称
func (c *Config) GetModelForTask(provider ProviderType, task string) string {
	var models ModelMapping

	switch provider {
	case ProviderTAL:
		models = c.TAL.Models
	case ProviderOpenAI:
		models = c.OpenAI.Models
	case ProviderClaude:
		models = c.Claude.Models
	default:
		return ""
	}

	switch task {
	case "image_analysis":
		return models.ImageAnalysis
	case "text_generation":
		return models.TextGeneration
	case "advanced_reasoning":
		return models.AdvancedReasoning
	case "voice_interaction":
		return models.VoiceInteraction
	case "video_analysis":
		return models.VideoAnalysis
	case "video_generation":
		return models.VideoGeneration
	default:
		return models.TextGeneration // 默认使用文本生成模型
	}
}

// GetAvailableProviders 获取所有可用的服务商列表
func (c *Config) GetAvailableProviders() []ProviderType {
	providers := []ProviderType{}

	// 检查各个服务商是否配置完整
	if c.TAL.TAL_MLOPS_APP_ID != "" && c.TAL.TAL_MLOPS_APP_KEY != "" {
		providers = append(providers, ProviderTAL)
	}
	if c.Spark.AppID != "" && c.Spark.APIKey != "" && c.Spark.APISecret != "" {
		providers = append(providers, ProviderSpark)
	}
	if c.OpenAI.APIKey != "" {
		providers = append(providers, ProviderOpenAI)
	}
	if c.Claude.APIKey != "" {
		providers = append(providers, ProviderClaude)
	}
	if c.Azure.APIKey != "" && c.Azure.Endpoint != "" {
		providers = append(providers, ProviderAzure)
	}
	if c.Baidu.APIKey != "" && c.Baidu.SecretKey != "" {
		providers = append(providers, ProviderBaidu)
	}

	return providers
}

// GetAIMode 获取AI模式
func (c *Config) GetAIMode() string {
	if c.AIMode == "" {
		return "internal" // 默认对内模式
	}
	return c.AIMode
}
