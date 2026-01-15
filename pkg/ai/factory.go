package ai

import (
	"fmt"
)

// AIFactory AI客户端工厂
type AIFactory struct {
	config *Config
}

// NewAIFactory 创建AI工厂
func NewAIFactory(config *Config) *AIFactory {
	return &AIFactory{
		config: config,
	}
}

// CreateClient 根据AI模式创建相应的客户端
// 对内环境：使用TAL AI模型
// 对外环境：使用开放AI模型（如OpenAI、Claude等）
func (f *AIFactory) CreateClient() (Client, error) {
	mode := f.config.GetAIMode()

	switch mode {
	case "internal":
		// 对内环境：优先使用TAL AI模型
		if f.isTALAvailable() {
			return NewTALClient(f.config.TAL)
		}
		// TAL不可用时降级到其他服务商
		return f.createFallbackClient()

	case "external":
		// 对外环境：使用开放AI模型
		availableProviders := f.config.GetAvailableProviders()

		// 优先级：OpenAI -> Claude -> Azure -> Baidu
		for _, provider := range []ProviderType{ProviderOpenAI, ProviderClaude, ProviderAzure, ProviderBaidu} {
			for _, available := range availableProviders {
				if available == provider {
					return NewClient(provider, f.config)
				}
			}
		}

		// 如果没有可用的外部服务商，使用本地模拟
		return f.createFallbackClient()

	default:
		return f.createFallbackClient()
	}
}

// isTALAvailable 检查TAL服务是否可用
func (f *AIFactory) isTALAvailable() bool {
	return f.config.TAL.TAL_MLOPS_APP_ID != "" && f.config.TAL.TAL_MLOPS_APP_KEY != ""
}

// createFallbackClient 创建降级客户端
func (f *AIFactory) createFallbackClient() (Client, error) {
	// 尝试创建任何可用的客户端
	availableProviders := f.config.GetAvailableProviders()

	if len(availableProviders) > 0 {
		return NewClient(availableProviders[0], f.config)
	}

	// 如果没有任何可用的客户端，返回错误
	return nil, fmt.Errorf("没有可用的AI服务商，请检查配置和环境变量")
}

// GetRecommendedProvider 根据任务类型推荐服务商
func (f *AIFactory) GetRecommendedProvider(taskType string) ProviderType {
	mode := f.config.GetAIMode()

	if mode == "internal" && f.isTALAvailable() {
		return ProviderTAL
	}

	// 根据任务类型推荐最适合的服务商
	switch taskType {
	case "advanced_reasoning":
		// 复杂推理任务优先使用TAL的qwen3-max
		if mode == "internal" && f.isTALAvailable() {
			return ProviderTAL
		}
		return ProviderOpenAI // 或者Claude

	case "image_analysis":
		// 图像分析任务
		if mode == "internal" && f.isTALAvailable() {
			return ProviderTAL // qwen3-vl-plus
		}
		return ProviderOpenAI // GPT-4o

	case "text_generation":
		// 文本生成任务
		if mode == "internal" && f.isTALAvailable() {
			return ProviderTAL // qwen-flash
		}
		return ProviderOpenAI // GPT-4

	default:
		return f.config.GetAvailableProviders()[0]
	}
}

// GetClientForTask 根据任务类型创建最适合的客户端
func (f *AIFactory) GetClientForTask(taskType string) (Client, error) {
	provider := f.GetRecommendedProvider(taskType)
	return NewClient(provider, f.config)
}
