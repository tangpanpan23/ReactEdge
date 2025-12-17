package ai

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// ExpressionPattern 表达模式
type ExpressionPattern struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Examples    []string `json:"examples"`
	Template    string   `json:"template"`
	Tags        []string `json:"tags"`
}

// UserProfile 用户画像
type UserProfile struct {
	PrimaryInterest string   `json:"primary_interest"` // 游戏/动漫/体育/科技/文艺
	ThinkingStyle   string   `json:"thinking_style"`   // 归纳/演绎/类比
	MetaphorStyle   string   `json:"metaphor_style"`   // 科技/文艺/生活
	Strengths       []string `json:"strengths"`
}

// ExpressionDNA 表达DNA分析结果
type ExpressionDNA struct {
	SharpenessScore  int                 `json:"sharpeness_score"`  // 犀利指数 0-100
	PersonalityTags  []string            `json:"personality_tags"`  // 个性标签
	UniquePatterns   []string            `json:"unique_patterns"`   // 独特模式
	ThinkingPattern  string              `json:"thinking_pattern"`  // 思维模式
	MetaphorStyle    string              `json:"metaphor_style"`    // 类比风格
	RhythmSignature  string              `json:"rhythm_signature"`  // 节奏特征
	UniquenessScore  int                 `json:"uniqueness_score"`  // 独特性分数
	Recommendations  []string            `json:"recommendations"`   // 优化建议
	NextChallenge    string              `json:"next_challenge"`    // 下次挑战
}

// HanStyleAI 韩寒风格AI引擎
type HanStyleAI struct {
	expressionPatterns []ExpressionPattern
	gameAnalogies      map[string][]string
	hanStyleCorpus     []string
	random             *rand.Rand
}

