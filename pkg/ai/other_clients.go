package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient OpenAI客户端
type OpenAIClient struct {
	*BaseClient
	config *OpenAIConfig
	client *openai.Client
}

// NewOpenAIClient 创建OpenAI客户端
func NewOpenAIClient(config OpenAIConfig) (*OpenAIClient, error) {
	openaiConfig := openai.DefaultConfig(config.APIKey)
	openaiConfig.BaseURL = config.BaseURL

	client := openai.NewClientWithConfig(openaiConfig)

	return &OpenAIClient{
		BaseClient: &BaseClient{
			provider: ProviderOpenAI,
		},
		config: &config,
		client: client,
	}, nil
}

// OpenAI客户端方法
func (c *OpenAIClient) GetAvailableModels() []string {
	return []string{"gpt-4o", "gpt-4", "gpt-3.5-turbo"}
}

func (c *OpenAIClient) ValidateModel(model string) bool {
	availableModels := c.GetAvailableModels()
	for _, availableModel := range availableModels {
		if availableModel == model {
			return true
		}
	}
	return false
}

func (c *OpenAIClient) AnalyzeImage(ctx context.Context, imageURL, prompt string) (*ImageAnalysisResult, error) {
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

func (c *OpenAIClient) GenerateQuestions(ctx context.Context, contextInfo string, category string) ([]Question, error) {
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

func (c *OpenAIClient) PolishNote(ctx context.Context, rawContent, contextInfo string) (*PolishedNote, error) {
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

func (c *OpenAIClient) TextToSpeech(ctx context.Context, text, voice, language string, speed float64) ([]byte, string, error) {
	// OpenAI TTS 简化实现
	return c.getDefaultAudioData(text), "wav", nil
}

func (c *OpenAIClient) AnalyzeVideo(ctx context.Context, videoData []byte, format, analysisType string, duration float64) (*VideoAnalysis, error) {
	return c.getDefaultVideoAnalysis(), nil
}

func (c *OpenAIClient) GenerateVideo(ctx context.Context, script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	return c.generateMockVideo(script, style, duration, scenes, voice, language)
}

func (c *OpenAIClient) GenerateReactionTemplates(ctx context.Context, scenario, style string) ([]ReactionTemplate, error) {
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

func (c *OpenAIClient) AnalyzeExpressionStyle(ctx context.Context, personName string, sampleText string) (*StyleAnalysis, error) {
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

func (c *OpenAIClient) SimulateDebate(ctx context.Context, scenario string, difficulty int, userStyle string) (*DebateSimulation, error) {
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

func (c *OpenAIClient) EvaluateReaction(ctx context.Context, userResponse, scenario, expectedStyle string) (*ReactionEvaluation, error) {
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

func (c *OpenAIClient) GetModelForTask(task string) string {
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
		return c.config.Models.TextGeneration // 默认使用文本生成模型
	}
}

// 默认方法实现
func (c *OpenAIClient) getDefaultImageAnalysis(imageURL, prompt string) *ImageAnalysisResult {
	return &ImageAnalysisResult{
		ObjectName:     "分析对象",
		Category:       "general",
		Description:    "由于AI服务暂时不可用，这里提供一个模拟的分析结果。在实际环境中，这个结果将由AI模型生成。",
		Confidence:     0.80,
		KeyFeatures:    []string{"特征分析", "形态识别", "内容描述"},
		ScientificName: "未知",
	}
}

func (c *OpenAIClient) getDefaultQuestions(category string) []Question {
	switch category {
	case "debate":
		return []Question{
			{
				Content:    "面对对方质疑时，你会如何回应？",
				Type:       "scenario",
				Difficulty: "basic",
				Purpose:    "练习基础反应能力",
			},
			{
				Content:    "如何在保持立场的同时缓和气氛？",
				Type:       "strategy",
				Difficulty: "intermediate",
				Purpose:    "学习策略性沟通",
			},
			{
				Content:    "面对情绪化的对手，如何控制对话节奏？",
				Type:       "evaluation",
				Difficulty: "advanced",
				Purpose:    "提升高级沟通技巧",
			},
		}
	default:
		return []Question{
			{
				Content:    "在这个沟通场景中，你的第一反应是什么？",
				Type:       "scenario",
				Difficulty: "basic",
				Purpose:    "建立反应意识",
			},
			{
				Content:    "如何用不同的风格表达同样的观点？",
				Type:       "strategy",
				Difficulty: "intermediate",
				Purpose:    "练习风格切换",
			},
			{
				Content:    "面对复杂沟通情境，你的核心策略是什么？",
				Type:       "evaluation",
				Difficulty: "advanced",
				Purpose:    "提升综合沟通能力",
			},
		}
	}
}

func (c *OpenAIClient) getDefaultPolishedNote(rawContent, contextInfo string) *PolishedNote {
	return &PolishedNote{
		Title:         "沟通训练记录",
		Summary:       "这是对用户沟通训练过程的总结记录",
		KeyPoints:     []string{"记录了训练过程", "总结了经验教训"},
		Questions:     []string{"如何改进沟通效果？"},
		Improvements:  []string{"增强表达清晰度", "优化反应速度"},
		FormattedText: rawContent,
	}
}

func (c *OpenAIClient) getDefaultVideoAnalysis() *VideoAnalysis {
	return &VideoAnalysis{
		Summary: &VideoSummary{
			Title:       "视频分析结果",
			Description: "视频内容分析",
			Keywords:    []string{"沟通", "表达", "分析"},
			Category:    "communication",
			Duration:    0,
		},
	}
}

func (c *OpenAIClient) generateMockVideo(script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	// 模拟视频数据
	mockVideoData := []byte("mock-video-data")
	format := "mp4"
	actualDuration := duration
	metadata := &VideoMetadata{
		Title:         "生成的视频",
		Description:   "基于脚本生成的沟通演示视频",
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	}
	return mockVideoData, format, actualDuration, metadata, nil
}

func (c *OpenAIClient) getDefaultReactionTemplates(scenario, style string) []ReactionTemplate {
	return []ReactionTemplate{
		{
			Scenario:    scenario,
			Steps:       []string{"识别情境", "选择风格", "组织回应"},
			KeyPhrases:  []string{"我理解你的观点", "让我来分享我的看法"},
			StyleNotes:  "采用" + style + "风格进行回应",
		},
	}
}

func (c *OpenAIClient) getDefaultStyleAnalysis(personName string) *StyleAnalysis {
	return &StyleAnalysis{
		PersonName: personName,
		LanguageFeatures: map[string]interface{}{
			"clarity": "清晰表达",
			"logic":   "逻辑严谨",
		},
		ThinkingPatterns: map[string]interface{}{
			"structure": "结构化思维",
			"strategy":  "策略性思考",
		},
		CommunicationStrategy: map[string]interface{}{
			"professional": "专业沟通",
			"empathy":      "共情表达",
		},
		PersonalTraits: map[string]interface{}{
			"authority": "权威感",
			"kindness":  "亲和力",
		},
		OverallScore: 8.5,
		StyleTags:    []string{"专业", "严谨", "共情"},
	}
}

func (c *OpenAIClient) getDefaultDebateSimulation(scenario string, difficulty int, userStyle string) *DebateSimulation {
	return &DebateSimulation{
		Scenario:        scenario,
		OpponentOpening: "这是我的观点，你觉得呢？",
		InteractionRounds: []DebateRound{
			{
				RoundNumber:      1,
				OpponentMove:     "我不同意你的观点",
				ExpectedResponse: "让我来解释我的看法",
				ReactionTips:     "保持冷静，清晰表达",
			},
		},
		KeyReactionPoints: []string{"观点冲突点", "共识建立点"},
		StyleSuggestions:  []string{"保持" + userStyle + "风格"},
		Difficulty:        difficulty,
	}
}

func (c *OpenAIClient) getDefaultReactionEvaluation() *ReactionEvaluation {
	return &ReactionEvaluation{
		ContentQuality: EvaluationItem{
			Score:       8.0,
			Description: "内容质量良好",
			Suggestions: []string{"继续保持逻辑清晰"},
		},
		StyleConformity: EvaluationItem{
			Score:       7.5,
			Description: "风格符合度较高",
			Suggestions: []string{"可以进一步加强风格特色"},
		},
		ReactionSpeed: EvaluationItem{
			Score:       6.0,
			Description: "反应速度一般",
			Suggestions: []string{"适当加快反应速度"},
		},
		CommunicationEffect: EvaluationItem{
			Score:       8.5,
			Description: "沟通效果良好",
			Suggestions: []string{"继续保持良好的沟通效果"},
		},
		OverallScore: 7.5,
		Strengths:     []string{"表达清晰", "逻辑性强"},
		Improvements:  []string{"提高反应速度", "增强互动性"},
	}
}

func (c *OpenAIClient) getDefaultAudioData(text string) []byte {
	// 模拟音频数据
	return []byte("mock-audio-data-" + text)
}

// ClaudeClient Claude客户端（占位符实现）
type ClaudeClient struct {
	*BaseClient
	config *ClaudeConfig
}

// NewClaudeClient 创建Claude客户端
func NewClaudeClient(config ClaudeConfig) (*ClaudeClient, error) {
	// TODO: 实现Claude API集成
	return &ClaudeClient{
		BaseClient: &BaseClient{
			provider: ProviderClaude,
		},
		config: &config,
	}, nil
}

// Claude客户端方法（简化为调用默认实现）
func (c *ClaudeClient) GetAvailableModels() []string {
	return []string{"claude-3-opus-20240229", "claude-3-haiku-20240307"}
}

func (c *ClaudeClient) ValidateModel(model string) bool {
	availableModels := c.GetAvailableModels()
	for _, availableModel := range availableModels {
		if availableModel == model {
			return true
		}
	}
	return false
}

func (c *ClaudeClient) AnalyzeImage(ctx context.Context, imageURL, prompt string) (*ImageAnalysisResult, error) {
	return c.getDefaultImageAnalysis(imageURL, prompt), nil
}

func (c *ClaudeClient) GenerateQuestions(ctx context.Context, contextInfo string, category string) ([]Question, error) {
	return c.getDefaultQuestions(category), nil
}

func (c *ClaudeClient) PolishNote(ctx context.Context, rawContent, contextInfo string) (*PolishedNote, error) {
	return c.getDefaultPolishedNote(rawContent, contextInfo), nil
}

func (c *ClaudeClient) TextToSpeech(ctx context.Context, text, voice, language string, speed float64) ([]byte, string, error) {
	return c.getDefaultAudioData(text), "wav", nil
}

func (c *ClaudeClient) AnalyzeVideo(ctx context.Context, videoData []byte, format, analysisType string, duration float64) (*VideoAnalysis, error) {
	return c.getDefaultVideoAnalysis(), nil
}

func (c *ClaudeClient) GenerateVideo(ctx context.Context, script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	return c.generateMockVideo(script, style, duration, scenes, voice, language)
}

func (c *ClaudeClient) GenerateReactionTemplates(ctx context.Context, scenario, style string) ([]ReactionTemplate, error) {
	return c.getDefaultReactionTemplates(scenario, style), nil
}

func (c *ClaudeClient) AnalyzeExpressionStyle(ctx context.Context, personName string, sampleText string) (*StyleAnalysis, error) {
	return c.getDefaultStyleAnalysis(personName), nil
}

func (c *ClaudeClient) SimulateDebate(ctx context.Context, scenario string, difficulty int, userStyle string) (*DebateSimulation, error) {
	return c.getDefaultDebateSimulation(scenario, difficulty, userStyle), nil
}

func (c *ClaudeClient) EvaluateReaction(ctx context.Context, userResponse, scenario, expectedStyle string) (*ReactionEvaluation, error) {
	return c.getDefaultReactionEvaluation(), nil
}

// 默认实现方法
func (c *ClaudeClient) getDefaultImageAnalysis(imageURL, prompt string) *ImageAnalysisResult {
	return &ImageAnalysisResult{
		ObjectName:     "分析对象",
		Category:       "general",
		Description:    "Claude AI分析结果（模拟）",
		Confidence:     0.85,
		KeyFeatures:    []string{"Claude分析特征"},
		ScientificName: "未知",
	}
}

func (c *ClaudeClient) getDefaultQuestions(category string) []Question {
	return []Question{
		{
			Content:    "Claude生成的问题示例",
			Type:       "scenario",
			Difficulty: "basic",
			Purpose:    "Claude AI训练",
		},
	}
}

func (c *ClaudeClient) getDefaultPolishedNote(rawContent, contextInfo string) *PolishedNote {
	return &PolishedNote{
		Title:         "Claude润色记录",
		Summary:       "Claude AI润色结果",
		KeyPoints:     []string{"Claude润色要点"},
		FormattedText: rawContent,
	}
}

func (c *ClaudeClient) getDefaultVideoAnalysis() *VideoAnalysis {
	return &VideoAnalysis{
		Summary: &VideoSummary{
			Title:       "Claude视频分析",
			Description: "Claude AI视频分析结果",
			Keywords:    []string{"Claude", "分析"},
			Category:    "educational",
			Duration:    60.0,
		},
	}
}

func (c *ClaudeClient) generateMockVideo(script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	mockVideoData := c.generateMockMP4Data(script, &VideoMetadata{
		Title:         "Claude生成的演示视频",
		Description:   fmt.Sprintf("Claude基于脚本生成的演示视频: %s", script[:min(100, len(script))]),
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	})

	return mockVideoData, "mp4", duration, &VideoMetadata{
		Title:         "Claude生成的演示视频",
		Description:   fmt.Sprintf("Claude基于脚本生成的演示视频: %s", script[:min(100, len(script))]),
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	}, nil
}

func (c *ClaudeClient) getDefaultReactionTemplates(scenario, style string) []ReactionTemplate {
	return []ReactionTemplate{
		{
			Scenario: "Claude生成的反应场景",
			Steps: []string{
				"Claude分析步骤1",
				"Claude分析步骤2",
			},
			KeyPhrases: []string{
				"Claude关键话术1",
				"Claude关键话术2",
			},
			StyleNotes: "Claude风格要点",
		},
	}
}

func (c *ClaudeClient) getDefaultStyleAnalysis(personName string) *StyleAnalysis {
	return &StyleAnalysis{
		PersonName: personName,
		LanguageFeatures: map[string]interface{}{
			"vocabulary": "Claude分析的词汇特点",
		},
		OverallScore: 8.5,
		StyleTags:    []string{"Claude分析", "AI生成"},
	}
}

func (c *ClaudeClient) getDefaultDebateSimulation(scenario string, difficulty int, userStyle string) *DebateSimulation {
	return &DebateSimulation{
		Scenario:        scenario,
		OpponentOpening: "Claude模拟的对手开场白...",
		KeyReactionPoints: []string{"Claude关键反应点"},
		StyleSuggestions: []string{"Claude风格建议"},
		Difficulty: difficulty,
	}
}

func (c *ClaudeClient) getDefaultReactionEvaluation() *ReactionEvaluation {
	return &ReactionEvaluation{
		OverallScore: 8.0,
		Strengths:    []string{"Claude评估优势"},
		Improvements: []string{"Claude改进建议"},
	}
}

func (c *ClaudeClient) getDefaultAudioData(text string) []byte {
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

func (c *ClaudeClient) generateMockMP4Data(script string, metadata *VideoMetadata) []byte {
	ftypBox := []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70,
		0x69, 0x73, 0x6F, 0x6D, 0x00, 0x00, 0x00, 0x01,
		0x69, 0x73, 0x6F, 0x6D, 0x61, 0x76, 0x63, 0x31,
	}

	scriptData := []byte(fmt.Sprintf("CLAUDE_GENERATED_VIDEO\nTitle: %s\nDescription: %s\nScript: %s\nLanguage: %s\nResolution: %s\n",
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

// AzureClient Azure AI客户端（占位符实现）
type AzureClient struct {
	*BaseClient
	config *AzureConfig
}

// NewAzureClient 创建Azure客户端
func NewAzureClient(config AzureConfig) (*AzureClient, error) {
	// TODO: 实现Azure OpenAI API集成
	return &AzureClient{
		BaseClient: &BaseClient{
			provider: ProviderAzure,
		},
		config: &config,
	}, nil
}

// Azure客户端方法（简化为返回默认值）
func (c *AzureClient) GetAvailableModels() []string {
	return []string{"gpt-4", "gpt-35-turbo"}
}

func (c *AzureClient) ValidateModel(model string) bool {
	availableModels := c.GetAvailableModels()
	for _, availableModel := range availableModels {
		if availableModel == model {
			return true
		}
	}
	return false
}

func (c *AzureClient) AnalyzeImage(ctx context.Context, imageURL, prompt string) (*ImageAnalysisResult, error) {
	return &ImageAnalysisResult{
		ObjectName:     "Azure分析对象",
		Category:       "general",
		Description:    "Azure AI分析结果（模拟）",
		Confidence:     0.82,
		KeyFeatures:    []string{"Azure分析特征"},
		ScientificName: "未知",
	}, nil
}

func (c *AzureClient) GenerateQuestions(ctx context.Context, contextInfo string, category string) ([]Question, error) {
	return []Question{
		{
			Content:    "Azure生成的问题示例",
			Type:       "scenario",
			Difficulty: "basic",
			Purpose:    "Azure AI训练",
		},
	}, nil
}

func (c *AzureClient) PolishNote(ctx context.Context, rawContent, contextInfo string) (*PolishedNote, error) {
	return &PolishedNote{
		Title:         "Azure润色记录",
		Summary:       "Azure AI润色结果",
		KeyPoints:     []string{"Azure润色要点"},
		FormattedText: rawContent,
	}, nil
}

func (c *AzureClient) TextToSpeech(ctx context.Context, text, voice, language string, speed float64) ([]byte, string, error) {
	return c.getDefaultAudioData(), "wav", nil
}

func (c *AzureClient) AnalyzeVideo(ctx context.Context, videoData []byte, format, analysisType string, duration float64) (*VideoAnalysis, error) {
	return &VideoAnalysis{
		Summary: &VideoSummary{
			Title:       "Azure视频分析",
			Description: "Azure AI视频分析结果",
			Keywords:    []string{"Azure", "分析"},
			Category:    "educational",
			Duration:    60.0,
		},
	}, nil
}

func (c *AzureClient) GenerateVideo(ctx context.Context, script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	return c.generateMockVideo(script, style, duration, scenes, voice, language)
}

func (c *AzureClient) GenerateReactionTemplates(ctx context.Context, scenario, style string) ([]ReactionTemplate, error) {
	return []ReactionTemplate{
		{
			Scenario: "Azure生成的反应场景",
			Steps: []string{
				"Azure分析步骤1",
				"Azure分析步骤2",
			},
			KeyPhrases: []string{
				"Azure关键话术1",
				"Azure关键话术2",
			},
			StyleNotes: "Azure风格要点",
		},
	}, nil
}

func (c *AzureClient) AnalyzeExpressionStyle(ctx context.Context, personName string, sampleText string) (*StyleAnalysis, error) {
	return &StyleAnalysis{
		PersonName: personName,
		LanguageFeatures: map[string]interface{}{
			"vocabulary": "Azure分析的词汇特点",
		},
		OverallScore: 8.2,
		StyleTags:    []string{"Azure分析", "AI生成"},
	}, nil
}

func (c *AzureClient) SimulateDebate(ctx context.Context, scenario string, difficulty int, userStyle string) (*DebateSimulation, error) {
	return &DebateSimulation{
		Scenario:        scenario,
		OpponentOpening: "Azure模拟的对手开场白...",
		KeyReactionPoints: []string{"Azure关键反应点"},
		StyleSuggestions: []string{"Azure风格建议"},
		Difficulty: difficulty,
	}, nil
}

func (c *AzureClient) EvaluateReaction(ctx context.Context, userResponse, scenario, expectedStyle string) (*ReactionEvaluation, error) {
	return &ReactionEvaluation{
		OverallScore: 7.8,
		Strengths:    []string{"Azure评估优势"},
		Improvements: []string{"Azure改进建议"},
	}, nil
}

func (c *AzureClient) getDefaultAudioData() []byte {
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

func (c *AzureClient) generateMockVideo(script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	mockVideoData := c.generateMockMP4Data(script, &VideoMetadata{
		Title:         "Azure生成的演示视频",
		Description:   fmt.Sprintf("Azure基于脚本生成的演示视频: %s", script[:min(100, len(script))]),
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	})

	return mockVideoData, "mp4", duration, &VideoMetadata{
		Title:         "Azure生成的演示视频",
		Description:   fmt.Sprintf("Azure基于脚本生成的演示视频: %s", script[:min(100, len(script))]),
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	}, nil
}

func (c *AzureClient) generateMockMP4Data(script string, metadata *VideoMetadata) []byte {
	ftypBox := []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70,
		0x69, 0x73, 0x6F, 0x6D, 0x00, 0x00, 0x00, 0x01,
		0x69, 0x73, 0x6F, 0x6D, 0x61, 0x76, 0x63, 0x31,
	}

	scriptData := []byte(fmt.Sprintf("AZURE_GENERATED_VIDEO\nTitle: %s\nDescription: %s\nScript: %s\nLanguage: %s\nResolution: %s\n",
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

// BaiduClient 百度AI客户端（占位符实现）
type BaiduClient struct {
	*BaseClient
	config *BaiduConfig
}

// NewBaiduClient 创建百度客户端
func NewBaiduClient(config BaiduConfig) (*BaiduClient, error) {
	// TODO: 实现百度AI API集成
	return &BaiduClient{
		BaseClient: &BaseClient{
			provider: ProviderBaidu,
		},
		config: &config,
	}, nil
}

// Baidu客户端方法（简化为返回默认值）
func (c *BaiduClient) GetAvailableModels() []string {
	return []string{"ernie-4.0", "ernie-3.5"}
}

func (c *BaiduClient) ValidateModel(model string) bool {
	availableModels := c.GetAvailableModels()
	for _, availableModel := range availableModels {
		if availableModel == model {
			return true
		}
	}
	return false
}

func (c *BaiduClient) AnalyzeImage(ctx context.Context, imageURL, prompt string) (*ImageAnalysisResult, error) {
	return &ImageAnalysisResult{
		ObjectName:     "百度分析对象",
		Category:       "general",
		Description:    "百度AI分析结果（模拟）",
		Confidence:     0.83,
		KeyFeatures:    []string{"百度分析特征"},
		ScientificName: "未知",
	}, nil
}

func (c *BaiduClient) GenerateQuestions(ctx context.Context, contextInfo string, category string) ([]Question, error) {
	return []Question{
		{
			Content:    "百度生成的问题示例",
			Type:       "scenario",
			Difficulty: "basic",
			Purpose:    "百度AI训练",
		},
	}, nil
}

func (c *BaiduClient) PolishNote(ctx context.Context, rawContent, contextInfo string) (*PolishedNote, error) {
	return &PolishedNote{
		Title:         "百度润色记录",
		Summary:       "百度AI润色结果",
		KeyPoints:     []string{"百度润色要点"},
		FormattedText: rawContent,
	}, nil
}

func (c *BaiduClient) TextToSpeech(ctx context.Context, text, voice, language string, speed float64) ([]byte, string, error) {
	return c.getDefaultAudioData(), "wav", nil
}

func (c *BaiduClient) AnalyzeVideo(ctx context.Context, videoData []byte, format, analysisType string, duration float64) (*VideoAnalysis, error) {
	return &VideoAnalysis{
		Summary: &VideoSummary{
			Title:       "百度视频分析",
			Description: "百度AI视频分析结果",
			Keywords:    []string{"百度", "分析"},
			Category:    "educational",
			Duration:    60.0,
		},
	}, nil
}

func (c *BaiduClient) GenerateVideo(ctx context.Context, script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	return c.generateMockVideo(script, style, duration, scenes, voice, language)
}

func (c *BaiduClient) GenerateReactionTemplates(ctx context.Context, scenario, style string) ([]ReactionTemplate, error) {
	return []ReactionTemplate{
		{
			Scenario: "百度生成的反应场景",
			Steps: []string{
				"百度分析步骤1",
				"百度分析步骤2",
			},
			KeyPhrases: []string{
				"百度关键话术1",
				"百度关键话术2",
			},
			StyleNotes: "百度风格要点",
		},
	}, nil
}

func (c *BaiduClient) AnalyzeExpressionStyle(ctx context.Context, personName string, sampleText string) (*StyleAnalysis, error) {
	return &StyleAnalysis{
		PersonName: personName,
		LanguageFeatures: map[string]interface{}{
			"vocabulary": "百度分析的词汇特点",
		},
		OverallScore: 8.3,
		StyleTags:    []string{"百度分析", "AI生成"},
	}, nil
}

func (c *BaiduClient) SimulateDebate(ctx context.Context, scenario string, difficulty int, userStyle string) (*DebateSimulation, error) {
	return &DebateSimulation{
		Scenario:        scenario,
		OpponentOpening: "百度模拟的对手开场白...",
		KeyReactionPoints: []string{"百度关键反应点"},
		StyleSuggestions: []string{"百度风格建议"},
		Difficulty: difficulty,
	}, nil
}

func (c *BaiduClient) EvaluateReaction(ctx context.Context, userResponse, scenario, expectedStyle string) (*ReactionEvaluation, error) {
	return &ReactionEvaluation{
		OverallScore: 8.1,
		Strengths:    []string{"百度评估优势"},
		Improvements: []string{"百度改进建议"},
	}, nil
}

func (c *BaiduClient) getDefaultAudioData() []byte {
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

func (c *BaiduClient) generateMockVideo(script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	mockVideoData := c.generateMockMP4Data(script, &VideoMetadata{
		Title:         "百度生成的演示视频",
		Description:   fmt.Sprintf("百度基于脚本生成的演示视频: %s", script[:min(100, len(script))]),
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	})

	return mockVideoData, "mp4", duration, &VideoMetadata{
		Title:         "百度生成的演示视频",
		Description:   fmt.Sprintf("百度基于脚本生成的演示视频: %s", script[:min(100, len(script))]),
		Scenes:        scenes,
		AudioLanguage: language,
		Resolution:    "1920x1080",
	}, nil
}

func (c *BaiduClient) generateMockMP4Data(script string, metadata *VideoMetadata) []byte {
	ftypBox := []byte{
		0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70,
		0x69, 0x73, 0x6F, 0x6D, 0x00, 0x00, 0x00, 0x01,
		0x69, 0x73, 0x6F, 0x6D, 0x61, 0x76, 0x63, 0x31,
	}

	scriptData := []byte(fmt.Sprintf("BAIDU_GENERATED_VIDEO\nTitle: %s\nDescription: %s\nScript: %s\nLanguage: %s\nResolution: %s\n",
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
