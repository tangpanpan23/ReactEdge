package ai

import (
	"context"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient OpenAI客户端
type OpenAIClient struct {
	*BaseClient
	client *openai.Client
	config *OpenAIConfig
}

// NewOpenAIClient 创建OpenAI客户端
func NewOpenAIClient(config OpenAIConfig) (*OpenAIClient, error) {
	clientConfig := openai.DefaultConfig(config.APIKey)
	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	}

	if config.Timeout > 0 {
		clientConfig.HTTPClient.Timeout = time.Duration(config.Timeout) * time.Second
	} else {
		clientConfig.HTTPClient.Timeout = 30 * time.Second
	}

	client := &OpenAIClient{
		BaseClient: &BaseClient{
			provider: ProviderOpenAI,
		},
		client: openai.NewClientWithConfig(clientConfig),
		config: &config,
	}

	return client, nil
}

// GetAvailableModels 获取可用模型列表
func (c *OpenAIClient) GetAvailableModels() []string {
	return []string{
		"gpt-4o",       // 最新的GPT-4优化版
		"gpt-4",        // GPT-4
		"gpt-4-turbo",  // GPT-4 Turbo
		"gpt-3.5-turbo", // GPT-3.5 Turbo
	}
}

// ValidateModel 验证模型是否可用
func (c *OpenAIClient) ValidateModel(model string) bool {
	availableModels := c.GetAvailableModels()
	for _, availableModel := range availableModels {
		if availableModel == model {
			return true
		}
	}
	return false
}

// AnalyzeImage 图像分析
func (c *OpenAIClient) AnalyzeImage(ctx context.Context, imageURL, prompt string) (*ImageAnalysisResult, error) {
	model := c.config.Models.ImageAnalysis
	if model == "" {
		model = "gpt-4o" // 默认使用GPT-4o进行图像分析
	}

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
		return c.getDefaultImageAnalysis(imageURL, prompt), nil
	}

	if len(resp.Choices) == 0 {
		return c.getDefaultImageAnalysis(imageURL, prompt), nil
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
func (c *OpenAIClient) GenerateQuestions(ctx context.Context, contextInfo string, category string) ([]Question, error) {
	model := c.config.Models.TextGeneration
	if model == "" {
		model = "gpt-4"
	}

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
		return c.getDefaultQuestions(category), nil
	}

	if len(resp.Choices) == 0 {
		return c.getDefaultQuestions(category), nil
	}

	content := resp.Choices[0].Message.Content
	questions := parseQuestionsFromJSON(content)
	if len(questions) == 0 {
		return c.getDefaultQuestions(category), nil
	}

	return questions, nil
}

// PolishNote 润色笔记
func (c *OpenAIClient) PolishNote(ctx context.Context, rawContent, contextInfo string) (*PolishedNote, error) {
	model := c.config.Models.TextGeneration
	if model == "" {
		model = "gpt-4"
	}

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
		return c.getDefaultPolishedNote(rawContent, contextInfo), nil
	}

	if len(resp.Choices) == 0 {
		return c.getDefaultPolishedNote(rawContent, contextInfo), nil
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
func (c *OpenAIClient) TextToSpeech(ctx context.Context, text, voice, language string, speed float64) ([]byte, string, error) {
	// OpenAI 目前不支持TTS，这里返回默认音频数据
	return c.getDefaultAudioData(text), "wav", nil
}

// AnalyzeVideo 视频分析
func (c *OpenAIClient) AnalyzeVideo(ctx context.Context, videoData []byte, format, analysisType string, duration float64) (*VideoAnalysis, error) {
	// 简化为返回默认分析结果
	return c.getDefaultVideoAnalysis(), nil
}

// GenerateVideo 视频生成
func (c *OpenAIClient) GenerateVideo(ctx context.Context, script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	// OpenAI 目前不支持视频生成，返回模拟数据
	return c.generateMockVideo(script, style, duration, scenes, voice, language)
}

// GenerateReactionTemplates 生成反应模板
func (c *OpenAIClient) GenerateReactionTemplates(ctx context.Context, scenario, style string) ([]ReactionTemplate, error) {
	model := c.config.Models.TextGeneration
	if model == "" {
		model = "gpt-4"
	}

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
		return c.getDefaultReactionTemplates(scenario, style), nil
	}

	if len(resp.Choices) == 0 {
		return c.getDefaultReactionTemplates(scenario, style), nil
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
		return c.getDefaultReactionTemplates(scenario, style), nil
	}

	return result.Templates, nil
}

