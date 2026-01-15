package ai

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// AIErrorHandler AI错误处理器
type AIErrorHandler struct {
	logger interface{} // 简化为interface{}，实际项目中应该使用具体的logger
}

// NewAIErrorHandler 创建AI错误处理器
func NewAIErrorHandler() *AIErrorHandler {
	return &AIErrorHandler{}
}

// HandleError 处理AI错误
func (h *AIErrorHandler) HandleError(err error, operation string) error {
	if err == nil {
		return nil
	}

	// 记录错误（实际项目中应该使用真实的logger）
	fmt.Printf("❌ AI操作%s失败: %v\n", operation, err)

	// 分类处理错误
	switch {
	case h.isRateLimitError(err):
		return h.handleRateLimitError(err)
	case h.isNetworkError(err):
		return h.handleNetworkError(err)
	case h.isInvalidRequestError(err):
		return h.handleInvalidRequestError(err)
	case h.isAuthenticationError(err):
		return h.handleAuthenticationError(err)
	default:
		return h.handleUnknownError(err)
	}
}

// isRateLimitError 判断是否为限流错误
func (h *AIErrorHandler) isRateLimitError(err error) bool {
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "rate limit") ||
		strings.Contains(errStr, "quota exceeded") ||
		strings.Contains(errStr, "too many requests")
}

// isNetworkError 判断是否为网络错误
func (h *AIErrorHandler) isNetworkError(err error) bool {
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "connection") ||
		strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "network") ||
		strings.Contains(errStr, "dial tcp")
}

// isInvalidRequestError 判断是否为无效请求错误
func (h *AIErrorHandler) isInvalidRequestError(err error) bool {
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "invalid") ||
		strings.Contains(errStr, "bad request") ||
		strings.Contains(errStr, "malformed")
}

// isAuthenticationError 判断是否为认证错误
func (h *AIErrorHandler) isAuthenticationError(err error) bool {
	errStr := strings.ToLower(err.Error())
	return strings.Contains(errStr, "unauthorized") ||
		strings.Contains(errStr, "authentication") ||
		strings.Contains(errStr, "api key")
}

// handleRateLimitError 处理限流错误
func (h *AIErrorHandler) handleRateLimitError(err error) error {
	fmt.Println("⚠️ 触发API限流，建议稍后重试")
	return errors.New("api_rate_limited")
}

// handleNetworkError 处理网络错误
func (h *AIErrorHandler) handleNetworkError(err error) error {
	fmt.Println("⚠️ 网络连接错误，请检查网络连接")
	return errors.New("network_error")
}

// handleInvalidRequestError 处理无效请求错误
func (h *AIErrorHandler) handleInvalidRequestError(err error) error {
	fmt.Println("⚠️ 请求参数无效，请检查输入")
	return errors.New("invalid_request")
}

// handleAuthenticationError 处理认证错误
func (h *AIErrorHandler) handleAuthenticationError(err error) error {
	fmt.Println("⚠️ API认证失败，请检查API密钥配置")
	return errors.New("authentication_failed")
}

// handleUnknownError 处理未知错误
func (h *AIErrorHandler) handleUnknownError(err error) error {
	fmt.Printf("⚠️ 未知AI错误: %v\n", err)
	return errors.New("unknown_ai_error")
}

// FallbackResponse 获取降级响应
func (h *AIErrorHandler) FallbackResponse(operation string) interface{} {
	switch operation {
	case "analyze_image":
		return h.defaultImageAnalysis()
	case "generate_questions":
		return h.defaultQuestions()
	case "polish_note":
		return h.defaultPolishedNote()
	case "generate_reaction_templates":
		return h.defaultReactionTemplates()
	case "analyze_expression_style":
		return h.defaultStyleAnalysis()
	case "simulate_debate":
		return h.defaultDebateSimulation()
	case "evaluate_reaction":
		return h.defaultReactionEvaluation()
	default:
		return h.defaultResponse()
	}
}

// 默认降级响应实现
func (h *AIErrorHandler) defaultImageAnalysis() *ImageAnalysisResult {
	return &ImageAnalysisResult{
		ObjectName:     "分析对象",
		Category:       "general",
		Description:    "AI服务暂时不可用，提供模拟分析结果",
		Confidence:     0.5,
		KeyFeatures:    []string{"模拟分析"},
		ScientificName: "未知",
	}
}

func (h *AIErrorHandler) defaultQuestions() []Question {
	return []Question{
		{
			Content:    "这个问题很有价值，让我们一起探讨",
			Type:       "scenario",
			Difficulty: "basic",
			Purpose:    "AI服务降级模式",
		},
	}
}

