package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"reactedge/internal/ai"
	"reactedge/internal/analysis"
	"reactedge/internal/challenge"
)

func main() {
	fmt.Println("🎤 AI酷表达实验室 · 韩寒特训版 - 演示版本")
	fmt.Println("================================================")
	fmt.Println()

	// 初始化AI引擎
	hanAI := ai.NewHanStyleAI()
	fmt.Printf("✅ 韩寒表达引擎已加载，包含 %d 个表达模式\n\n", len(hanAI.GetExpressionPatterns()))

	// 初始化挑战管理器
	challengeManager := challenge.NewManager(hanAI)

	// 模拟用户ID
	userID := "demo_user"

	// 开始挑战
	fmt.Println("🔥 欢迎来到【酷表达实验室】· 韩寒特训版")
	fmt.Println("🎯 今日挑战：课堂突击提问")
	fmt.Println()

	// 阶段1: AI解构韩寒表达法
	fmt.Println("🧠 AI解构【韩寒表达法】三大武器：")
	fmt.Println()
	fmt.Println("1. **反常规视角** 🌪️")
	fmt.Println("   普通人：赞美书店变多 → 文化繁荣")
	fmt.Println("   韩寒式：\"当书店开始比拼装修而不是书目，这和奶茶店比杯子颜值有什么区别？\"")
	fmt.Println()
	fmt.Println("2. **精准文化类比** 🎬")
	fmt.Println("   把抽象概念变成具体场景：")
	fmt.Println("   \"这就像电影院里全是爆米花味，但没人在意放的是什么电影\"")
	fmt.Println()
	fmt.Println("3. **节奏打断技巧** ⚡")
	fmt.Println("   在对方预期处突然转折：")
	fmt.Println("   \"很多人说这是好事...(停顿)但好事有时候是最可怕的陷阱\"")
	fmt.Println()
	fmt.Println("🛠️ 你的专属工具箱：")
	fmt.Println("【反问模板】\"难道...就代表...?\"")
	fmt.Println("【类比模板】\"这就像...其实不过是...\"")
	fmt.Println("【转折模板】\"表面上看是...实际上暴露了...\"")
	fmt.Println()

	// 等待用户继续
	waitForEnter("按回车键继续到个性化模板生成...")

	// 阶段2: 生成个性化模板
	fmt.Println("🤖 AI为你生成【个性化应答模板】")
	fmt.Println()

	// 模拟用户输入来探测兴趣
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("为了生成更适合你的模板，请简单描述一下你的兴趣爱好：")
	fmt.Println("(例如：我喜欢玩游戏、看动漫、打篮球、编程、学音乐...)")
	fmt.Print("你的兴趣：")

	interestInput, _ := reader.ReadString('\n')
	interestInput = strings.TrimSpace(interestInput)

	// 从输入中探测用户画像
	userProfile := hanAI.DetectUserProfile(interestInput)
	fmt.Printf("✅ AI探测到你的偏好：%s\n", userProfile.PrimaryInterest)
	fmt.Printf("✅ 为你生成【%s版】应答模板：\n\n", getInterestDisplayName(userProfile.PrimaryInterest))

	// 生成个性化模板
	template := hanAI.GeneratePersonalizedTemplate(userProfile, "\"你对网红书店遍地开花这种现象，怎么看？\"")
	fmt.Println(template)
	fmt.Println()

	fmt.Println("💡 你的表达框架：")
	fmt.Println("（1）游戏类比切入 → 吸引同龄人")
	fmt.Println("（2）对比转折 → 展现思辨")
	fmt.Println("（3）现象本质 → 提升深度")
	fmt.Println("（4）犀利反问 → 留下印象")
	fmt.Println()

	waitForEnter("按回车键开始你的表达挑战...")

	// 阶段3: 表达挑战
	fmt.Println("🎤 现在请用你的风格回答！")
	fmt.Println("⏱️ 15秒思考，45秒回答")
	fmt.Println("📚 场景：语文课上，老师突然点名：\"你对网红书店遍地开花这种现象，怎么看？\"")
	fmt.Println()
	fmt.Println("请在这里输入你的回答：")

	// 模拟倒计时
	fmt.Print("开始思考 (15秒)... ")
	for i := 15; i > 0; i-- {
		fmt.Printf("%d ", i)
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("时间到！")
	fmt.Println()

	// 获取用户回答
	speech, _ := reader.ReadString('\n')
	speech = strings.TrimSpace(speech)

	if speech == "" {
		fmt.Println("😅 看来你需要更多时间思考。让我们看看AI的示例回答：")
		speech = "我觉得这就像《塞尔达》里到处是神庙但解谜都很简单——数量多了，质量却被稀释了。表面上是书店繁荣，实际上暴露了我们用'打卡'代替'阅读'的虚荣。如果书店变成拍照背景板，那和游戏里的贴图BUG有什么区别？"
		fmt.Printf("示例回答：%s\n\n", speech)
	}

	// 提交回答并分析
	fmt.Println("🎯 正在分析你的表达DNA...")
	time.Sleep(2 * time.Second)

	// 阶段4: DNA分析结果
	dna := hanAI.AnalyzeExpressionDNA(speech, userProfile)

	// 语音分析
	speechAnalyzer := analysis.NewSpeechAnalyzer()
	duration := 45 * time.Second // 假设45秒的回答时间
	speechResult := speechAnalyzer.AnalyzeText(speech, duration)
	speechTips := speechAnalyzer.GetSpeechTips(speechResult)

	fmt.Println("📊 你的【表达DNA分析报告】")
	fmt.Println()
	fmt.Printf("🔥 犀利指数：%d/100\n", dna.SharpenessScore)

	fmt.Println("🎯 个性标签：")
	for _, tag := range dna.PersonalityTags {
		fmt.Printf("   • %s\n", tag)
	}

	fmt.Println("💎 发现你的独家表达模式：")
	for _, pattern := range dna.UniquePatterns {
		fmt.Printf("   • %s\n", pattern)
	}

	fmt.Printf("🧠 思维模式：%s\n", dna.ThinkingPattern)
	fmt.Printf("🎨 类比风格：%s\n", dna.MetaphorStyle)
	fmt.Printf("🎵 节奏特征：%s\n", dna.RhythmSignature)
	fmt.Printf("✨ 独特性分数：%d/100\n", dna.UniquenessScore)

	fmt.Println()
	fmt.Println("🎤 语音表现分析：")
	fmt.Printf("   • 字数：%d字\n", speechResult.WordCount)
	fmt.Printf("   • 句子数：%d句\n", speechResult.SentenceCount)
	fmt.Printf("   • 语速：%.1f字/分钟\n", speechResult.WordsPerMinute)
	fmt.Printf("   • 节奏分数：%d/100\n", speechResult.RhythmScore)
	fmt.Printf("   • 清晰度分数：%d/100\n", speechResult.ClarityScore)
	fmt.Printf("   • 信心分数：%d/100\n", speechResult.ConfidenceScore)

	fmt.Println("🎯 语音优化建议：")
	for _, tip := range speechTips {
		fmt.Printf("   • %s\n", tip)
	}

	fmt.Println()
	fmt.Println("🆙 AI综合优化建议：")
	for _, rec := range dna.Recommendations {
		fmt.Printf("   • %s\n", rec)
	}

	fmt.Printf("🎮 明日挑战预告：%s\n", dna.NextChallenge)
	fmt.Println()

	fmt.Println("🎉 挑战完成！")
	fmt.Println("这不止是一次训练。AI发现了你的独特表达天赋，明天的挑战会围绕这个优势继续设计。")
	fmt.Println()
	fmt.Println("💪 传统教育试图把所有人教成同一个'优秀模样'，而我们的AI引擎，专门发现并放大你独有的表达天赋。")
	fmt.Println("🌟 每天3分钟，不是学习套路，而是让你的个性表达变得更犀利、更有影响力。")
	fmt.Println()
	fmt.Println("感谢体验 AI酷表达实验室！有任何问题欢迎反馈。")
}

func waitForEnter(message string) {
	fmt.Println(message)
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getInterestDisplayName(interest string) string {
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