// NewHanStyleAI 创建韩寒风格AI引擎
func NewHanStyleAI() *HanStyleAI {
	rand.Seed(time.Now().UnixNano())

	ai := &HanStyleAI{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	ai.initializeExpressionPatterns()
	ai.initializeGameAnalogies()
	ai.initializeHanCorpus()

	return ai
}

// initializeExpressionPatterns 初始化表达模式库
func (ai *HanStyleAI) initializeExpressionPatterns() {
	ai.expressionPatterns = []ExpressionPattern{
		{
			Name:        "反常规视角",
			Description: "从意想不到的角度切入问题",
			Examples: []string{
				"当书店开始比拼装修而不是书目，这和奶茶店比杯子颜值有什么区别？",
				"很多人说这是好事...(停顿)但好事有时候是最可怕的陷阱",
			},
			Template: "表面上看是{A}，但实际上{B}，这就像{C}",
			Tags:     []string{"视角转换", "反问", "类比"},
		},
		{
			Name:        "精准文化类比",
			Description: "用具体场景解释抽象概念",
			Examples: []string{
				"这就像电影院里全是爆米花味，但没人在意放的是什么电影",
				"现在的网红就像手机里的APP，更新快，卸载也快",
			},
			Template: "这就像{A}，表面{B}，实际上{C}",
			Tags:     []string{"类比", "文化", "现象解构"},
		},
		{
			Name:        "节奏打断技巧",
			Description: "在预期处突然转折制造冲击",
			Examples: []string{
				"大家都说读书很重要...(突然加速)但读什么书更重要",
				"表面繁荣...(停顿)内里空虚",
			},
			Template: "大家都说{A}...(转折)但{B}才是关键",
			Tags:     []string{"节奏", "转折", "强调"},
		},
		{
			Name:        "本质追问法",
			Description: "不断追问现象背后的本质",
			Examples: []string{
				"网红书店热背后，到底在迎合什么需求？",
				"我们是在热爱阅读，还是在热爱被看见？",
			},
			Template: "表面{D}，但{E}？我们是在{F}，还是{G}？",
			Tags:     []string{"追问", "本质", "反思"},
		},
	}
}

// initializeGameAnalogies 初始化游戏类比库
func (ai *HanStyleAI) initializeGameAnalogies() {
	ai.gameAnalogies = map[string][]string{
		"游戏": {
			"这就像《塞尔达》里到处是神庙但解谜都很简单——数量多了，质量却被稀释了",
			"现在的现象就像游戏里的外挂，短期好用，但破坏了游戏平衡",
			"这和MOBA游戏一样，团队重要性超过个人英雄主义",
			"就像RPG游戏，升级不能只看经验值，更要看技能成长",
		},
		"动漫": {
			"这就像动漫里的热血少年，明明弱小却总想挑战强者",
			"现在的文化现象像动漫里的模板剧情，套路化却总有人买账",
			"就像《死亡笔记》，权力越大，责任越重，却也越危险",
		},
		"体育": {
			"这就像足球比赛，战术重要，但个人灵感决定胜负",
			"现在的竞争像马拉松，不是冲刺就能赢的",
			"就像篮球明星，数据好看但团队配合才出冠军",
		},
		"科技": {
			"这就像智能手机，功能越来越多，但很多人只会刷抖音",
			"现在的创新像AI发展，潜力巨大但伦理风险同样大",
			"就像互联网思维，连接一切却也放大了一切问题",
		},
		"文艺": {
			"这就像诗歌朗诵，形式优美但内容空洞最可怕",
			"现在的创作像绘画临摹，技术娴熟却缺少灵魂",
			"就像音乐节，热闹喧嚣但真正在听歌的人不多",
		},
	}
}

// initializeHanCorpus 初始化韩寒语料库（简化版）
func (ai *HanStyleAI) initializeHanCorpus() {
	ai.hanStyleCorpus = []string{
		"很多人说这是好事，但好事有时候是最可怕的陷阱",
		"表面上看是繁荣，实际上暴露了我们用'打卡'代替'阅读'的虚荣",
		"当书店开始比拼装修而不是书目，这和奶茶店比杯子颜值有什么区别",
		"现在的网红就像手机里的APP，更新快，卸载也快",
		"我们是在热爱阅读，还是在热爱被看见",
		"这就像电影院里全是爆米花味，但没人在意放的是什么电影",
		"大家都说读书很重要，但读什么书更重要",
		"表面繁荣，内里空虚",
	}
}

// GetExpressionPatterns 获取表达模式库
func (ai *HanStyleAI) GetExpressionPatterns() []ExpressionPattern {
	return ai.expressionPatterns
}

// GeneratePersonalizedTemplate 生成个性化模板
func (ai *HanStyleAI) GeneratePersonalizedTemplate(profile UserProfile, topic string) string {
	var template strings.Builder

	// 根据兴趣选择类比
	analogies := ai.gameAnalogies[profile.PrimaryInterest]
	if len(analogies) == 0 {
		analogies = ai.gameAnalogies["游戏"] // 默认游戏类比
	}

	gameAnalogy := analogies[ai.random.Intn(len(analogies))]

	// 构建个性化模板
	template.WriteString(fmt.Sprintf("老师，我觉得%s——\n", gameAnalogy))
	template.WriteString("数量多了，质量却被稀释了。\n")
	template.WriteString("表面上是书店繁荣，实际上暴露了我们用'打卡'代替'阅读'的虚荣。\n")
	template.WriteString("如果书店变成拍照背景板，那和游戏里的贴图BUG有什么区别？")

	return template.String()
}

// AnalyzeExpressionDNA 分析表达DNA
func (ai *HanStyleAI) AnalyzeExpressionDNA(userSpeech string, profile UserProfile) ExpressionDNA {
	// 简单的分析逻辑（实际项目中会更复杂）
	words := strings.Fields(userSpeech)
	wordCount := len(words)

	// 计算犀利指数
	sharpenessScore := 60 + ai.random.Intn(40) // 60-99随机

	// 检测思维模式
	thinkingPattern := ai.detectThinkingPattern(userSpeech)

	// 检测类比风格
	metaphorStyle := ai.detectMetaphorStyle(userSpeech)

	// 检测节奏特征
	rhythmSignature := ai.detectRhythmSignature(userSpeech)

	// 计算独特性
	uniquenessScore := 70 + ai.random.Intn(30)

	// 生成个性标签
	personalityTags := ai.generatePersonalityTags(profile, userSpeech)

	// 生成独特模式
	uniquePatterns := ai.generateUniquePatterns(userSpeech)

	// 生成优化建议
	recommendations := ai.generateRecommendations(userSpeech)

	// 生成下次挑战
	nextChallenge := ai.generateNextChallenge(profile)

	return ExpressionDNA{
		SharpenessScore: sharpenessScore,
		PersonalityTags: personalityTags,
		UniquePatterns:  uniquePatterns,
		ThinkingPattern: thinkingPattern,
		MetaphorStyle:   metaphorStyle,
		RhythmSignature: rhythmSignature,
		UniquenessScore: uniquenessScore,
		Recommendations: recommendations,
		NextChallenge:   nextChallenge,
	}
}

// detectThinkingPattern 检测思维模式
func (ai *HanStyleAI) detectThinkingPattern(speech string) string {
	if strings.Contains(speech, "就像") || strings.Contains(speech, "好像") {
		return "类比思维"
	}
	if strings.Contains(speech, "因为") || strings.Contains(speech, "所以") {
		return "逻辑推理"
	}
	if strings.Contains(speech, "我觉得") || strings.Contains(speech, "我认为") {
		return "主观判断"
	}
	return "现象描述"
}

// detectMetaphorStyle 检测类比风格
func (ai *HanStyleAI) detectMetaphorStyle(speech string) string {
	if strings.Contains(speech, "游戏") || strings.Contains(speech, "塞尔达") {
		return "游戏思维"
	}
	if strings.Contains(speech, "电影") || strings.Contains(speech, "音乐") {
		return "文艺表达"
	}
	if strings.Contains(speech, "手机") || strings.Contains(speech, "互联网") {
		return "科技视角"
	}
	return "生活类比"
}

// detectRhythmSignature 检测节奏特征
func (ai *HanStyleAI) detectRhythmSignature(speech string) string {
	sentences := strings.Split(speech, "。")
	if len(sentences) > 3 {
		return "层次递进"
	}
	if strings.Contains(speech, "但") || strings.Contains(speech, "却") {
		return "转折对比"
	}
	return "直线表达"
}

// generatePersonalityTags 生成个性标签
func (ai *HanStyleAI) generatePersonalityTags(profile UserProfile, speech string) []string {
	tags := []string{}

	if profile.PrimaryInterest == "游戏" {
		tags = append(tags, "游戏思维者")
	}

	if strings.Contains(speech, "本质") || strings.Contains(speech, "真正") {
		tags = append(tags, "本质追问者")
	}

	if strings.Count(speech, "？") > 1 {
		tags = append(tags, "犀利发问者")
	}

	if len(tags) == 0 {
		tags = append(tags, "潜力新星")
	}

	return tags
}

// generateUniquePatterns 生成独特模式
func (ai *HanStyleAI) generateUniquePatterns(speech string) []string {
	patterns := []string{}

	if strings.Contains(speech, "具体案例 → 抽象问题 → 犀利反问") {
		patterns = append(patterns, "三段式思维")
	}

	patterns = append(patterns, "敢于挑战常规")

	return patterns
}

// generateRecommendations 生成优化建议
func (ai *HanStyleAI) generateRecommendations(speech string) []string {
	recommendations := []string{}

	if !strings.Contains(speech, "就像") {
		recommendations = append(recommendations,
			"可以尝试加入具体类比，让观点更生动")
	}

	if strings.Count(speech, "很") > 2 {
		recommendations = append(recommendations,
			"减少'很'字的使用，尝试更精准的形容词")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations,
			"你的表达已经很犀利了，继续保持！")
	}

	return recommendations
}

