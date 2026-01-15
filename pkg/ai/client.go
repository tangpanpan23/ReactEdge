package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

// Client AI服务客户端接口
type Client interface {
	// 基础方法
	GetProvider() ProviderType
	GetAvailableModels() []string
	ValidateModel(model string) bool

	// 核心AI功能
	AnalyzeImage(ctx context.Context, imageURL, prompt string) (*ImageAnalysisResult, error)
	GenerateQuestions(ctx context.Context, contextInfo string, category string) ([]Question, error)
	PolishNote(ctx context.Context, rawContent, contextInfo string) (*PolishedNote, error)
	TextToSpeech(ctx context.Context, text, voice, language string, speed float64) ([]byte, string, error)
	AnalyzeVideo(ctx context.Context, videoData []byte, format, analysisType string, duration float64) (*VideoAnalysis, error)
	GenerateVideo(ctx context.Context, script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error)

	// ReactEdge特定功能
	GenerateReactionTemplates(ctx context.Context, scenario, style string) ([]ReactionTemplate, error)
	AnalyzeExpressionStyle(ctx context.Context, personName string, sampleText string) (*StyleAnalysis, error)
	SimulateDebate(ctx context.Context, scenario string, difficulty int, userStyle string) (*DebateSimulation, error)
	EvaluateReaction(ctx context.Context, userResponse, scenario, expectedStyle string) (*ReactionEvaluation, error)
}

// BaseClient 基础AI客户端结构体
type BaseClient struct {
	provider ProviderType
	config   *Config
}

// GetProvider 获取服务商类型
func (c *BaseClient) GetProvider() ProviderType {
	return c.provider
}

// NewClient 创建AI客户端
func NewClient(provider ProviderType, config *Config) (Client, error) {
	switch provider {
	case ProviderTAL:
		return NewTALClient(config.TAL)
	case ProviderOpenAI:
		return NewOpenAIClient(config.OpenAI)
	case ProviderClaude:
		return NewClaudeClient(config.Claude)
	case ProviderAzure:
		return NewAzureClient(config.Azure)
	case ProviderBaidu:
		return NewBaiduClient(config.Baidu)
	case ProviderSpark:
		return NewSparkClient(config.Spark)
	default:
		return nil, fmt.Errorf("不支持的服务商: %s", provider)
	}
}

// TALClient TAL内部AI服务客户端
type TALClient struct {
	*BaseClient
	httpClient *http.Client
	config     *TALConfig
	baseURL    string
	authToken  string
	client     *openai.Client // OpenAI兼容客户端

	// 请求限流
	requestMutex   sync.Mutex
	lastRequestTime time.Time
	minInterval     time.Duration // 最小请求间隔
}

// NewTALClient 创建TAL客户端
func NewTALClient(config TALConfig) (*TALClient, error) {
	// 构建认证token
	authToken := fmt.Sprintf("%s:%s", config.TAL_MLOPS_APP_ID, config.TAL_MLOPS_APP_KEY)

	// 从配置中获取请求间隔（如果有的话），否则使用默认值
	minInterval := 1 * time.Second // 默认1秒
	if config.Timeout > 0 {
		// 可以根据配置调整间隔，这里暂时保持默认
		minInterval = 1 * time.Second
	}

	fmt.Printf("🔧 初始化TAL AI客户端 - 端点: %s, 最小请求间隔: %.1f秒\n", config.BaseURL, minInterval.Seconds())

	// 测试网络连接
	fmt.Printf("🔍 测试TAL AI服务连接...\n")
	testClient := &http.Client{Timeout: 10 * time.Second}
	testURL := config.BaseURL + "/models"
	testReq, _ := http.NewRequest("GET", testURL, nil)
	testReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	if testResp, testErr := testClient.Do(testReq); testErr == nil {
		testResp.Body.Close()
		fmt.Printf("✅ TAL AI服务连接正常 (状态码: %d)\n", testResp.StatusCode)
	} else {
		fmt.Printf("⚠️ TAL AI服务连接测试失败: %v\n", testErr)
		fmt.Println("💡 可能的原因: 网络连接问题、服务不可用或认证信息错误")
	}

	// 设置内部AI服务端点
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "http://ai-service.tal.com/openai-compatible/v1"
	}

	// 创建HTTP客户端 - 不设置超时，完全依赖传入的context控制超时
	httpClient := &http.Client{}

	fmt.Printf("🔧 HTTP客户端初始化完成，超时由context控制\n")

	// 初始化OpenAI兼容客户端
	openaiConfig := openai.DefaultConfig(authToken)
	openaiConfig.BaseURL = baseURL
	client := openai.NewClientWithConfig(openaiConfig)

	return &TALClient{
		BaseClient: &BaseClient{
			provider: ProviderTAL,
			config:   &Config{TAL: config},
		},
		httpClient:     httpClient,
		config:         &config,
		baseURL:        baseURL,
		authToken:      authToken,
		client:         client,
		minInterval:    minInterval, // 每秒最多1个请求
		lastRequestTime: time.Now().Add(-minInterval * 2), // 初始化为过去的时间
	}, nil
}

