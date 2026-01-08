package ai

import (
	"context"
	"fmt"
	"sync"
)

// Manager AI服务管理器
type Manager struct {
	config    *Config
	client    Client
	providers map[ProviderType]Client
	mutex     sync.RWMutex
}

// NewManager 创建AI服务管理器
func NewManager(configPath string) (*Manager, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载AI配置失败: %w", err)
	}

	if err := config.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("AI配置验证失败: %w", err)
	}

	manager := &Manager{
		config:    config,
		providers: make(map[ProviderType]Client),
	}

	// 初始化可用的AI客户端
	if err := manager.initializeClients(); err != nil {
		return nil, fmt.Errorf("初始化AI客户端失败: %w", err)
	}

	// 设置默认客户端
	defaultProvider := ProviderType(config.DefaultProvider)
	if client, exists := manager.providers[defaultProvider]; exists {
		manager.client = client
	} else {
		// 如果默认服务商不可用，使用第一个可用的客户端
		for _, client := range manager.providers {
			manager.client = client
			fmt.Printf("⚠️ 默认服务商%s不可用，使用%s作为默认客户端\n", defaultProvider, client.GetProvider())
			break
		}
	}

	fmt.Printf("✅ AI服务管理器初始化完成，默认服务商: %s\n", manager.client.GetProvider())
	fmt.Printf("📊 可用服务商: ")
	for provider := range manager.providers {
		fmt.Printf("%s ", provider)
	}
	fmt.Println()

	return manager, nil
}

// initializeClients 初始化所有可用的AI客户端
func (m *Manager) initializeClients() error {
	availableProviders := m.config.GetAvailableProviders()

	for _, provider := range availableProviders {
		client, err := NewClient(provider, m.config)
		if err != nil {
			fmt.Printf("⚠️ 初始化%s客户端失败: %v\n", provider, err)
			continue
		}

		m.providers[provider] = client
		fmt.Printf("✅ %s客户端初始化成功，支持%d个模型\n", provider, len(client.GetAvailableModels()))
	}

	if len(m.providers) == 0 {
		return fmt.Errorf("没有可用的AI服务商，请检查配置和环境变量")
	}

	return nil
}

// GetClient 获取当前默认客户端
func (m *Manager) GetClient() Client {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.client
}

// GetClientByProvider 根据服务商获取客户端
func (m *Manager) GetClientByProvider(provider ProviderType) (Client, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if client, exists := m.providers[provider]; exists {
		return client, nil
	}

	return nil, fmt.Errorf("服务商%s不可用", provider)
}

// SwitchProvider 切换默认服务商
func (m *Manager) SwitchProvider(provider ProviderType) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if client, exists := m.providers[provider]; exists {
		m.client = client
		fmt.Printf("✅ 已切换到%s服务商\n", provider)
		return nil
	}

	return fmt.Errorf("服务商%s不可用", provider)
}

// GetAvailableProviders 获取所有可用服务商
func (m *Manager) GetAvailableProviders() []ProviderType {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	providers := make([]ProviderType, 0, len(m.providers))
	for provider := range m.providers {
		providers = append(providers, provider)
	}

	return providers
}

// GetConfig 获取配置
func (m *Manager) GetConfig() *Config {
	return m.config
}

// AnalyzeImage 图像分析（使用默认客户端）
func (m *Manager) AnalyzeImage(ctx context.Context, imageURL, prompt string) (*ImageAnalysisResult, error) {
	return m.client.AnalyzeImage(ctx, imageURL, prompt)
}

// GenerateQuestions 生成问题（使用默认客户端）
func (m *Manager) GenerateQuestions(ctx context.Context, contextInfo string, category string) ([]Question, error) {
	return m.client.GenerateQuestions(ctx, contextInfo, category)
}

// PolishNote 润色笔记（使用默认客户端）
func (m *Manager) PolishNote(ctx context.Context, rawContent, contextInfo string) (*PolishedNote, error) {
	return m.client.PolishNote(ctx, rawContent, contextInfo)
}

// TextToSpeech 文字转语音（使用默认客户端）
func (m *Manager) TextToSpeech(ctx context.Context, text, voice, language string, speed float64) ([]byte, string, error) {
	return m.client.TextToSpeech(ctx, text, voice, language, speed)
}

// AnalyzeVideo 视频分析（使用默认客户端）
func (m *Manager) AnalyzeVideo(ctx context.Context, videoData []byte, format, analysisType string, duration float64) (*VideoAnalysis, error) {
	return m.client.AnalyzeVideo(ctx, videoData, format, analysisType, duration)
}

// GenerateVideo 视频生成（使用默认客户端）
func (m *Manager) GenerateVideo(ctx context.Context, script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	return m.client.GenerateVideo(ctx, script, style, duration, scenes, voice, language)
}

// GenerateReactionTemplates 生成反应模板（使用默认客户端）
func (m *Manager) GenerateReactionTemplates(ctx context.Context, scenario, style string) ([]ReactionTemplate, error) {
	return m.client.GenerateReactionTemplates(ctx, scenario, style)
}

// AnalyzeExpressionStyle 分析表达风格（使用默认客户端）
func (m *Manager) AnalyzeExpressionStyle(ctx context.Context, personName string, sampleText string) (*StyleAnalysis, error) {
	return m.client.AnalyzeExpressionStyle(ctx, personName, sampleText)
}