func (h *AIErrorHandler) defaultPolishedNote() *PolishedNote {
	return &PolishedNote{
		Title:       "沟通记录",
		Summary:     "AI服务暂时不可用",
		KeyPoints:   []string{"记录已保存"},
		Questions:   []string{"稍后重试AI分析"},
		FormattedText: "原始内容已保存",
	}
}

func (h *AIErrorHandler) defaultReactionTemplates() []ReactionTemplate {
	return []ReactionTemplate{
		{
			Scenario:    "通用场景",
			Steps:       []string{"保持冷静", "认真倾听", "适当回应"},
			KeyPhrases:  []string{"我理解你的观点", "让我们一起探讨"},
			StyleNotes:  "AI服务降级模式",
		},
	}
}

func (h *AIErrorHandler) defaultStyleAnalysis() *StyleAnalysis {
	return &StyleAnalysis{
		PersonName: "分析对象",
		LanguageFeatures: map[string]interface{}{
			"clarity": "清晰度分析",
		},
		ThinkingPatterns: map[string]interface{}{
			"logic": "逻辑分析",
		},
		CommunicationStrategy: map[string]interface{}{
			"strategy": "策略分析",
		},
		PersonalTraits: map[string]interface{}{
			"traits": "特征分析",
		},
		OverallScore: 7.0,
		StyleTags:    []string{"分析中"},
	}
}

func (h *AIErrorHandler) defaultDebateSimulation() *DebateSimulation {
	return &DebateSimulation{
		Scenario:        "辩论场景",
		OpponentOpening: "这是我的观点",
		InteractionRounds: []DebateRound{
			{
				RoundNumber:      1,
				OpponentMove:     "不同意见",
				ExpectedResponse: "理解并回应",
				ReactionTips:     "保持专业",
			},
		},
		KeyReactionPoints: []string{"关键点"},
		StyleSuggestions:  []string{"专业回应"},
		Difficulty:        1,
	}
}

func (h *AIErrorHandler) defaultReactionEvaluation() *ReactionEvaluation {
	return &ReactionEvaluation{
		ContentQuality: EvaluationItem{
			Score:       7.0,
			Description: "内容质量良好",
			Suggestions: []string{"保持当前水平"},
		},
		StyleConformity: EvaluationItem{
			Score:       6.5,
			Description: "风格符合度一般",
			Suggestions: []string{"适当调整风格"},
		},
		ReactionSpeed: EvaluationItem{
			Score:       7.5,
			Description: "反应速度良好",
			Suggestions: []string{"继续保持"},
		},
		CommunicationEffect: EvaluationItem{
			Score:       7.0,
			Description: "沟通效果良好",
			Suggestions: []string{"继续优化"},
		},
		OverallScore: 7.0,
		Strengths:     []string{"基础扎实"},
		Improvements:  []string{"细节优化"},
	}
}

func (h *AIErrorHandler) defaultResponse() interface{} {
	return map[string]string{
		"status":  "degraded",
		"message": "AI服务暂时不可用，使用降级模式",
	}
}

// AICircuitBreaker AI熔断器
type AICircuitBreaker struct {
	failureCount int
	lastFailure  time.Time
	state        string // "closed", "open", "half-open"
	timeout      time.Duration
	maxFailures  int
}

// NewAICircuitBreaker 创建AI熔断器
func NewAICircuitBreaker(maxFailures int, timeout time.Duration) *AICircuitBreaker {
	return &AICircuitBreaker{
		state:       "closed",
		maxFailures: maxFailures,
		timeout:     timeout,
	}
}

// Call 执行带熔断器的调用
func (cb *AICircuitBreaker) Call(operation func() error) error {
	if cb.state == "open" {
		if time.Since(cb.lastFailure) > cb.timeout {
			cb.state = "half-open"
		} else {
			return errors.New("circuit breaker is open")
		}
	}

	err := operation()
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()
	return nil
}

// recordFailure 记录失败
func (cb *AICircuitBreaker) recordFailure() {
	cb.failureCount++
	cb.lastFailure = time.Now()

	if cb.failureCount >= cb.maxFailures {
		cb.state = "open"
		fmt.Println("🔌 AI熔断器开启，暂时停止AI调用")
	}
}

// recordSuccess 记录成功
func (cb *AICircuitBreaker) recordSuccess() {
	cb.failureCount = 0
	cb.state = "closed"
	if cb.state == "half-open" {
		fmt.Println("🔄 AI熔断器半开，恢复正常调用")
	}
}

// IsOpen 检查熔断器是否开启
func (cb *AICircuitBreaker) IsOpen() bool {
	return cb.state == "open"
}