// AnalyzeExpressionStyle 分析表达风格
func (c *OpenAIClient) AnalyzeExpressionStyle(ctx context.Context, personName string, sampleText string) (*StyleAnalysis, error) {
	model := c.config.Models.AdvancedReasoning
	if model == "" {
		model = "gpt-4"
	}

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
		return c.getDefaultStyleAnalysis(personName), nil
	}

	if len(resp.Choices) == 0 {
		return c.getDefaultStyleAnalysis(personName), nil
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
		return c.getDefaultStyleAnalysis(personName), nil
	}

	return &result, nil
}

// SimulateDebate 模拟辩论
func (c *OpenAIClient) SimulateDebate(ctx context.Context, scenario string, difficulty int, userStyle string) (*DebateSimulation, error) {
	model := c.config.Models.AdvancedReasoning
	if model == "" {
		model = "gpt-4"
	}

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
		return c.getDefaultDebateSimulation(scenario, difficulty, userStyle), nil
	}

	if len(resp.Choices) == 0 {
		return c.getDefaultDebateSimulation(scenario, difficulty, userStyle), nil
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
		return c.getDefaultDebateSimulation(scenario, difficulty, userStyle), nil
	}

	return &result, nil
}

// EvaluateReaction 评估反应
func (c *OpenAIClient) EvaluateReaction(ctx context.Context, userResponse, scenario, expectedStyle string) (*ReactionEvaluation, error) {
	model := c.config.Models.AdvancedReasoning
	if model == "" {
		model = "gpt-4"
	}

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
		return c.getDefaultReactionEvaluation(), nil
	}

	if len(resp.Choices) == 0 {
		return c.getDefaultReactionEvaluation(), nil
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
		return c.getDefaultReactionEvaluation(), nil
	}

	return &result, nil
}

// 默认实现方法（复用TALClient的实现）
func (c *OpenAIClient) getDefaultImageAnalysis(imageURL, prompt string) *ImageAnalysisResult {
	return &ImageAnalysisResult{
		ObjectName:     "分析对象",
		Category:       "general",
		Description:    "由于AI服务暂时不可用，这里提供一个模拟的分析结果。",
		Confidence:     0.80,
		KeyFeatures:    []string{"特征分析", "形态识别", "内容描述"},
		ScientificName: "未知",
	}
}

func (c *OpenAIClient) getDefaultQuestions(category string) []Question {
	return []Question{
		{
			Content:    "在这个沟通场景中，你的第一反应是什么？",
			Type:       "scenario",
			Difficulty: "basic",
			Purpose:    "建立反应意识",
		},
	}
}

func (c *OpenAIClient) getDefaultPolishedNote(rawContent, contextInfo string) *PolishedNote {
	return &PolishedNote{
		Title:         "反应训练记录",
		Summary:       "这是反应训练的记录总结",
		KeyPoints:     []string{"记录了训练过程", "总结了经验教训"},
		FormattedText: rawContent,
	}
}

func (c *OpenAIClient) getDefaultVideoAnalysis() *VideoAnalysis {
	return &VideoAnalysis{
		Summary: &VideoSummary{
			Title:       "训练视频",
			Description: "职场沟通反应训练演示",
			Keywords:    []string{"沟通", "训练", "反应"},
			Category:    "educational",
			Duration:    60.0,
		},
	}
}

func (c *OpenAIClient) generateMockVideo(script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	mockVideoData := c.generateMockMP4Data(script, &VideoMetadata{
		Title:         "AI生成的演示视频",
		Description:   fmt.Sprintf("基于脚本生成的演示视频: %s", script[:min(100, len(script))]),
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	})

	return mockVideoData, "mp4", duration, &VideoMetadata{
		Title:         "AI生成的演示视频",
		Description:   fmt.Sprintf("基于脚本生成的演示视频: %s", script[:min(100, len(script))]),
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	}, nil
}

