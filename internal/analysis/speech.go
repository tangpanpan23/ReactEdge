package analysis

import (
	"regexp"
	"strings"
	"time"
)

// SpeechAnalyzer 语音分析器
type SpeechAnalyzer struct {
	wordCount     int
	sentenceCount int
	questionCount int
	exclamationCount int
	periodCount   int
	commaCount    int
	duration      time.Duration
	wordsPerMinute float64
}

// SpeechResult 语音分析结果
type SpeechResult struct {
	Text            string         `json:"text"`
	WordCount       int            `json:"word_count"`
	SentenceCount   int            `json:"sentence_count"`
	QuestionCount   int            `json:"question_count"`
	Duration        time.Duration  `json:"duration"`
	WordsPerMinute  float64        `json:"words_per_minute"`
	PauseCount      int            `json:"pause_count"`
	RhythmScore     int            `json:"rhythm_score"`
	ClarityScore    int            `json:"clarity_score"`
	ConfidenceScore int            `json:"confidence_score"`
}

// NewSpeechAnalyzer 创建语音分析器
func NewSpeechAnalyzer() *SpeechAnalyzer {
	return &SpeechAnalyzer{}
}

// AnalyzeText 分析文本内容（模拟语音转文字后的分析）
func (sa *SpeechAnalyzer) AnalyzeText(text string, duration time.Duration) *SpeechResult {
	sa.analyzeText(text)
	sa.duration = duration

	// 计算语速（字/分钟）
	if duration.Seconds() > 0 {
		sa.wordsPerMinute = float64(sa.wordCount) / duration.Minutes()
	}

	// 计算节奏分数（基于标点符号分布）
	rhythmScore := sa.calculateRhythmScore(text)

	// 计算清晰度分数（基于句子结构）
	clarityScore := sa.calculateClarityScore(text)

	// 计算信心分数（基于表达完整性）
	confidenceScore := sa.calculateConfidenceScore(text)

	return &SpeechResult{
		Text:            text,
		WordCount:       sa.wordCount,
		SentenceCount:   sa.sentenceCount,
		QuestionCount:   sa.questionCount,
		Duration:        duration,
		WordsPerMinute:  sa.wordsPerMinute,
		PauseCount:      sa.calculatePauseCount(text),
		RhythmScore:     rhythmScore,
		ClarityScore:    clarityScore,
		ConfidenceScore: confidenceScore,
	}
}

// analyzeText 分析文本基本特征
func (sa *SpeechAnalyzer) analyzeText(text string) {
	words := strings.Fields(text)
	sa.wordCount = len(words)

	// 计算句子数
	sentences := regexp.MustCompile(`[。！？.!?]`).FindAllString(text, -1)
	sa.sentenceCount = len(sentences)

	// 计算问号数
	questions := regexp.MustCompile(`[？?]`).FindAllString(text, -1)
	sa.questionCount = len(questions)

	// 计算感叹号数
	exclamations := regexp.MustCompile(`[！!]`).FindAllString(text, -1)
	sa.exclamationCount = len(exclamations)

	// 计算句号数
	periods := regexp.MustCompile(`[。.]`).FindAllString(text, -1)
	sa.periodCount = len(periods)

	// 计算逗号数
	commas := regexp.MustCompile(`[，,]`).FindAllString(text, -1)
	sa.commaCount = len(commas)
}