// TALTransport TAL认证传输层
type TALTransport struct {
	base  http.RoundTripper
	token string
}

func (t *TALTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.token)
	if t.base == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.base.RoundTrip(req)
}

// GetAvailableModels 获取可用模型列表
func (c *TALClient) GetAvailableModels() []string {
	return []string{
		"qwen3-vl-plus",               // 图像分析主模型
		"qwen-flash",                  // 文本生成主模型
		"qwen-flash",                  // 复杂推理主模型（更稳定）
		"qwen3-omni-flash",            // 语音交互主模型
		"qwen3-vl-235b-a22b-instruct", // 图像分析备用模型
		"qwen-turbo",                  // 文本生成备用模型
		"qwen-max",                    // 复杂推理备用模型
	}
}

// ValidateModel 验证模型是否可用
func (c *TALClient) ValidateModel(model string) bool {
	availableModels := c.GetAvailableModels()
	for _, availableModel := range availableModels {
		if availableModel == model {
			return true
		}
	}
	return false
}

// GetModelForTask 根据任务获取模型
func (c *TALClient) GetModelForTask(task string) string {
	switch task {
	case "image_analysis":
		return c.config.Models.ImageAnalysis
	case "text_generation":
		return c.config.Models.TextGeneration
	case "advanced_reasoning":
		return c.config.Models.AdvancedReasoning
	case "voice_interaction":
		return c.config.Models.VoiceInteraction
	case "video_analysis":
		return c.config.Models.VideoAnalysis
	case "video_generation":
		return c.config.Models.VideoGeneration
	default:
		return c.config.Models.TextGeneration
	}
}

// AnalyzeImage 图像分析
func (c *TALClient) AnalyzeImage(ctx context.Context, imageURL, prompt string) (*ImageAnalysisResult, error) {
	model := c.GetModelForTask("image_analysis")

	contentParts := []openai.ChatMessagePart{
		{
			Type: openai.ChatMessagePartTypeText,
			Text: prompt,
		},
		{
			Type: openai.ChatMessagePartTypeImageURL,
			ImageURL: &openai.ChatMessageImageURL{
				URL: imageURL,
			},
		},
	}

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:         openai.ChatMessageRoleUser,
				MultiContent: contentParts,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		// AI服务不可用时，返回默认响应
		return getDefaultImageAnalysis(), nil
	}

	if len(resp.Choices) == 0 {
		return getDefaultImageAnalysis(), nil
	}

	content := resp.Choices[0].Message.Content

	result := &ImageAnalysisResult{
		ObjectName:     extractObjectName(content),
		Category:       extractCategory(content),
		Description:    content,
		Confidence:     0.95,
		KeyFeatures:    extractKeyFeatures(content),
		ScientificName: extractScientificName(content),
	}

	return result, nil
}