func (c *OpenAIClient) getDefaultReactionTemplates(scenario, style string) []ReactionTemplate {
	return []ReactionTemplate{
		{
			Scenario: "领导质疑项目进展",
			Steps: []string{
				"先倾听完整质疑内容",
				"用数据回应具体问题",
				"说明改进措施",
			},
			KeyPhrases: []string{
				"您提到的这个问题确实重要",
				"根据我们的数据统计...",
				"我们已经制定了以下改进方案",
			},
			StyleNotes: "保持专业、数据驱动的沟通风格",
		},
	}
}

func (c *OpenAIClient) getDefaultStyleAnalysis(personName string) *StyleAnalysis {
	return &StyleAnalysis{
		PersonName: personName,
		LanguageFeatures: map[string]interface{}{
			"vocabulary": "专业术语丰富",
		},
		OverallScore: 8.0,
		StyleTags:    []string{"专业", "严谨"},
	}
}

func (c *OpenAIClient) getDefaultDebateSimulation(scenario string, difficulty int, userStyle string) *DebateSimulation {
	return &DebateSimulation{
		Scenario:        scenario,
		OpponentOpening: "我认为这个方案有很大问题...",
		KeyReactionPoints: []string{"数据支撑", "逻辑推理"},
		StyleSuggestions: []string{"保持冷静", "用事实说话"},
		Difficulty: difficulty,
	}
}

func (c *OpenAIClient) getDefaultReactionEvaluation() *ReactionEvaluation {
	return &ReactionEvaluation{
		OverallScore: 7.5,
		Strengths:    []string{"逻辑清晰", "表达专业"},
		Improvements: []string{"增加互动性", "注意语速控制"},
	}
}

func (c *OpenAIClient) getDefaultAudioData(text string) []byte {
	wavHeader := []byte{
		0x52, 0x49, 0x46, 0x46, 0x24, 0x08, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45,
		0x66, 0x6D, 0x74, 0x20, 0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
		0x80, 0x3E, 0x00, 0x00, 0x80, 0x3E, 0x00, 0x00, 0x02, 0x00, 0x10, 0x00,
		0x64, 0x61, 0x74, 0x61, 0x00, 0x08, 0x00, 0x00,
	}

	audioData := make([]byte, 2048)
	for i := range audioData {
		audioData[i] = byte((i * 37) % 256)
	}

	return append(wavHeader, audioData...)
}

func (c *OpenAIClient) generateMockMP4Data(script string, metadata *VideoMetadata) []byte {
	ftypBox := []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70,
		0x69, 0x73, 0x6F, 0x6D, 0x00, 0x00, 0x00, 0x01,
		0x69, 0x73, 0x6F, 0x6D, 0x61, 0x76, 0x63, 0x31,
	}

	scriptData := []byte(fmt.Sprintf("AI_GENERATED_VIDEO\nTitle: %s\nDescription: %s\nScript: %s\nLanguage: %s\nResolution: %s\n",
		metadata.Title, metadata.Description, script, metadata.AudioLanguage, metadata.Resolution))

	mdatSize := uint32(8 + len(scriptData))
	mdatBox := make([]byte, 8+len(scriptData))
	mdatBox[0] = byte(mdatSize >> 24)
	mdatBox[1] = byte(mdatSize >> 16)
	mdatBox[2] = byte(mdatSize >> 8)
	mdatBox[3] = byte(mdatSize)
	mdatBox[4] = 'm'
	mdatBox[5] = 'd'
	mdatBox[6] = 'a'
	mdatBox[7] = 't'
	copy(mdatBox[8:], scriptData)

	moovBox := []byte{
		0x00, 0x00, 0x00, 0x6C, 0x6D, 0x6F, 0x6F, 0x76,
	}

	videoData := append(ftypBox, mdatBox...)
	videoData = append(videoData, moovBox...)

	if len(videoData) < 1024 {
		padding := make([]byte, 1024-len(videoData))
		videoData = append(videoData, padding...)
	}

	return videoData
}