// SimulateDebate 模拟辩论（使用默认客户端）
func (m *Manager) SimulateDebate(ctx context.Context, scenario string, difficulty int, userStyle string) (*DebateSimulation, error) {
	return m.client.SimulateDebate(ctx, scenario, difficulty, userStyle)
}

// EvaluateReaction 评估反应（使用默认客户端）
func (m *Manager) EvaluateReaction(ctx context.Context, userResponse, scenario, expectedStyle string) (*ReactionEvaluation, error) {
	return m.client.EvaluateReaction(ctx, userResponse, scenario, expectedStyle)
}

// ReactEdge增强功能

// GeneratePersonalizedTraining 生成个性化训练计划
func (m *Manager) GeneratePersonalizedTraining(ctx context.Context, userProfile map[string]interface{}, currentLevel int) (*PersonalizedTraining, error) {
	prompt := fmt.Sprintf(`基于用户画像生成个性化临场反应训练计划：

用户画像：%v
当前等级：%d

请生成包含以下内容的训练计划：
1. 主要训练重点
2. 推荐的训练场景
3. 难度递进建议
4. 每周训练安排
5. 预期效果评估

返回JSON格式的训练计划。`, userProfile, currentLevel)

	// 使用高级推理模型生成训练计划
	client := m.GetClient()
	questions, err := client.GenerateQuestions(ctx, fmt.Sprintf("训练计划生成：%s", prompt), "training")
	if err != nil {
		return m.getDefaultPersonalizedTraining(userProfile, currentLevel), nil
	}

	// 这里简化为基于问题的训练计划
	training := &PersonalizedTraining{
		UserLevel:      currentLevel,
		MainFocus:      []string{"反应速度", "内容质量", "风格适应"},
		RecommendedScenarios: []string{"述职答辩", "分享会提问", "争辩冲突"},
		WeeklyPlan:     m.generateWeeklyPlan(currentLevel),
		ExpectedOutcomes: []string{"提升反应速度20%", "增强内容逻辑性", "掌握多种沟通风格"},
	}

	return training, nil
}

// generateWeeklyPlan 生成每周训练计划
func (m *Manager) generateWeeklyPlan(level int) []WeeklySession {
	sessions := []WeeklySession{}

	switch level {
	case 1:
		sessions = []WeeklySession{
			{Day: 1, Focus: "基础反应训练", Duration: 15, Scenarios: []string{"简单问答"}},
			{Day: 2, Focus: "风格适应训练", Duration: 20, Scenarios: []string{"正式场合"}},
			{Day: 3, Focus: "压力模拟训练", Duration: 15, Scenarios: []string{"时间限制"}},
			{Day: 4, Focus: "反馈分析", Duration: 10, Scenarios: []string{"自我评估"}},
			{Day: 5, Focus: "综合训练", Duration: 25, Scenarios: []string{"混合场景"}},
		}
	case 2:
		sessions = []WeeklySession{
			{Day: 1, Focus: "高级反应训练", Duration: 20, Scenarios: []string{"复杂问题"}},
			{Day: 2, Focus: "多风格切换", Duration: 25, Scenarios: []string{"不同场合"}},
			{Day: 3, Focus: "辩论模拟", Duration: 30, Scenarios: []string{"观点冲突"}},
			{Day: 4, Focus: "实时反馈", Duration: 15, Scenarios: []string{"AI评估"}},
			{Day: 5, Focus: "挑战训练", Duration: 35, Scenarios: []string{"高难度场景"}},
		}
	default:
		sessions = []WeeklySession{
			{Day: 1, Focus: "专家级训练", Duration: 30, Scenarios: []string{"专业辩论"}},
			{Day: 2, Focus: "危机应对", Duration: 35, Scenarios: []string{"紧急情况"}},
			{Day: 3, Focus: "领导沟通", Duration: 30, Scenarios: []string{"高层对话"}},
			{Day: 4, Focus: "公众演讲", Duration: 40, Scenarios: []string{"大型会议"}},
			{Day: 5, Focus: "大师挑战", Duration: 45, Scenarios: []string{"终极考验"}},
		}
	}

	return sessions
}

// getDefaultPersonalizedTraining 默认个性化训练计划
func (m *Manager) getDefaultPersonalizedTraining(userProfile map[string]interface{}, level int) *PersonalizedTraining {
	return &PersonalizedTraining{
		UserLevel:      level,
		MainFocus:      []string{"反应速度提升", "内容质量优化", "风格适应性"},
		RecommendedScenarios: []string{"述职报告", "分享会", "团队讨论"},
		WeeklyPlan:     m.generateWeeklyPlan(level),
		ExpectedOutcomes: []string{"显著提升临场反应能力", "增强沟通自信心", "掌握多种应对策略"},
	}
}

// 数据结构

// PersonalizedTraining 个性化训练计划
type PersonalizedTraining struct {
	UserLevel           int             `json:"user_level"`
	MainFocus           []string        `json:"main_focus"`
	RecommendedScenarios []string       `json:"recommended_scenarios"`
	WeeklyPlan          []WeeklySession `json:"weekly_plan"`
	ExpectedOutcomes    []string        `json:"expected_outcomes"`
}

// WeeklySession 每周训练 session
type WeeklySession struct {
	Day       int      `json:"day"`
	Focus     string   `json:"focus"`
	Duration  int      `json:"duration"` // 分钟
	Scenarios []string `json:"scenarios"`
}