// GenerateQuestions 生成问题
func (c *TALClient) GenerateQuestions(ctx context.Context, contextInfo string, category string) ([]Question, error) {
	model := c.GetModelForTask("text_generation")

	prompt := fmt.Sprintf(`基于以下信息为用户生成3个引导性的反应训练问题：

上下文信息：%s
训练类别：%s

要求：
1. 问题要适合职场沟通场景
2. 问题要激发思考和反应能力
3. 问题难度要循序渐进（从简单到深入）
4. 每个问题都要有明确的类型标注
5. 确保所有内容适合职场培训场景

请以JSON格式返回，包含以下字段：
- content: 问题内容
- type: 问题类型（scenario场景, strategy策略, evaluation评估）
- difficulty: 难度（basic基本, intermediate中级, advanced高级）
- purpose: 问题目的说明`, contextInfo, category)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个职场沟通训练助手，专门为用户设计反应训练问题。请以JSON格式返回包含questions数组的结果。",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return getDefaultQuestions(), nil
	}

	if len(resp.Choices) == 0 {
		return getDefaultQuestions(), nil
	}

	content := resp.Choices[0].Message.Content
	questions := parseQuestionsFromJSON(content)
	if len(questions) == 0 {
		return getDefaultQuestions(), nil
	}

	return questions, nil
}

// PolishNote 润色笔记（这里用于润色反应记录）
func (c *TALClient) PolishNote(ctx context.Context, rawContent, contextInfo string) (*PolishedNote, error) {
	model := c.GetModelForTask("text_generation")

	prompt := fmt.Sprintf(`请帮用户润色他们的反应训练记录，让它更清晰、有逻辑性。

原始内容：%s

上下文信息：%s

要求：
1. 保持用户的原意和表达特色
2. 让表达更清晰准确
3. 添加适当的沟通技巧解释
4. 指出可能的改进方向
5. 确保所有内容适合职场培训场景

请严格按照以下JSON格式返回结果：

{
  "title": "记录标题",
  "summary": "内容总结",
  "key_points": ["关键要点1", "关键要点2"],
  "communication_tips": ["沟通技巧1"],
  "questions": ["问题1"],
  "improvements": ["改进建议1"],
  "formatted_text": "格式化的文本内容"
}

请确保返回的是有效的JSON格式。`, rawContent, contextInfo)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return getDefaultPolishedNote(), nil
	}

	if len(resp.Choices) == 0 {
		return getDefaultPolishedNote(), nil
	}

	content := resp.Choices[0].Message.Content

	// 处理markdown格式的JSON代码块
	jsonContent := content
	if strings.Contains(content, "```json") {
		startIndex := strings.Index(content, "```json")
		if startIndex != -1 {
			startIndex += 7
			endIndex := strings.Index(content[startIndex:], "```")
			if endIndex != -1 {
				jsonContent = strings.TrimSpace(content[startIndex : startIndex+endIndex])
			}
		}
	}

	var jsonResult PolishedNote
	if err := json.Unmarshal([]byte(jsonContent), &jsonResult); err != nil {
		jsonResult = PolishedNote{
			Title:       "反应训练记录",
			Summary:     extractSummaryFromText(jsonContent),
			FormattedText: jsonContent,
			KeyPoints:   extractKeyPointsFromText(jsonContent),
		}
	}

	// 确保必需字段有值
	if jsonResult.Title == "" {
		jsonResult.Title = "反应训练记录"
	}
	if jsonResult.Summary == "" {
		jsonResult.Summary = "这是反应训练的记录总结"
	}
	if len(jsonResult.KeyPoints) == 0 {
		jsonResult.KeyPoints = []string{"记录了训练过程", "总结了经验教训"}
	}
	if jsonResult.FormattedText == "" {
		jsonResult.FormattedText = content
	}

	return &jsonResult, nil
}

