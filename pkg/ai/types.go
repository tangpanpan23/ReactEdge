package ai

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ImageAnalysisResult 图像分析结果
type ImageAnalysisResult struct {
	ObjectName     string   `json:"object_name"`
	Category       string   `json:"category"`
	Confidence     float64  `json:"confidence"`
	Description    string   `json:"description"`
	KeyFeatures    []string `json:"key_features"`
	ScientificName string   `json:"scientific_name"`
}

// Question 问题结构
type Question struct {
	Content    string `json:"content"`
	Type       string `json:"type"`
	Difficulty string `json:"difficulty"`
	Purpose    string `json:"purpose"`
}

// PolishedNote 润色后的笔记
type PolishedNote struct {
	Title             string   `json:"title"`
	Summary           string   `json:"summary"`
	KeyPoints         []string `json:"key_points"`
	ScientificConcepts []string `json:"scientific_concepts,omitempty"`
	CommunicationTips []string `json:"communication_tips,omitempty"`
	Questions         []string `json:"questions"`
	Improvements      []string `json:"improvements,omitempty"`
	Connections       []string `json:"connections,omitempty"`
	FormattedText     string   `json:"formatted_text"`
}

// VideoAnalysis 视频分析结果
type VideoAnalysis struct {
	Scenes  []*SceneAnalysis  `json:"scenes"`
	Objects []*ObjectDetection `json:"objects"`
	Emotions []*EmotionAnalysis `json:"emotions"`
	Texts   []*TextRecognition `json:"texts"`
	Audio   []*AudioAnalysis  `json:"audio"`
	Summary *VideoSummary     `json:"summary"`
}

// SceneAnalysis 场景分析
type SceneAnalysis struct {
	Timestamp   float64 `json:"timestamp"`
	SceneType   string  `json:"scene_type"`
	Description string  `json:"description"`
	Confidence  float64 `json:"confidence"`
}

// ObjectDetection 物体检测
type ObjectDetection struct {
	Timestamp  float64       `json:"timestamp"`
	ObjectName string        `json:"object_name"`
	Confidence float64       `json:"confidence"`
	Bbox       *BoundingBox `json:"bbox"`
}

// EmotionAnalysis 情感分析
type EmotionAnalysis struct {
	Timestamp  float64 `json:"timestamp"`
	Emotion    string  `json:"emotion"`
	Confidence float64 `json:"confidence"`
}

// TextRecognition 文字识别
type TextRecognition struct {
	Timestamp float64       `json:"timestamp"`
	Text      string        `json:"text"`
	Language  string        `json:"language"`
	Confidence float64      `json:"confidence"`
	Bbox      *BoundingBox `json:"bbox"`
}

// AudioAnalysis 音频分析
type AudioAnalysis struct {
	Timestamp    float64 `json:"timestamp"`
	Transcription string `json:"transcription"`
	Language     string  `json:"language"`
	Confidence   float64 `json:"confidence"`
}

// VideoSummary 视频总结
type VideoSummary struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Keywords    []string `json:"keywords"`
	Category    string   `json:"category"`
	Duration    float64  `json:"duration"`
}

// VideoMetadata 视频元数据
type VideoMetadata struct {
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Scenes        []string `json:"scenes"`
	AudioLanguage string   `json:"audio_language"`
	Resolution    string   `json:"resolution"`
}

// BoundingBox 边界框
type BoundingBox struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// ReactionTemplate 反应模板
type ReactionTemplate struct {
	Scenario    string   `json:"scenario"`
	Steps       []string `json:"steps"`
	KeyPhrases  []string `json:"key_phrases"`
	StyleNotes  string   `json:"style_notes"`
}

// StyleAnalysis 风格分析
type StyleAnalysis struct {
	PersonName       string            `json:"person_name"`
	LanguageFeatures map[string]interface{} `json:"language_features"`
	ThinkingPatterns map[string]interface{} `json:"thinking_patterns"`
	CommunicationStrategy map[string]interface{} `json:"communication_strategy"`
	PersonalTraits   map[string]interface{} `json:"personal_traits"`
	OverallScore     float64           `json:"overall_score"`
	StyleTags        []string          `json:"style_tags"`
}

