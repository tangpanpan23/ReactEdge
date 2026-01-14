package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"reactedge/internal/ai"
)

func main() {
	fmt.Println("🎭 职场沟通风格演示系统 · 言刃 ReactEdge")
	fmt.Println("========================================")
	fmt.Println()

	// 初始化AI引擎
	hanAI := ai.NewHanStyleAI()
	fmt.Printf("✅ AI风格模仿引擎已加载，包含 %d 个表达模式\n", len(hanAI.GetExpressionPatterns()))
	fmt.Println("   支持康辉、董卿、韩寒、成铭四人风格")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// 第一步：选择名人风格
	fmt.Println("🎭 第一步：请选择你的目标表达风格")
	fmt.Println()
	fmt.Println("1️⃣  康辉（专业得体）- 沉稳权威，适合正式场合")
	fmt.Println("2️⃣  董卿（温婉大气）- 情感共鸣，适合沟通交流")
	fmt.Println("3️⃣  韩寒（犀利风格）- 反常规视角，适合辩论表达")
	fmt.Println("4️⃣  成铭（逻辑严谨）- 理性分析，适合策略破局")
	fmt.Println()

	selectedStyle := ""
	selectedStyleName := ""

	for {
		fmt.Print("请选择风格 (1-4): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			selectedStyle = "kanghui"
			selectedStyleName = "康辉（专业得体）"
			fmt.Println("✅ 已选择：康辉（专业得体）- 沉稳权威，适合正式场合")
		case "2":
			selectedStyle = "dongqing"
			selectedStyleName = "董卿（温婉大气）"
			fmt.Println("✅ 已选择：董卿（温婉大气）- 情感共鸣，适合沟通交流")
		case "3":
			selectedStyle = "hanhan"
			selectedStyleName = "韩寒（犀利风格）"
			fmt.Println("✅ 已选择：韩寒（犀利风格）- 反常规视角，适合辩论表达")
		case "4":
			selectedStyle = "chengming"
			selectedStyleName = "成铭（逻辑严谨）"
			fmt.Println("✅ 已选择：成铭（逻辑严谨）- 理性分析，适合策略破局")
		default:
			fmt.Println("❌ 无效选择，请输入1-4之间的数字")
			continue
		}
		break
	}

	fmt.Println()

	// 第二步：选择经典讲话内容
	fmt.Println("📚 第二步：请选择经典讲话内容参考")
	fmt.Println()

	classicContent := map[string][]string{
		"kanghui": {
			"《新闻联播》疫情报道（2020年）",
			"《新闻周刊》节目主持内容",
			"中央电视台大型晚会主持词",
		},
		"dongqing": {
			"《中国诗词大会》总决赛主持词",
			"《朗读者》节目串联词",
			"《故事里的中国》系列节目",
		},
		"hanhan": {
			"博客文章《一座城池》（完整版）",
			"演讲稿《我所理解的生活》",
			"微博经典长文（2010-2020年）",
		},
		"chengming": {
			"《奇葩说》经典辩论回合",
			"《超级演说家》演讲内容",
			"商业演讲和TED演讲",
		},
	}

	contentOptions := classicContent[selectedStyle]
	for i, content := range contentOptions {
		fmt.Printf("%d️⃣  %s\n", i+1, content)
	}
	fmt.Println()

	selectedContent := ""
	for {
		fmt.Print("请选择经典内容 (1-3): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input >= "1" && input <= "3" {
			idx := int(input[0] - '1')
			selectedContent = contentOptions[idx]
			fmt.Printf("✅ 已选择：%s\n", selectedContent)
			break
		} else {
			fmt.Println("❌ 无效选择，请输入1-3之间的数字")
		}
	}

	fmt.Println()

	// 第三步：输入职场问题
	fmt.Println("💼 第三步：请输入你的职场问题")
	fmt.Println()
	fmt.Println("例如：")
	fmt.Println("- \"领导问我这个项目的ROI为什么这么低？\"")
	fmt.Println("- \"分享会上有人质疑我的技术方案不可行\"")
	fmt.Println("- \"同事说我这个想法太不切实际了\"")
	fmt.Println()

	fmt.Print("请输入你的问题：")
	userQuestion, _ := reader.ReadString('\n')
	userQuestion = strings.TrimSpace(userQuestion)

	if userQuestion == "" {
		userQuestion = "领导问我这个项目的ROI为什么这么低？"
		fmt.Printf("使用示例问题：%s\n", userQuestion)
	}

	fmt.Println()

	// 第四步：获得风格化回答
	fmt.Println("🤖 第四步：生成风格化回答")
	fmt.Println()

	fmt.Printf("🎯 基于【%s】风格，参考【%s】\n", selectedStyleName, selectedContent)
	fmt.Printf("❓ 问题：%s\n", userQuestion)
	fmt.Println()

	// 生成风格化回答
	response := hanAI.GenerateStyleResponse(selectedStyle, userQuestion, selectedContent)

	fmt.Printf("💬 %s式回答：\n", selectedStyleName)
	fmt.Println()
	fmt.Println(response)
	fmt.Println()

	// 提供一些使用建议
	fmt.Println("💡 风格解析：")
	switch selectedStyle {
	case "kanghui":
		fmt.Println("• 专业得体：用数据和事实支撑观点，展现权威性")
		fmt.Println("• 适用场合：正式汇报、述职答辩、技术讨论")
	case "dongqing":
		fmt.Println("• 温婉大气：注重情感共鸣，温和有礼的沟通方式")
		fmt.Println("• 适用场合：跨部门协调、客户沟通、团队建设")
	case "hanhan":
		fmt.Println("• 犀利直接：直言不讳，反常规视角，追求观点冲击力")
		fmt.Println("• 适用场合：应对质疑、辩论冲突、观点交锋")
	case "chengming":
		fmt.Println("• 逻辑严谨：层层递进，策略性思维，掌控局面")
		fmt.Println("• 适用场合：方案辩护、利益谈判、危机应对")
	}

	fmt.Println()
	fmt.Println("🎉 演示完成！")
	fmt.Println("你可以继续尝试不同风格和问题，体验各种沟通方式的效果。")
	fmt.Println()
	fmt.Println("感谢体验 职场沟通风格演示系统 · 言刃 ReactEdge！")
}

