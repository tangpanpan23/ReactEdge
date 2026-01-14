package challenge

import (
	"fmt"
	"time"

	"reactedge/internal/ai"
)

// ChallengePhase 挑战阶段
type ChallengePhase int

const (
	PhaseWelcome ChallengePhase = iota
	PhaseAIDeconstruction
	PhasePersonalizedTemplate
	PhaseRecording
	PhaseDNAAnalysis
	PhaseComplete
)

// ChallengeState 挑战状态
type ChallengeState struct {
	CurrentPhase    ChallengePhase     `json:"current_phase"`
	StartTime       time.Time         `json:"start_time"`
	PhaseStartTime  time.Time         `json:"phase_start_time"`
	UserProfile     ai.UserProfile    `json:"user_profile"`
	CurrentTopic    string            `json:"current_topic"`
	UserSpeech      string            `json:"user_speech"`
	ExpressionDNA   *ai.ExpressionDNA `json:"expression_dna,omitempty"`
	PersonalizedTemplate string       `json:"personalized_template"`
	TimeRemaining   int               `json:"time_remaining"` // 秒
}

// ChallengeManager 挑战管理器
type ChallengeManager struct {
	hanAI     *ai.HanStyleAI
	challenges map[string]*ChallengeState
}

// NewManager 创建挑战管理器
func NewManager(hanAI *ai.HanStyleAI) *ChallengeManager {
	return &ChallengeManager{
		hanAI:     hanAI,
		challenges: make(map[string]*ChallengeState),
	}
}

// StartChallenge 开始新挑战
func (cm *ChallengeManager) StartChallenge(userID string) *ChallengeState {
	now := time.Now()
	state := &ChallengeState{
		CurrentPhase:   PhaseWelcome,
		StartTime:      now,
		PhaseStartTime: now,
		CurrentTopic:   cm.getRandomTopic(),
		TimeRemaining:  180, // 3分钟
	}

	cm.challenges[userID] = state
	return state
}

// GetChallengeState 获取挑战状态
func (cm *ChallengeManager) GetChallengeState(userID string) *ChallengeState {
	return cm.challenges[userID]
}

// AdvancePhase 推进到下一阶段
func (cm *ChallengeManager) AdvancePhase(userID string) *ChallengeState {
	state := cm.challenges[userID]
	if state == nil {
		return nil
	}

	now := time.Now()
	state.PhaseStartTime = now

	switch state.CurrentPhase {
	case PhaseWelcome:
		state.CurrentPhase = PhaseAIDeconstruction
		state.TimeRemaining = 180 - int(now.Sub(state.StartTime).Seconds())
	case PhaseAIDeconstruction:
		state.CurrentPhase = PhasePersonalizedTemplate
		// 生成个性化模板
		if state.UserProfile.PrimaryInterest == "" {
			// 如果还没有用户画像，使用默认值
			state.UserProfile = ai.UserProfile{PrimaryInterest: "游戏"}
		}
		state.PersonalizedTemplate = cm.hanAI.GeneratePersonalizedTemplate(state.UserProfile, state.CurrentTopic)
		state.TimeRemaining = 180 - int(now.Sub(state.StartTime).Seconds())
	case PhasePersonalizedTemplate:
		state.CurrentPhase = PhaseRecording
		state.TimeRemaining = 180 - int(now.Sub(state.StartTime).Seconds())
	case PhaseRecording:
		state.CurrentPhase = PhaseDNAAnalysis
		// 分析表达DNA
		if state.UserSpeech != "" {
			state.ExpressionDNA = &ai.ExpressionDNA{}
			*state.ExpressionDNA = cm.hanAI.AnalyzeExpressionDNA(state.UserSpeech, state.UserProfile)
		}
		state.TimeRemaining = 180 - int(now.Sub(state.StartTime).Seconds())
	case PhaseDNAAnalysis:
		state.CurrentPhase = PhaseComplete
		state.TimeRemaining = 0
	}

	return state
}

// SubmitSpeech 提交用户语音
func (cm *ChallengeManager) SubmitSpeech(userID, speech string) *ChallengeState {
	state := cm.challenges[userID]
	if state == nil {
		return nil
	}

	state.UserSpeech = speech

	// 从语音中探测用户画像
	state.UserProfile = cm.hanAI.DetectUserProfile(speech)

	return state
}

// UpdateProfile 更新用户画像
func (cm *ChallengeManager) UpdateProfile(userID string, profile ai.UserProfile) *ChallengeState {
	state := cm.challenges[userID]
	if state == nil {
		return nil
	}

	state.UserProfile = profile

	// 重新生成个性化模板
	state.PersonalizedTemplate = cm.hanAI.GeneratePersonalizedTemplate(profile, state.CurrentTopic)

	return state
}