// DebateSimulation 辩论模拟
type DebateSimulation struct {
	Scenario         string   `json:"scenario"`
	OpponentOpening  string   `json:"opponent_opening"`
	InteractionRounds []DebateRound `json:"interaction_rounds"`
	KeyReactionPoints []string `json:"key_reaction_points"`
	StyleSuggestions []string `json:"style_suggestions"`
	Difficulty       int      `json:"difficulty"`
}

// DebateRound 辩论回合
type DebateRound struct {
	RoundNumber   int    `json:"round_number"`
	OpponentMove  string `json:"opponent_move"`
	ExpectedResponse string `json:"expected_response"`
	ReactionTips  string `json:"reaction_tips"`
}

// ReactionEvaluation 反应评估
type ReactionEvaluation struct {
	ContentQuality     EvaluationItem `json:"content_quality"`
	StyleConformity    EvaluationItem `json:"style_conformity"`
	ReactionSpeed      EvaluationItem `json:"reaction_speed"`
	CommunicationEffect EvaluationItem `json:"communication_effect"`
	OverallScore       float64        `json:"overall_score"`
	Strengths          []string       `json:"strengths"`
	Improvements       []string       `json:"improvements"`
}

// EvaluationItem 评估项
type EvaluationItem struct {
	Score       float64 `json:"score"`
	Description string  `json:"description"`
	Suggestions []string `json:"suggestions"`
}

// 默认实现方法

// getDefaultImageAnalysis 默认图像分析结果
func getDefaultImageAnalysis() *ImageAnalysisResult {
	return &ImageAnalysisResult{
		ObjectName:     "分析对象",
		Category:       "general",
		Description:    "由于AI服务暂时不可用，这里提供一个模拟的分析结果。在实际环境中，这个结果将由AI模型生成。",
		Confidence:     0.80,
		KeyFeatures:    []string{"特征分析", "形态识别", "内容描述"},
		ScientificName: "未知",
	}
}

// getDefaultQuestions 默认问题列表
func getDefaultQuestions() []Question {
	// 返回固定的默认问题列表
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
			Content:    "面对复杂情况，如何快速做出有效反应？",
			Type:       "evaluation",
			Difficulty: "advanced",
			Purpose:    "提升综合能力",
		},
	}
}

// getDefaultPolishedNote 默认润色结果
func getDefaultPolishedNote() *PolishedNote {
	return &PolishedNote{
		Title:         "反应训练记录",
		Summary:       "这是反应训练的记录总结",
		KeyPoints:     []string{"记录了训练过程", "总结了经验教训"},
		CommunicationTips: []string{"注意倾听", "清晰表达", "保持礼貌"},
		Questions:     []string{"还有哪些地方可以改进？"},
		Improvements:  []string{"增加练习频率", "尝试不同场景"},
		FormattedText: "这是润色后的内容。由于AI服务暂时不可用，这里提供一个模拟的润色结果。在实际环境中，这个结果将由AI模型生成。",
	}
}

// getDefaultVideoAnalysis 默认视频分析结果
func getDefaultVideoAnalysis() *VideoAnalysis {
	return &VideoAnalysis{
		Scenes: []*SceneAnalysis{
			{
				Timestamp:   0.0,
				SceneType:   "training",
				Description: "职场沟通训练场景",
				Confidence:  0.85,
			},
		},
		Objects: []*ObjectDetection{
			{
				Timestamp:  5.0,
				ObjectName: "演示对象",
				Confidence: 0.82,
			},
		},
		Summary: &VideoSummary{
			Title:       "训练视频",
			Description: "职场沟通反应训练演示",
			Keywords:    []string{"沟通", "训练", "反应"},
			Category:    "educational",
			Duration:    60.0,
		},
	}
}

// generateMockVideo 生成模拟视频
func (c *TALClient) generateMockVideo(script, style string, duration float64, scenes []string, voice, language string) ([]byte, string, float64, *VideoMetadata, error) {
	// 生成模拟的MP4文件数据
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

// getDefaultReactionTemplates 默认反应模板
func getDefaultReactionTemplates() []ReactionTemplate {
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
		{
			Scenario: "同事提出不同意见",
			Steps: []string{
				"认可对方观点的价值",
				"说明自己的考虑角度",
				"寻求共识或妥协方案",
			},
			KeyPhrases: []string{
				"您的观点很有道理",
				"从另一个角度来看...",
				"我们可以考虑这样的方案",
			},
			StyleNotes: "开放包容、寻求共赢的沟通方式",
		},
	}
}