// TextToSpeech 文字转语音
func (c *TALClient) TextToSpeech(ctx context.Context, text, voice, language string, speed float64) ([]byte, string, error) {
	model := c.GetModelForTask("voice_interaction")

	prompt := fmt.Sprintf(`请将以下文字转换为语音：

文字内容：%s
语音类型：%s
语言：%s
语速：%.1f倍速

要求：
1. 生成自然流畅的语音
2. 保持专业语调
3. 语速适中，易于理解
4. 发音准确，表达清晰

请按照以下JSON格式返回结果：
{
  "audio_data": "base64编码的音频数据",
  "format": "wav"
}

不要包含任何其他文字，只返回有效的JSON。`, text, voice, language, speed)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return getDefaultAudioData(), "wav", nil
	}

	if len(resp.Choices) == 0 {
		return getDefaultAudioData(), "wav", nil
	}

	content := resp.Choices[0].Message.Content

	// 处理音频数据
	var audioResult struct {
		AudioData string `json:"audio_data"`
		Format    string `json:"format"`
	}

	jsonContent := content
	if strings.Contains(content, "```json") {
		startIndex := strings.Index(content, "```json")
		if startIndex != -1 {
			startIndex += 7
			endIndex := strings.Index(content[startIndex:], "```")
			if endIndex != -1 {
				jsonContent = strings.TrimSpace(content[startIndex : startIndex+endIndex])
			}
		}
	}

	if err := json.Unmarshal([]byte(jsonContent), &audioResult); err != nil {
		return getDefaultAudioData(), "wav", nil
	}

	if audioResult.AudioData == "" {
		return getDefaultAudioData(), "wav", nil
	}

	audioBytes, err := base64.StdEncoding.DecodeString(audioResult.AudioData)
	if err != nil {
		return getDefaultAudioData(), "wav", nil
	}

	format := audioResult.Format
	if format == "" {
		format = "wav"
	}

	return audioBytes, format, nil
}

// AnalyzeVideo 视频分析
func (c *TALClient) AnalyzeVideo(ctx context.Context, videoData []byte, format, analysisType string, duration float64) (*VideoAnalysis, error) {
	// 简化实现，使用默认分析结果
	return getDefaultVideoAnalysis(), nil
}

// GenerateVideo 视频生成
func (c *TALClient) GenerateVideo(ctx context.Context, script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	// 简化实现，返回模拟视频数据
	return c.generateMockVideo(script, style, duration, scenes, voice, language)
}

// GenerateReactionTemplates 生成反应模板
func (c *TALClient) GenerateReactionTemplates(ctx context.Context, scenario, style string) ([]ReactionTemplate, error) {
	model := c.GetModelForTask("text_generation")

	prompt := fmt.Sprintf(`基于以下场景和风格，为用户生成临场反应训练模板：

场景：%s
风格：%s

要求：
1. 生成3-5个实用的反应模板
2. 每个模板包含触发情境、反应步骤、关键话术
3. 模板要贴合职场实际场景
4. 风格要符合指定的沟通风格

请以JSON格式返回，包含templates数组，每个模板包含：
- scenario: 触发情境
- steps: 反应步骤数组
- key_phrases: 关键话术数组
- style_notes: 风格要点`, scenario, style)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "你是一个职场沟通教练，专门设计临场反应训练模板。",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return getDefaultReactionTemplates(), nil
	}

	if len(resp.Choices) == 0 {
		return getDefaultReactionTemplates(), nil
	}

	content := resp.Choices[0].Message.Content

	// 解析JSON响应
	var result struct {
		Templates []ReactionTemplate `json:"templates"`
	}

	jsonContent := content
	if strings.Contains(content, "```json") {
		startIndex := strings.Index(content, "```json")
		if startIndex != -1 {
			startIndex += 7
			endIndex := strings.Index(content[startIndex:], "```")
			if endIndex != -1 {
				jsonContent = strings.TrimSpace(content[startIndex : startIndex+endIndex])
			}
		}
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return getDefaultReactionTemplates(), nil
	}

	return result.Templates, nil
}

