package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"reactedge/internal/ai"
	"reactedge/internal/analysis"
	"reactedge/internal/challenge"
)

func main() {
	fmt.Println("🎤 AI酷表达实验室 · 言刃 ReactEdge - 演示版本")
	fmt.Println("================================================")
	fmt.Println()

	// 初始化AI引擎
	hanAI := ai.NewHanStyleAI()
	fmt.Printf("✅ 多风格表达引擎已加载，包含 %d 个表达模式\n", len(hanAI.GetExpressionPatterns()))
	fmt.Println("   支持康辉、韩寒、董卿、黄执中等顶尖人物风格")
	fmt.Println()

	// 初始化挑战管理器
	challengeManager := challenge.NewManager(hanAI)

	// 模拟用户ID
	userID := "demo_user"

	// 场景选择
	fmt.Println("🏫 请选择训练场景：")
	fmt.Println()
	fmt.Println("1. 📚 课堂挑战 - 语文课突发提问")
	fmt.Println("2. 🚪 电梯挑战 - 30秒向CEO汇报")
	fmt.Println("3. 🔥 危机应对 - 临时记者会")
	fmt.Println("4. 💰 投资答辩 - 面对投资人质询")
	fmt.Println("5. 👨‍👩‍👧‍👦 家庭调解 - 化解亲人矛盾")
	fmt.Println("6. 🌍 跨文化沟通 - 处理文化误解")
	fmt.Println()

	selectedScenario := ""

	for {
		fmt.Print("请选择场景 (1-6): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			selectedScenario = "classroom"
			fmt.Println("✅ 已选择：课堂挑战 - 突发提问应对训练")
		case "2":
			selectedScenario = "elevator"
			fmt.Println("✅ 已选择：电梯挑战 - 30秒汇报训练")
		case "3":
			selectedScenario = "crisis"
			fmt.Println("✅ 已选择：危机应对 - 舆情处理训练")
		case "4":
			selectedScenario = "investment"
			fmt.Println("✅ 已选择：投资答辩 - 融资沟通训练")
		case "5":
			selectedScenario = "family"
			fmt.Println("✅ 已选择：家庭调解 - 情感智慧训练")
		case "6":
			selectedScenario = "cultural"
			fmt.Println("✅ 已选择：跨文化沟通 - 文化敏感训练")
		default:
			fmt.Println("❌ 无效选择，请输入1-6之间的数字")
			continue
		}
		break
	}

	fmt.Println()
	fmt.Println("🎯 请选择你的目标表达风格：")

	// 风格选择
	fmt.Println("🔥 欢迎来到【酷表达实验室】· 言刃 ReactEdge")
	fmt.Println("🎯 请选择你的目标表达风格：")
	fmt.Println()
	fmt.Println("1. 📰 康辉（央视型）- 沉稳权威，适合正式场合")
	fmt.Println("2. 🌪️ 韩寒（犀利型）- 反常规视角，适合辩论表达")
	fmt.Println("3. 💝 董卿（共情型）- 情感共鸣，适合沟通交流")
	fmt.Println("4. 🧠 黄执中（辩论型）- 逻辑重构，适合理性分析")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	selectedStyle := ""

	for {
		fmt.Print("请选择风格 (1-4): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			selectedStyle = "kanghui"
			fmt.Println("✅ 已选择：康辉（央视型）- 沉稳权威风格")
		case "2":
			selectedStyle = "hanhan"
			fmt.Println("✅ 已选择：韩寒（犀利型）- 反常规视角风格")
		case "3":
			selectedStyle = "dongqing"
			fmt.Println("✅ 已选择：董卿（共情型）- 情感共鸣风格")
		case "4":
			selectedStyle = "huangzhizhong"
			fmt.Println("✅ 已选择：黄执中（辩论型）- 逻辑重构风格")
		default:
			fmt.Println("❌ 无效选择，请输入1-4之间的数字")
			continue
		}
		break
	}

	fmt.Println()
	fmt.Printf("🎯 今日挑战：%s\n", getScenarioDisplayName(selectedScenario))
	fmt.Println()

	// 显示场景描述和挑战要求
	switch selectedScenario {
	case "classroom":
		fmt.Println("📚 场景：语文课上，老师突然点名：\"你对网红书店遍地开花这种现象，怎么看？\"")
		fmt.Println("⏰ 时间：45秒回答，要有自己的观点和见解")
	case "elevator":
		fmt.Println("🚪 场景：电梯偶遇CEO，他问：\"你觉得我们这个项目最核心的价值是什么？\"")
		fmt.Println("⏰ 时间：30秒回答，要突出重点，逻辑清晰")
	case "crisis":
		fmt.Println("🔥 场景：公司突发负面舆情，你作为发言人召开临时记者会")
		fmt.Println("   记者问：\"请问公司对此次事件如何回应？\"")
		fmt.Println("⏰ 时间：60秒回答，要安抚情绪，提供解决方案")
	case "investment":
		fmt.Println("💰 场景：投资人会议上，第3位投资人问：\"你们的商业模式真的可持续吗？\"")
		fmt.Println("⏰ 时间：45秒回答，要有数据支撑，展现信心")
	case "family":
		fmt.Println("👨‍👩‍👧‍👦 场景：父母吵架后，母亲生气地说：\"你爸根本不在乎这个家！\"")
		fmt.Println("   你需要：引导双方冷静，寻找共识")
		fmt.Println("⏰ 时间：60秒引导，要共情倾听，理性分析")
	case "cultural":
		fmt.Println("🌍 场景：国际会议上，外方合作伙伴说：\"你们的方式太官僚了，完全没有效率\"")
		fmt.Println("   背景：对方来自崇尚个人主义文化")
		fmt.Println("⏰ 时间：45秒回应，要尊重差异，寻求理解")
	}

	fmt.Println()

	// 阶段1: AI解构目标风格表达法
	fmt.Printf("🧠 AI解构【%s】表达武器：\n", getStyleDisplayName(selectedStyle))
	fmt.Println()

	switch selectedStyle {
	case "kanghui":
		fmt.Println("1. **事实数据构建** 📊")
		fmt.Println("   用政策文件和数据支撑观点：")
		fmt.Println("   \"事实上，根据最新统计数据显示...\"")
		fmt.Println()
		fmt.Println("2. **三层结构推进** 🏗️")
		fmt.Println("   国家-社会-个体层次递进：")
		fmt.Println("   \"从国家层面来看...社会角度分析...对个体而言...\"")
		fmt.Println()
		fmt.Println("3. **关键词停顿** ⏸️")
		fmt.Println("   关键处适当停顿增强权威感：")
		fmt.Println("   \"这个数据...(停顿)非常重要\"")
		fmt.Println()
		fmt.Println("🛠️ 康辉式工具箱：")
		fmt.Println("【数据引用】\"事实上...数据显示...\"")
		fmt.Println("【层次递进】\"从...层面来看...\"")
		fmt.Println("【权威停顿】在关键数据后停顿")

	case "dongqing":
		fmt.Println("1. **情感共鸣** 💝")
		fmt.Println("   建立情感连接：")
		fmt.Println("   \"我完全能够感受到大家的这种心情...\"")
		fmt.Println()
		fmt.Println("2. **个人故事引入** 📖")
		fmt.Println("   用故事拉近距离：")
		fmt.Println("   \"让我想起曾经的一个经历...\"")
		fmt.Println()
		fmt.Println("3. **价值升华** ✨")
		fmt.Println("   引导至更高层次：")
		fmt.Println("   \"这不仅仅是...更是关乎...\"")
		fmt.Println()
		fmt.Println("🛠️ 董卿式工具箱：")
		fmt.Println("【共鸣表达】\"我能感受到...\"")
		fmt.Println("【故事导入】\"让我想起...\"")
		fmt.Println("【价值提升】\"这关乎...\"")

	case "huangzhizhong":
		fmt.Println("1. **问题重定义** 🔄")
		fmt.Println("   重新框架化问题：")
		fmt.Println("   \"这个问题不应该这样问，我们应该思考...\"")
		fmt.Println()
		fmt.Println("2. **利害分析** ⚖️")
		fmt.Println("   分析各方利益得失：")
		fmt.Println("   \"这样做对谁有利？对谁有弊？\"")
		fmt.Println()
		fmt.Println("3. **选择构建** 🎯")
		fmt.Println("   构建清晰的选择路径：")
		fmt.Println("   \"我们有三个选择：第一...第二...第三...\"")
		fmt.Println()
		fmt.Println("🛠️ 黄执中式工具箱：")
		fmt.Println("【问题重构】\"真正的问题是...\"")
		fmt.Println("【利弊分析】\"这样做的好处是...坏处是...\"")
		fmt.Println("【选择框架】\"我们面临的选择是...\"")

	default: // hanhan
		fmt.Println("1. **反常规视角** 🌪️")
		fmt.Println("   从意想不到的角度切入：")
		fmt.Println("   \"当书店开始比拼装修而不是书目，这和奶茶店比杯子颜值有什么区别？\"")
		fmt.Println()
		fmt.Println("2. **精准文化类比** 🎬")
		fmt.Println("   把抽象概念变成具体场景：")
		fmt.Println("   \"这就像电影院里全是爆米花味，但没人在意放的是什么电影\"")
		fmt.Println()
		fmt.Println("3. **节奏打断技巧** ⚡")
		fmt.Println("   在对方预期处突然转折：")
		fmt.Println("   \"很多人说这是好事...(停顿)但好事有时候是最可怕的陷阱\"")
		fmt.Println()
		fmt.Println("🛠️ 韩寒式工具箱：")
		fmt.Println("【反问模板】\"难道...就代表...?\"")
		fmt.Println("【类比模板】\"这就像...其实不过是...\"")
		fmt.Println("【转折模板】\"表面上看是...实际上暴露了...\"")
	}

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

	// 根据场景调整时间要求
	timeLimit := 45
	switch selectedScenario {
	case "elevator":
		timeLimit = 30
		fmt.Println("⏱️ 10秒思考，30秒回答")
	case "crisis", "family":
		timeLimit = 60
		fmt.Println("⏱️ 15秒思考，60秒回答")
	default:
		fmt.Println("⏱️ 15秒思考，45秒回答")
	}

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
	duration := time.Duration(timeLimit) * time.Second // 根据场景调整时间
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

	// 根据选择的风格给出针对性分析
	fmt.Println()
	fmt.Printf("🎭 与【%s】风格相似度分析：\n", getStyleDisplayName(selectedStyle))

	rand.Seed(time.Now().UnixNano())
	switch selectedStyle {
	case "kanghui":
		fmt.Printf("   • 权威感：%d/100 (数据引用和停顿运用)\n", 75+rand.Intn(20))
		fmt.Printf("   • 层次感：%d/100 (结构化表达能力)\n", 70+rand.Intn(25))
		fmt.Printf("   • 专业度：%d/100 (事实支撑程度)\n", 72+rand.Intn(23))
	case "dongqing":
		fmt.Printf("   • 共情力：%d/100 (情感连接能力)\n", 78+rand.Intn(17))
		fmt.Printf("   • 亲和力：%d/100 (拉近距离技巧)\n", 80+rand.Intn(15))
		fmt.Printf("   • 温暖感：%d/100 (人文关怀程度)\n", 75+rand.Intn(20))
	case "huangzhizhong":
		fmt.Printf("   • 逻辑性：%d/100 (推理结构完整度)\n", 82+rand.Intn(13))
		fmt.Printf("   • 重构力：%d/100 (问题框架重定义)\n", 79+rand.Intn(16))
		fmt.Printf("   • 辩证性：%d/100 (多角度分析能力)\n", 76+rand.Intn(19))
	default: // hanhan
		fmt.Printf("   • 犀利度：%d/100 (反常规视角运用)\n", 85+rand.Intn(10))
		fmt.Printf("   • 类比力：%d/100 (文化类比创意度)\n", 80+rand.Intn(15))
		fmt.Printf("   • 冲击力：%d/100 (转折技巧掌握)\n", 78+rand.Intn(17))
	}

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
	// 根据选择的场景和风格给出个性化训练建议
	fmt.Println()
	fmt.Printf("🎯 你的【%s】风格在【%s】场景的训练计划：\n",
		getStyleDisplayName(selectedStyle), getScenarioDisplayName(selectedScenario))

	// 场景特定的建议
	switch selectedScenario {
	case "elevator":
		fmt.Println("🏢 电梯汇报要点：时间宝贵，突出核心价值，逻辑简洁")
		fmt.Println("💡 建议练习：每天用30秒总结一个项目的核心价值")
	case "crisis":
		fmt.Println("🔥 危机应对要点：先安抚情绪，再提供事实，最后展望未来")
		fmt.Println("💡 建议练习：准备3套不同程度的危机应对话术")
	case "investment":
		fmt.Println("💰 投资答辩要点：数据支撑信心，风险透明，愿景清晰")
		fmt.Println("💡 建议练习：准备项目5大关键数据的快速调用")
	case "family":
		fmt.Println("👨‍👩‍👧‍👦 家庭调解要点：倾听共情，理性分析，引导共识")
		fmt.Println("💡 建议练习：练习\"我理解你的感受，同时...\"的句式")
	case "cultural":
		fmt.Println("🌍 跨文化要点：尊重差异，寻求共性，建立桥梁")
		fmt.Println("💡 建议练习：学习不同文化的沟通偏好和禁忌")
	}

	fmt.Println()
	fmt.Printf("🎯 通用【%s】风格提升计划：\n", getStyleDisplayName(selectedStyle))

	switch selectedStyle {
	case "kanghui":
		fmt.Println("   📅 第1-3天：练习权威停顿 - 在关键数据后停顿2秒")
		fmt.Println("   📅 第4-7天：掌握三层结构 - 国家→社会→个体的递进表达")
		fmt.Println("   📅 第8-14天：数据引用特训 - 快速调取3个相关数据支撑观点")
	case "dongqing":
		fmt.Println("   📅 第1-3天：情感共鸣练习 - 从\"我能感受到\"开始每段表达")
		fmt.Println("   📅 第4-7天：故事导入技巧 - 准备5个生活故事用于拉近距离")
		fmt.Println("   📅 第8-14天：价值升华训练 - 将具体问题提升到普遍价值层面")
	case "huangzhizhong":
		fmt.Println("   📅 第1-3天：问题重定义 - 练习用\"真正的问题是...\"重新框架")
		fmt.Println("   📅 第4-7天：利害分析法 - 每个观点分析正反两方面")
		fmt.Println("   📅 第8-14天：选择构建训练 - 为复杂问题构建清晰的选择路径")
	default: // hanhan
		fmt.Println("   📅 第1-3天：反常规视角 - 练习从相反角度切入问题")
		fmt.Println("   📅 第4-7天：文化类比特训 - 每天创造3个新颖的文化类比")
		fmt.Println("   📅 第8-14天：转折冲击训练 - 掌握在预期处突然转折的时机")
	}

	fmt.Println()
	fmt.Println("🗓️  21天后，你将能熟练运用这种风格，在各种场合自信表达！")
	fmt.Println()
	fmt.Println("感谢体验 AI酷表达实验室 · 言刃 ReactEdge！有任何问题欢迎反馈。")
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

func getStyleDisplayName(style string) string {
	names := map[string]string{
		"kanghui":       "康辉（央视型）",
		"hanhan":        "韩寒（犀利型）",
		"dongqing":      "董卿（共情型）",
		"huangzhizhong": "黄执中（辩论型）",
	}

	if name, ok := names[style]; ok {
		return name
	}
	return "韩寒（犀利型）"
}

func getScenarioDisplayName(scenario string) string {
	names := map[string]string{
		"classroom":  "课堂突击提问",
		"elevator":   "电梯汇报挑战",
		"crisis":     "舆情危机应对",
		"investment": "投资人答辩",
		"family":     "家庭矛盾调解",
		"cultural":   "跨文化沟通",
	}

	if name, ok := names[scenario]; ok {
		return name
	}
	return "课堂突击提问"
}