// getDefaultStyleAnalysis 默认风格分析
func getDefaultStyleAnalysis() *StyleAnalysis {
	return &StyleAnalysis{
		PersonName: "康辉",
		LanguageFeatures: map[string]interface{}{
			"vocabulary": "专业术语丰富",
			"sentence_structure": "逻辑清晰",
		},
		ThinkingPatterns: map[string]interface{}{
			"logic_structure": "层层递进",
			"argumentation": "数据支撑",
		},
		CommunicationStrategy: map[string]interface{}{
			"position_expression": "坚定但不强硬",
			"conflict_handling": "寻求共识",
		},
		PersonalTraits: map[string]interface{}{
			"unique_identifiers": []string{"专业性强", "逻辑严谨"},
			"style_labels": []string{"专业型", "分析型"},
		},
		OverallScore: 8.5,
		StyleTags:    []string{"专业", "严谨", "建设性"},
	}
}

// getDefaultDebateSimulation 默认辩论模拟
func getDefaultDebateSimulation() *DebateSimulation {
	return &DebateSimulation{
		Scenario:        "项目方案讨论",
		OpponentOpening: "我认为这个方案有很大问题...",
		InteractionRounds: []DebateRound{
			{
				RoundNumber: 1,
				OpponentMove: "这个方案成本太高了",
				ExpectedResponse: "让我们从投资回报率的角度来看",
				ReactionTips: "用数据反驳，保持专业",
			},
		},
		KeyReactionPoints: []string{"数据支撑", "逻辑推理", "风格一致"},
		StyleSuggestions: []string{"保持冷静", "用事实说话", "寻求共赢"},
		Difficulty: 2,
	}
}

// getDefaultReactionEvaluation 默认反应评估
func getDefaultReactionEvaluation() *ReactionEvaluation {
	return &ReactionEvaluation{
		ContentQuality: EvaluationItem{
			Score:       7.5,
			Description: "内容逻辑清晰，表达准确",
			Suggestions: []string{"可以增加更多具体数据", "注意表达的简洁性"},
		},
		StyleConformity: EvaluationItem{
			Score:       8.0,
			Description: "基本符合期望风格",
			Suggestions: []string{"可以更自然一些"},
		},
		ReactionSpeed: EvaluationItem{
			Score:       7.0,
			Description: "反应速度适中",
			Suggestions: []string{"可以适当加快反应速度"},
		},
		CommunicationEffect: EvaluationItem{
			Score:       8.5,
			Description: "沟通效果良好",
			Suggestions: []string{"注意倾听对方的反馈"},
		},
		OverallScore: 7.8,
		Strengths:    []string{"逻辑清晰", "表达专业", "数据支撑"},
		Improvements: []string{"增加互动性", "注意语速控制", "多练习不同场景"},
	}
}

// getDefaultAudioData 生成默认音频数据
func getDefaultAudioData() []byte {
	// 生成一个简单的WAV文件头部 + 模拟音频数据
	wavHeader := []byte{
		0x52, 0x49, 0x46, 0x46, // "RIFF"
		0x24, 0x08, 0x00, 0x00, // 文件大小
		0x57, 0x41, 0x56, 0x45, // "WAVE"
		0x66, 0x6D, 0x74, 0x20, // "fmt "
		0x10, 0x00, 0x00, 0x00, // fmt chunk大小
		0x01, 0x00,             // 格式：PCM
		0x01, 0x00,             // 声道数：1
		0x80, 0x3E, 0x00, 0x00, // 采样率：16000
		0x80, 0x3E, 0x00, 0x00, // 字节率
		0x02, 0x00,             // 块对齐
		0x10, 0x00,             // 位深度：16
		0x64, 0x61, 0x74, 0x61, // "data"
		0x00, 0x08, 0x00, 0x00, // 数据大小
	}

	// 生成一些模拟的音频数据
	audioData := make([]byte, 2048)
	for i := range audioData {
		audioData[i] = byte((i * 37) % 256) // 简单的伪随机数据
	}

	return append(wavHeader, audioData...)
}