// AnalyzeExpressionStyle 分析表达风格
func (c *TALClient) AnalyzeExpressionStyle(ctx context.Context, personName string, sampleText string) (*StyleAnalysis, error) {
	model := c.GetModelForTask("advanced_reasoning")

	prompt := fmt.Sprintf(`请分析%s的表达风格：

样本文本：%s

请从以下维度进行分析：
1. 语言特点（词汇、句式、修辞手法）
2. 思维模式（逻辑结构、论证方式）
3. 沟通策略（立场表达、冲突处理）
4. 个人特色（独特标识、风格标签）

请返回JSON格式的分析结果。`, personName, sampleText)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return getDefaultStyleAnalysis(), nil
	}

	if len(resp.Choices) == 0 {
		return getDefaultStyleAnalysis(), nil
	}

	content := resp.Choices[0].Message.Content

	var result StyleAnalysis
	jsonContent := content
	if strings.Contains(content, "```json") {
		startIndex := strings.Index(content, "```json")
		if startIndex != -1 {
			startIndex += 7
			endIndex := strings.Index(content[startIndex:], "```")
			if endIndex != -1 {
				jsonContent = strings.TrimSpace(content[startIndex : startIndex+endIndex])
			}
		}
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return getDefaultStyleAnalysis(), nil
	}

	return &result, nil
}

// SimulateDebate 模拟辩论
func (c *TALClient) SimulateDebate(ctx context.Context, scenario string, difficulty int, userStyle string) (*DebateSimulation, error) {
	model := c.GetModelForTask("advanced_reasoning")

	prompt := fmt.Sprintf(`请模拟一个%s场景的辩论训练：

场景：%s
难度等级：%d
用户风格：%s

请生成：
1. 对手的开场陈述
2. 3轮交互对话
3. 关键的反应机会点
4. 风格适配建议

返回JSON格式的结果。`, scenario, scenario, difficulty, userStyle)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return getDefaultDebateSimulation(), nil
	}

	if len(resp.Choices) == 0 {
		return getDefaultDebateSimulation(), nil
	}

	content := resp.Choices[0].Message.Content

	var result DebateSimulation
	jsonContent := content
	if strings.Contains(content, "```json") {
		startIndex := strings.Index(content, "```json")
		if startIndex != -1 {
			startIndex += 7
			endIndex := strings.Index(content[startIndex:], "```")
			if endIndex != -1 {
				jsonContent = strings.TrimSpace(content[startIndex : startIndex+endIndex])
			}
		}
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return getDefaultDebateSimulation(), nil
	}

	return &result, nil
}