// generateNextChallenge 生成下次挑战
func (ai *HanStyleAI) generateNextChallenge(profile UserProfile) string {
	challenges := []string{
		"如果用游戏机制解释内卷现象，你会怎么比喻？",
		"谈谈你对'网红经济'的看法，用一个生活场景来类比",
		"如果传统行业要和互联网结合，你认为关键是什么？",
		"你觉得现在的教育像什么？用一个你熟悉的事物来比喻",
	}

	return challenges[ai.random.Intn(len(challenges))]
}

// DetectUserProfile 从用户输入中探测用户画像
func (ai *HanStyleAI) DetectUserProfile(speech string) UserProfile {
	profile := UserProfile{}

	// 检测主要兴趣
	if strings.Contains(speech, "游戏") || strings.Contains(speech, "塞尔达") {
		profile.PrimaryInterest = "游戏"
	} else if strings.Contains(speech, "动漫") || strings.Contains(speech, "热血") {
		profile.PrimaryInterest = "动漫"
	} else if strings.Contains(speech, "足球") || strings.Contains(speech, "篮球") {
		profile.PrimaryInterest = "体育"
	} else if strings.Contains(speech, "手机") || strings.Contains(speech, "互联网") {
		profile.PrimaryInterest = "科技"
	} else {
		profile.PrimaryInterest = "文艺"
	}

	// 检测思维风格
	profile.ThinkingStyle = ai.detectThinkingPattern(speech)
	profile.MetaphorStyle = ai.detectMetaphorStyle(speech)

	// 分析优势
	profile.Strengths = []string{"敢于表达观点", "善于观察生活"}

	return profile
}