// generateMockMP4Data 生成模拟的MP4文件数据
func (c *TALClient) generateMockMP4Data(script string, metadata *VideoMetadata) []byte {
	// MP4文件的基本结构
	ftypBox := []byte{
		0x00, 0x00, 0x00, 0x20, // box size (32 bytes)
		0x66, 0x74, 0x79, 0x70, // "ftyp"
		0x69, 0x73, 0x6F, 0x6D, // major_brand: isom
		0x00, 0x00, 0x00, 0x01, // minor_version
		0x69, 0x73, 0x6F, 0x6D, // compatible_brands[0]: isom
		0x61, 0x76, 0x63, 0x31, // compatible_brands[1]: avc1
	}

	// 创建包含脚本信息的"视频"数据
	scriptData := []byte(fmt.Sprintf("AI_GENERATED_VIDEO\nTitle: %s\nDescription: %s\nScript: %s\nLanguage: %s\nResolution: %s\n",
		metadata.Title, metadata.Description, script, metadata.AudioLanguage, metadata.Resolution))

	// mdat box (媒体数据)
	mdatSize := uint32(8 + len(scriptData))
	mdatBox := make([]byte, 8+len(scriptData))
	// 大端字节序写入size
	mdatBox[0] = byte(mdatSize >> 24)
	mdatBox[1] = byte(mdatSize >> 16)
	mdatBox[2] = byte(mdatSize >> 8)
	mdatBox[3] = byte(mdatSize)
	// type
	mdatBox[4] = 'm'
	mdatBox[5] = 'd'
	mdatBox[6] = 'a'
	mdatBox[7] = 't'
	// data
	copy(mdatBox[8:], scriptData)

	// 简化的moov box
	moovBox := []byte{
		0x00, 0x00, 0x00, 0x6C, // box size (108 bytes)
		0x6D, 0x6F, 0x6F, 0x76, // "moov"
	}

	// 组合成完整的MP4文件
	videoData := append(ftypBox, mdatBox...)
	videoData = append(videoData, moovBox...)

	// 确保文件大小合理
	if len(videoData) < 1024 {
		padding := make([]byte, 1024-len(videoData))
		videoData = append(videoData, padding...)
	}

	return videoData
}

// 辅助函数

func extractObjectName(content string) string {
	if len(content) > 50 {
		return content[:50] + "..."
	}
	return content
}

func extractCategory(content string) string {
	if strings.Contains(content, "领导") || strings.Contains(content, "汇报") {
		return "workplace"
	}
	return "general"
}

func extractKeyFeatures(content string) []string {
	return []string{"特征1", "特征2"}
}

func extractScientificName(content string) string {
	return "未知"
}

func parseQuestionsFromJSON(content string) []Question {
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

	// 尝试解析JSON
	var result struct {
		Questions []Question `json:"questions"`
	}

	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return []Question{}
	}

	return result.Questions
}

func extractSummaryFromText(text string) string {
	if len(text) > 100 {
		return text[:100] + "..."
	}
	return text
}

func extractKeyPointsFromText(text string) []string {
	points := strings.Split(text, "。")
	var keyPoints []string
	for _, point := range points {
		point = strings.TrimSpace(point)
		if point != "" && len(point) > 5 {
			keyPoints = append(keyPoints, point)
			if len(keyPoints) >= 3 {
				break
			}
		}
	}

	if len(keyPoints) == 0 {
		keyPoints = []string{"学习了新的知识", "发现了有趣的现象"}
	}

	return keyPoints
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// getDefaultVideoData 默认视频数据
func getDefaultVideoData() []byte {
	// 返回一个小的模拟视频数据
	return []byte{0x00, 0x00, 0x00, 0x20, 0x66, 0x74, 0x79, 0x70, 0x6D, 0x70, 0x34, 0x32, 0x00, 0x00, 0x00, 0x00}
}

// getDefaultVideoMetadata 默认视频元数据
func getDefaultVideoMetadata() *VideoMetadata {
	return &VideoMetadata{
		Title:       "模拟视频",
		Description: "由于AI服务暂时不可用，这里提供一个模拟的视频结果。在实际环境中，这个结果将由AI模型生成。",
		Scenes:      []string{"场景1", "场景2"},
		Resolution:  "1920x1080",
	}
}