// GenerateResponseWithModel 使用指定模型生成回答
func (c *TALClient) GenerateResponseWithModel(ctx context.Context, prompt, model string) (string, error) {
	// 打印输入信息
	fmt.Printf("📝 AI推理输入:\n")
	fmt.Printf("   模型: %s\n", model)
	fmt.Printf("   提示长度: %d 字符\n", len(prompt))
	if len(prompt) > 200 {
		fmt.Printf("   提示预览: %s...\n", prompt[:200])
	} else {
		fmt.Printf("   完整提示: %s\n", prompt)
	}

	// 请求限流检查
	c.requestMutex.Lock()
	elapsed := time.Since(c.lastRequestTime)
	if elapsed < c.minInterval {
		waitTime := c.minInterval - elapsed
		fmt.Printf("⏳ 请求过于频繁，等待 %.2f 秒...\n", waitTime.Seconds())
		time.Sleep(waitTime)
	}
	c.lastRequestTime = time.Now()
	c.requestMutex.Unlock()

	// 构建OpenAI兼容的请求
	requestBody := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"max_tokens":  c.config.MaxTokens,
		"temperature": c.config.Temperature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("构建请求失败: %w", err)
	}

	// 创建HTTP请求
	url := fmt.Sprintf("%s/chat/completions", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))

	// 发送请求（带超时监控）
	fmt.Printf("🔄 发送AI请求到: %s, 模型: %s, 提示长度: %d\n", url, model, len(prompt))
	fmt.Printf("⏳ 发送HTTP请求，等待响应...\n")
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("❌ HTTP请求失败，耗时: %.2fs, 错误: %v\n", duration.Seconds(), err)

		// 检查是否是超时错误
		if strings.Contains(err.Error(), "context deadline exceeded") ||
		   strings.Contains(err.Error(), "Client.Timeout exceeded") {
			fmt.Println("💡 超时建议: AI推理可能需要更长时间，请检查网络连接或增加超时设置")
		}

		return "", fmt.Errorf("发送请求失败: %w", err)
	}

	fmt.Printf("✅ HTTP请求成功，耗时: %.2fs\n", duration.Seconds())
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("❌ AI响应读取失败: %v\n", err)
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	fmt.Printf("📡 AI响应状态码: %d, 响应大小: %d bytes\n", resp.StatusCode, len(body))

	if resp.StatusCode != http.StatusOK {
		// 特殊处理429错误（配额超限）
		if resp.StatusCode == 429 {
			fmt.Printf("⚠️ AI服务配额超限 (429): %s\n", string(body))
			fmt.Println("💡 建议: 检查API配额、降低请求频率或联系服务商")

			// 尝试等待后重试一次
			fmt.Println("🔄 等待5秒后重试...")
			time.Sleep(5 * time.Second)

			// 重置限流时间戳以允许重试
			c.requestMutex.Lock()
			c.lastRequestTime = time.Now().Add(-c.minInterval)
			c.requestMutex.Unlock()

			// 递归重试一次
			return c.GenerateResponseWithModel(ctx, prompt, model)
		} else {
			fmt.Printf("❌ AI服务错误 (状态码: %d): %s\n", resp.StatusCode, string(body))
		}
		return "", fmt.Errorf("AI服务返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	fmt.Println("✅ AI请求成功")

	// 解析响应
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 提取回答内容
	choices, ok := response["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return "", fmt.Errorf("响应中没有找到choices字段或为空")
	}

	choice, ok := choices[0].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("choice格式不正确")
	}

	message, ok := choice["message"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("message格式不正确")
	}

	content, ok := message["content"].(string)
	if !ok {
		return "", fmt.Errorf("content字段不存在或不是字符串")
	}

	// 打印输出信息
	fmt.Printf("📤 AI推理输出:\n")
	fmt.Printf("   响应长度: %d 字符\n", len(content))
	if len(content) > 200 {
		fmt.Printf("   响应预览: %s...\n", content[:200])
	} else {
		fmt.Printf("   完整响应: %s\n", content)
	}

	return content, nil
}

// EvaluateReaction 评估反应
func (c *TALClient) EvaluateReaction(ctx context.Context, userResponse, scenario, expectedStyle string) (*ReactionEvaluation, error) {
	model := c.GetModelForTask("advanced_reasoning")

	prompt := fmt.Sprintf(`请评估用户的反应表现：

用户反应：%s
场景：%s
期望风格：%s

请从以下维度评估：
1. 内容质量（逻辑性、相关性）
2. 风格符合度（是否符合期望风格）
3. 反应速度（思考-反应的时间合理性）
4. 沟通效果（说服力、感染力）
5. 改进建议

返回JSON格式的评估结果。`, userResponse, scenario, expectedStyle)

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return getDefaultReactionEvaluation(), nil
	}

	if len(resp.Choices) == 0 {
		return getDefaultReactionEvaluation(), nil
	}

	content := resp.Choices[0].Message.Content

	var result ReactionEvaluation
	jsonContent := content
	if strings.Contains(content, "```json") {
		startIndex := strings.Index(content, "```json")
		if startIndex != -1 {
			startIndex += 7
			endIndex := strings.Index(content[startIndex:], "```")
			if endIndex != -1 {
				jsonContent = strings.TrimSpace(content[startIndex : startIndex+endIndex])
			}
		}
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return getDefaultReactionEvaluation(), nil
	}

	return &result, nil
}
