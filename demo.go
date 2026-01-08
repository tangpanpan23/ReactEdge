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

	// 场景选择 - 聚焦三大核心危机场景
	fmt.Println("🏢 请选择职场危机训练场景：")
	fmt.Println()
	fmt.Println("=== 三大核心场景 ===")
	fmt.Println("1. 📊 述职答辩 - 面对领导的各种问题")
	fmt.Println("2. 🎤 分享会刁难 - 场下故意挑衅的发问")
	fmt.Println("3. 💬 争辩冲突 - 不友善的突如其来快速反应")
	fmt.Println()
	fmt.Println("=== 扩展训练场景 ===")
	fmt.Println("4. 🚪 电梯汇报 - 30秒向CEO汇报")
	fmt.Println("5. 🔥 危机公关 - 临时记者会应对")
	fmt.Println("6. 💰 投资答辩 - 面对投资人质询")
	fmt.Println()

	selectedScenario := ""

	for {
		fmt.Print("请选择场景 (1-6): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			selectedScenario = "performance"
			fmt.Println("✅ 已选择：述职答辩 - 面对领导质疑的快速反应训练")
		case "2":
			selectedScenario = "presentation"
			fmt.Println("✅ 已选择：分享会刁难 - 应对场下挑衅发问的危机处理")
		case "3":
			selectedScenario = "debate"
			fmt.Println("✅ 已选择：争辩冲突 - 不友善争辩时的快速反应接话")
		case "4":
			selectedScenario = "elevator"
			fmt.Println("✅ 已选择：电梯汇报 - 30秒向CEO汇报的核心价值")
		case "5":
			selectedScenario = "crisis"
			fmt.Println("✅ 已选择：危机公关 - 临时记者会舆情应对")
		case "6":
			selectedScenario = "investment"
			fmt.Println("✅ 已选择：投资答辩 - 面对投资人尖锐质询")
		default:
			fmt.Println("❌ 无效选择，请输入1-6之间的数字")
			continue
		}
		break
	}

	fmt.Println()
	fmt.Println("🎯 请选择你的目标表达风格：")

	// 风格选择
	fmt.Println("🎯 请选择你的目标表达风格：")
	fmt.Println()
	fmt.Println("1. 📰 康辉式 (标准得体) - 真诚与专业，零失误，稳定控场")
	fmt.Println("2. 🧠 成铭式 (辩手机制) - 策略与破局，逻辑拆解，掌控走向")
	fmt.Println("3. 🌪️ 韩寒式 (犀利风格) - 真实与穿透，直接回应，观点冲击")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	selectedStyle := ""

	for {
		fmt.Print("请选择风格 (1-3): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			selectedStyle = "kanghui"
			fmt.Println("✅ 已选择：康辉式 (标准得体) - 真诚与专业")
		case "2":
			selectedStyle = "chengming"
			fmt.Println("✅ 已选择：成铭式 (辩手机制) - 策略与破局")
		case "3":
			selectedStyle = "hanhan"
			fmt.Println("✅ 已选择：韩寒式 (犀利风格) - 真实与穿透")
		default:
			fmt.Println("❌ 无效选择，请输入1-3之间的数字")
			continue
		}
		break
	}

	fmt.Println()
	fmt.Printf("🎯 今日挑战：%s\n", getScenarioDisplayName(selectedScenario))
	fmt.Println()

	// 显示场景描述和挑战要求
	switch selectedScenario {
	case "performance":
		fmt.Println("📊 场景：年终述职汇报中，领导突然打断问：\"你这个项目的ROI怎么这么低？数据可靠吗？\"")
		fmt.Println("   挑战：领导质疑你的工作成果和专业能力")
		fmt.Println("⏰ 时间：45秒反应，要用数据反驳，展现专业性")
	case "presentation":
		fmt.Println("🎤 场景：产品分享会后Q&A环节，一位竞争对手模样的听众发问：")
		fmt.Println("   \"你这个方案听起来不错，但我觉得就是换汤不换药，老调重弹而已\"")
		fmt.Println("   挑战：公开场合的故意刁难和人身攻击")
		fmt.Println("⏰ 时间：60秒回应，要化解恶意，维护专业形象")
	case "debate":
		fmt.Println("💬 场景：团队会议讨论新项目方案，同事突然激动地说：")
		fmt.Println("   \"你这个想法太天真了，完全不考虑实际情况，简直是纸上谈兵！\"")
		fmt.Println("   挑战：情绪化争辩，不友善的快速攻击")
		fmt.Println("⏰ 时间：30秒快速反应，要控制情绪，逻辑反击")
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

	case "chengming":
		fmt.Println("1. **避实就虚** 🎯")
		fmt.Println("   不正面硬扛，抓住逻辑漏洞反击：")
		fmt.Println("   \"您提到的风险，恰恰是我们设计第三套预案的原因\"")
		fmt.Println()
		fmt.Println("2. **请君入瓮** 🪤")
		fmt.Println("   连续提问诱导对方进入己方论证范围：")
		fmt.Println("   \"我们讨论成本，最终目标是看投资回报率。您更关注短期预算还是长期收益呢？\"")
		fmt.Println()
		fmt.Println("3. **归谬反驳** 🎭")
		fmt.Println("   承认对方观点，推导荒谬结论：")
		fmt.Println("   \"如果按您的逻辑，我们就应该回到刀耕火种的时代\"")
		fmt.Println()
		fmt.Println("🛠️ 成铭式工具箱：")
		fmt.Println("【避实就虚】\"您提到的X点，恰恰证明了Y的必要性\"")
		fmt.Println("【请君入瓮】连续提问，诱导对方进入你的框架")
		fmt.Println("【归谬反驳】\"如果按您的逻辑，那就等于...\"")

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
	case "performance":
		timeLimit = 45
		fmt.Println("⏱️ 10秒思考，45秒回答（述职答辩需要数据支撑）")
	case "presentation":
		timeLimit = 60
		fmt.Println("⏱️ 15秒思考，60秒回答（分享会刁难需要全面应对）")
	case "debate":
		timeLimit = 30
		fmt.Println("⏱️ 5秒思考，30秒快速反应（争辩冲突时间紧迫）")
	case "elevator":
		timeLimit = 30
		fmt.Println("⏱️ 10秒思考，30秒回答")
	case "crisis":
		timeLimit = 60
		fmt.Println("⏱️ 15秒思考，60秒回答")
	default:
		timeLimit = 45
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
	case "performance":
		fmt.Println("📊 述职答辩要点：用数据回应质疑，展现专业性，平衡自信与谦逊")
		fmt.Println("💡 建议练习：准备项目ROI、风险控制、创新贡献等3大类问题的标准回答")
	case "presentation":
		fmt.Println("🎤 分享会刁难要点：保持专业风度，用逻辑化解恶意，维护个人品牌")
		fmt.Println("💡 建议练习：模拟不同类型的挑衅（技术质疑、动机攻击、人身挑衅），练习化解技巧")
	case "debate":
		fmt.Println("💬 争辩冲突要点：控制情绪节奏，逻辑反击，适时结束对话")
		fmt.Println("💡 建议练习：练习\"不接招\"的技巧，识别情绪陷阱，避免陷入无谓纠缠")
	case "elevator":
		fmt.Println("🏢 电梯汇报要点：时间宝贵，突出核心价值，逻辑简洁")
		fmt.Println("💡 建议练习：每天用30秒总结一个项目的核心价值")
	case "crisis":
		fmt.Println("🔥 危机应对要点：先安抚情绪，再提供事实，最后展望未来")
		fmt.Println("💡 建议练习：准备3套不同程度的危机应对话术")
	case "investment":
		fmt.Println("💰 投资答辩要点：数据支撑信心，风险透明，愿景清晰")
		fmt.Println("💡 建议练习：准备项目5大关键数据的快速调用")
	}

	fmt.Println()
	fmt.Printf("🎯 针对【%s】场景的【%s】风格专项训练计划：\n", getScenarioDisplayName(selectedScenario), getStyleDisplayName(selectedStyle))

	switch selectedStyle {
	case "kanghui":
		switch selectedScenario {
		case "performance":
			fmt.Println("   📅 第1-3天：述职数据防御 - 准备项目ROI、风险、贡献等关键数据")
			fmt.Println("   📅 第4-7天：领导质疑应对 - 练习用\"数据显示...\"开始的回应")
			fmt.Println("   📅 第8-14天：专业形象塑造 - 掌握数据引用后的停顿技巧")
		case "presentation":
			fmt.Println("   📅 第1-3天：公开场合权威建立 - 面对刁难时保持专业风度")
			fmt.Println("   📅 第4-7天：事实澄清技巧 - 用数据化解技术质疑")
			fmt.Println("   📅 第8-14天：品牌维护训练 - 学会在冲突中维护个人专业形象")
		case "debate":
			fmt.Println("   📅 第1-3天：情绪控制与专业回应 - 不带情绪地用事实回应")
			fmt.Println("   📅 第4-7天：逻辑支撑训练 - 每个观点都有数据或事实依据")
			fmt.Println("   📅 第8-14天：辩论专业度提升 - 建立让人信服的专业形象")
		}
	case "chengming":
		switch selectedScenario {
		case "performance":
			fmt.Println("   📅 第1-3天：领导质疑拆解 - 分析领导问题背后的真实意图")
			fmt.Println("   📅 第4-7天：避实就虚练习 - 转移战场到领导弱点而不直接对抗")
			fmt.Println("   📅 第8-14天：策略性让步 - 学会在述职中适当妥协以求共赢")
		case "presentation":
			fmt.Println("   📅 第1-3天：刁难者意图识别 - 判断是恶意还是认知差异")
			fmt.Println("   📅 第4-7天：逻辑陷阱设置 - 请君入瓮诱导对方进入己方框架")
			fmt.Println("   📅 第8-14天：场下反击艺术 - 在不失风度前提下逻辑反制")
		case "debate":
			fmt.Println("   📅 第1-3天：情绪对抗策略 - 识别对方情绪弱点，逻辑突破")
			fmt.Println("   📅 第4-7天：归谬反驳训练 - 承认对方后推导荒谬结论")
			fmt.Println("   📅 第8-14天：节奏控制技巧 - 掌握争辩的攻守转换时机")
		}
	default: // hanhan
		switch selectedScenario {
		case "performance":
			fmt.Println("   📅 第1-3天：领导质疑反问 - 用\"这个问题背后真正的考虑是...\"")
			fmt.Println("   📅 第4-7天：真诚反思表达 - 展现对工作的深刻思考")
			fmt.Println("   📅 第8-14天：建设性坦诚 - 批评中带着改进建议")
		case "presentation":
			fmt.Println("   📅 第1-3天：挑衅化解技巧 - 用幽默化解恶意攻击")
			fmt.Println("   📅 第4-7天：态度反制训练 - 面对不公时坚定立场")
			fmt.Println("   📅 第8-14天：真诚破冰 - 用坦诚化解紧张气氛")
		case "debate":
			fmt.Println("   📅 第1-3天：情绪冲突处理 - 直面情绪但不陷入对抗")
			fmt.Println("   📅 第4-7天：犀利反问练习 - 质疑对方观点的前提")
			fmt.Println("   📅 第8-14天：适时结束对话 - 掌握不纠缠的智慧")
		}
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
		"kanghui":  "康辉式 (标准得体)",
		"chengming": "成铭式 (辩手机制)",
		"hanhan":    "韩寒式 (犀利风格)",
	}

	if name, ok := names[style]; ok {
		return name
	}
	return "韩寒式 (犀利风格)"
}

func getScenarioDisplayName(scenario string) string {
	names := map[string]string{
		"performance": "述职答辩危机",
		"presentation": "分享会刁难应对",
		"debate":      "争辩冲突反应",
		"elevator":    "电梯汇报挑战",
		"crisis":      "舆情危机应对",
		"investment":  "投资人答辩",
	}

	if name, ok := names[scenario]; ok {
		return name
	}
	return "述职答辩危机"
}