// GetPhaseContent 获取当前阶段的内容
func (cm *ChallengeManager) GetPhaseContent(state *ChallengeState) map[string]interface{} {
	content := map[string]interface{}{
		"phase": state.CurrentPhase,
		"time_remaining": state.TimeRemaining,
		"topic": state.CurrentTopic,
	}

	switch state.CurrentPhase {
	case PhaseWelcome:
		content["title"] = "🎤 欢迎来到【酷表达实验室】· 韩寒特训版"
		content["subtitle"] = "🎯 今日挑战：课堂突击提问"
		content["description"] = fmt.Sprintf("📚 场景：语文课上，老师突然点名：%s\n\n⏰ 要求：15秒思考，45秒回答，要有自己的观点\n\n🔄 AI将全程分析你的\"表达DNA\"", state.CurrentTopic)

	case PhaseAIDeconstruction:
		content["title"] = "🧠 AI解构【韩寒表达法】三大武器"
		content["weapons"] = []map[string]interface{}{
			{
				"name": "反常规视角 🌪️",
				"description": "普通人：赞美书店变多 → 文化繁荣\n韩寒式：\"当书店开始比拼装修而不是书目，这和奶茶店比杯子颜值有什么区别？\"",
			},
			{
				"name": "精准文化类比 🎬",
				"description": "把抽象概念变成具体场景：\n\"这就像电影院里全是爆米花味，但没人在意放的是什么电影\"",
			},
			{
				"name": "节奏打断技巧 ⚡",
				"description": "在对方预期处突然转折：\n\"很多人说这是好事...(停顿)但好事有时候是最可怕的陷阱\"",
			},
		}
		content["tools"] = []string{
			"【反问模板】\"难道...就代表...?\"",
			"【类比模板】\"这就像...其实不过是...\"",
			"【转折模板】\"表面上看是...实际上暴露了...\"",
		}

	case PhasePersonalizedTemplate:
		content["title"] = "🤖 AI为你生成【个性化应答模板】"
		content["profile_detection"] = fmt.Sprintf("✅ AI探测到你的偏好：%s", state.UserProfile.PrimaryInterest)
		content["template_title"] = fmt.Sprintf("✅ 为你生成【%s版】应答模板：", cm.getInterestDisplayName(state.UserProfile.PrimaryInterest))
		content["template"] = state.PersonalizedTemplate
		content["framework"] = []string{
			"（1）游戏类比切入 → 吸引同龄人",
			"（2）对比转折 → 展现思辨",
			"（3）现象本质 → 提升深度",
			"（4）犀利反问 → 留下印象",
		}

	case PhaseRecording:
		content["title"] = "🎤 现在请用你的风格回答！"
		content["instruction"] = "⏱️ 15秒思考 → 45秒发言"
		content["tips"] = "（思考时AI显示关键词提示：游戏、质量、虚荣、本质...）"
		content["topic"] = state.CurrentTopic

	case PhaseDNAAnalysis:
		if state.ExpressionDNA != nil {
			content["title"] = "📊 你的【表达DNA分析报告】"
			content["sharpeness_score"] = state.ExpressionDNA.SharpenessScore
			content["personality_tags"] = state.ExpressionDNA.PersonalityTags
			content["unique_patterns"] = state.ExpressionDNA.UniquePatterns
			content["thinking_pattern"] = state.ExpressionDNA.ThinkingPattern
			content["metaphor_style"] = state.ExpressionDNA.MetaphorStyle
			content["recommendations"] = state.ExpressionDNA.Recommendations
			content["next_challenge"] = state.ExpressionDNA.NextChallenge
		}

	case PhaseComplete:
		content["title"] = "🎉 挑战完成！"
		content["message"] = "这不止是一次训练。AI发现了你的独特表达天赋，明天的挑战会围绕这个优势继续设计。"
	}

	return content
}

// getRandomTopic 获取随机话题
func (cm *ChallengeManager) getRandomTopic() string {
	topics := []string{
		"\"你对网红书店遍地开花这种现象，怎么看？\"",
		"\"你觉得现在的短视频平台改变了我们的注意力，怎么评价？\"",
		"\"谈谈你对'内卷'这个词的理解\"",
		"\"你认为人工智能会替代哪些工作？\"",
		"\"现在的校园生活和以前有什么不同？\"",
	}

	return topics[time.Now().UnixNano()%int64(len(topics))]
}

// getInterestDisplayName 获取兴趣显示名称
func (cm *ChallengeManager) getInterestDisplayName(interest string) string {
	names := map[string]string{
		"游戏": "游戏玩家",
		"动漫": "动漫爱好者",
		"体育": "体育迷",
		"科技": "科技达人",
		"文艺": "文艺青年",
	}

	if name, ok := names[interest]; ok {
		return name
	}
	return "玩家"
}

// GetCurrentPhaseDuration 获取当前阶段建议时长（秒）
func (cm *ChallengeManager) GetCurrentPhaseDuration(phase ChallengePhase) int {
	durations := map[ChallengePhase]int{
		PhaseWelcome:             30,  // 0.5分钟
		PhaseAIDeconstruction:    60,  // 1分钟
		PhasePersonalizedTemplate: 42,  // 0.7分钟
		PhaseRecording:           30,  // 0.5分钟
		PhaseDNAAnalysis:         18,  // 0.3分钟
	}

	if duration, ok := durations[phase]; ok {
		return duration
	}
	return 30
}