// calculateRhythmScore 计算节奏分数
func (sa *SpeechAnalyzer) calculateRhythmScore(text string) int {
	score := 50 // 基础分数

	// 根据标点符号密度调整
	totalPunctuation := sa.questionCount + sa.exclamationCount + sa.periodCount + sa.commaCount
	punctuationDensity := float64(totalPunctuation) / float64(sa.wordCount) * 100

	if punctuationDensity > 10 {
		score += 20 // 标点丰富，节奏感强
	} else if punctuationDensity < 3 {
		score -= 10 // 标点稀少，可能节奏平淡
	}

	// 检查是否有停顿词
	pauseWords := []string{"嗯", "啊", "这个", "那个", "就是说"}
	pauseCount := 0
	for _, word := range pauseWords {
		pauseCount += strings.Count(text, word)
	}

	if pauseCount > 3 {
		score -= 15 // 太多停顿词
	} else if pauseCount == 0 {
		score += 10 // 表达流畅
	}

	// 确保分数在0-100范围内
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

// calculateClarityScore 计算清晰度分数
func (sa *SpeechAnalyzer) calculateClarityScore(text string) int {
	score := 60 // 基础分数

	// 根据句子长度调整
	if sa.sentenceCount > 0 {
		avgSentenceLength := float64(sa.wordCount) / float64(sa.sentenceCount)

		if avgSentenceLength > 20 {
			score -= 20 // 句子过长，清晰度下降
		} else if avgSentenceLength < 5 {
			score -= 10 // 句子过短，可能表达不完整
		} else {
			score += 15 // 句子长度适中
		}
	}

	// 检查是否有重复词
	words := strings.Fields(text)
	wordFreq := make(map[string]int)
	for _, word := range words {
		wordFreq[word]++
	}

	repeatCount := 0
	for _, count := range wordFreq {
		if count > 3 {
			repeatCount++
		}
	}

	if repeatCount > 0 {
		score -= repeatCount * 10 // 重复词过多
	}

	// 检查表达完整性
	if sa.questionCount > 0 && sa.sentenceCount > 2 {
		score += 10 // 有问题句且有足够解释
	}

	// 确保分数在0-100范围内
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

// calculateConfidenceScore 计算信心分数
func (sa *SpeechAnalyzer) calculateConfidenceScore(text string) int {
	score := 55 // 基础分数

	// 检查表达确定性词
	confidentWords := []string{"我认为", "我相信", "我确定", "绝对", "肯定", "确实"}
	hesitantWords := []string{"可能", "也许", "大概", "不太确定", "我觉得", "应该是"}

	confidentCount := 0
	hesitantCount := 0

	for _, word := range confidentWords {
		confidentCount += strings.Count(text, word)
	}

	for _, word := range hesitantWords {
		hesitantCount += strings.Count(text, word)
	}

	score += confidentCount * 5
	score -= hesitantCount * 3

	// 检查表达长度（较长的表达通常更有信心）
	if sa.wordCount > 50 {
		score += 15
	} else if sa.wordCount < 20 {
		score -= 10
	}

	// 检查是否有具体例子
	if strings.Contains(text, "比如") || strings.Contains(text, "例如") ||
	   strings.Contains(text, "就像") || strings.Contains(text, "比如说") {
		score += 10 // 有具体例子，说明有思考深度
	}

	// 确保分数在0-100范围内
	if score > 100 {
		score = 100
	} else if score < 0 {
		score = 0
	}

	return score
}

// calculatePauseCount 计算停顿次数（基于标点符号和停顿词）
func (sa *SpeechAnalyzer) calculatePauseCount(text string) int {
	pauseWords := []string{"嗯", "啊", "这个", "那个", "就是说", "也就是说"}
	pauseCount := 0

	for _, word := range pauseWords {
		pauseCount += strings.Count(text, word)
	}

	// 标点符号也算作停顿
	pauseCount += sa.commaCount + sa.periodCount + sa.questionCount + sa.exclamationCount

	return pauseCount
}

// GetSpeechTips 获取语音建议
func (sa *SpeechAnalyzer) GetSpeechTips(result *SpeechResult) []string {
	tips := []string{}

	// 语速建议
	if result.WordsPerMinute > 200 {
		tips = append(tips, "语速稍快，建议适当放慢，让听众有时间消化观点")
	} else if result.WordsPerMinute < 120 {
		tips = append(tips, "语速稍慢，可以适当加快节奏，增加表现力")
	} else {
		tips = append(tips, "语速适中，很好地控制了表达节奏")
	}

	// 节奏建议
	if result.RhythmScore < 50 {
		tips = append(tips, "节奏可以更丰富，适当使用停顿和语调变化")
	} else if result.RhythmScore > 80 {
		tips = append(tips, "节奏控制很好，有感染力的表达方式")
	}

	// 清晰度建议
	if result.ClarityScore < 60 {
		tips = append(tips, "表达可以更清晰，建议使用更简洁的句子结构")
	} else if result.ClarityScore > 80 {
		tips = append(tips, "表达清晰明了，逻辑结构良好")
	}

	// 信心建议
	if result.ConfidenceScore < 50 {
		tips = append(tips, "可以更坚定地表达观点，减少犹豫词的使用")
	} else if result.ConfidenceScore > 75 {
		tips = append(tips, "自信的表达，观点阐述很有说服力")
	}

	return tips
}
