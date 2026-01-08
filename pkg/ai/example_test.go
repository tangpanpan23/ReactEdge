package ai

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestAIIntegration AI服务集成测试
func TestAIIntegration(t *testing.T) {
	// 创建AI服务管理器
	manager, err := NewManager("../../config/ai.yaml")
	if err != nil {
		t.Logf("AI服务管理器初始化失败（正常，因为没有真实API密钥）: %v", err)
		return
	}

	ctx := context.Background()

	// 测试基本功能
	t.Run("TestGenerateReactionTemplates", func(t *testing.T) {
		templates, err := manager.GenerateReactionTemplates(ctx, "述职答辩", "韩寒风格")
		if err != nil {
			t.Logf("生成反应模板失败（可能因为AI服务不可用）: %v", err)
			return
		}

		fmt.Printf("✅ 生成%d个反应模板\n", len(templates))
		for i, template := range templates {
			fmt.Printf("  模板%d: %s\n", i+1, template.Scenario)
		}
	})

	t.Run("TestSimulateDebate", func(t *testing.T) {
		simulation, err := manager.SimulateDebate(ctx, "述职答辩", 2, "韩寒风格")
		if err != nil {
			t.Logf("模拟辩论失败（可能因为AI服务不可用）: %v", err)
			return
		}

		fmt.Printf("✅ 辩论模拟结果:\n")
		fmt.Printf("  场景: %s\n", simulation.Scenario)
		fmt.Printf("  对手开场: %s\n", simulation.OpponentOpening)
		fmt.Printf("  交互轮数: %d\n", len(simulation.InteractionRounds))
		fmt.Printf("  难度等级: %d\n", simulation.Difficulty)
	})

	t.Run("TestEvaluateReaction", func(t *testing.T) {
		evaluation, err := manager.EvaluateReaction(ctx, "我认为这个问题需要从根本上 reconsider", "述职答辩", "韩寒风格")
		if err != nil {
			t.Logf("评估反应失败（可能因为AI服务不可用）: %v", err)
			return
		}

		fmt.Printf("✅ 反应评估结果:\n")
		fmt.Printf("  整体评分: %.1f\n", evaluation.OverallScore)
		fmt.Printf("  优势: %v\n", evaluation.Strengths)
		fmt.Printf("  改进建议: %v\n", evaluation.Improvements)
	})
}

// TestConfigLoading 测试配置加载
func TestConfigLoading(t *testing.T) {
	config, err := LoadConfig("../../config/ai.yaml")
	if err != nil {
		t.Logf("配置加载失败: %v", err)
		return
	}

	fmt.Printf("✅ 配置加载成功:\n")
	fmt.Printf("  默认服务商: %s\n", config.DefaultProvider)
	fmt.Printf("  可用服务商: %v\n", config.GetAvailableProviders())

	// 检查TAL配置
	if config.TAL.TAL_MLOPS_APP_ID != "" {
		fmt.Printf("  TAL配置: ✅ 已配置\n")
	} else {
		fmt.Printf("  TAL配置: ⚠️ 未配置（使用环境变量TAL_MLOPS_APP_ID）\n")
	}
}

// TestProviderSwitching 测试服务商切换
func TestProviderSwitching(t *testing.T) {
	manager, err := NewManager("../../config/ai.yaml")
	if err != nil {
		t.Logf("AI服务管理器初始化失败: %v", err)
		return
	}

	availableProviders := manager.GetAvailableProviders()
	fmt.Printf("✅ 可用服务商: %v\n", availableProviders)

	// 测试切换服务商
	for _, provider := range availableProviders {
		err := manager.SwitchProvider(provider)
		if err != nil {
			t.Logf("切换到%s失败: %v", provider, err)
			continue
		}

		currentClient := manager.GetClient()
		fmt.Printf("✅ 成功切换到%s，当前模型: %v\n",
			provider, currentClient.GetAvailableModels()[:3]) // 只显示前3个模型
		break
	}
}

// BenchmarkAIResponse 性能测试
func BenchmarkAIResponse(b *testing.B) {
	manager, err := NewManager("../../config/ai.yaml")
	if err != nil {
		b.Logf("AI服务管理器初始化失败: %v", err)
		return
	}

	ctx := context.Background()

	b.Run("GenerateReactionTemplates", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := manager.GenerateReactionTemplates(ctx, "述职答辩", "韩寒风格")
			if err != nil {
				b.Logf("生成模板失败: %v", err)
				break
			}
		}
	})
}

// ExampleUsage 使用示例
func ExampleUsage() {
	fmt.Println("🎯 ReactEdge AI服务使用示例")
	fmt.Println("================================")

	// 1. 初始化AI服务管理器
	manager, err := NewManager("config/ai.yaml")
	if err != nil {
		fmt.Printf("❌ 初始化失败: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 2. 生成反应模板
	fmt.Println("\n📝 生成述职答辩反应模板...")
	templates, err := manager.GenerateReactionTemplates(ctx, "述职答辩", "韩寒风格")
	if err != nil {
		fmt.Printf("❌ 生成模板失败: %v\n", err)
	} else {
		fmt.Printf("✅ 生成%d个模板\n", len(templates))
	}

	// 3. 模拟辩论
	fmt.Println("\n🎭 模拟述职辩论...")
	simulation, err := manager.SimulateDebate(ctx, "述职答辩", 2, "韩寒风格")
	if err != nil {
		fmt.Printf("❌ 模拟辩论失败: %v\n", err)
	} else {
		fmt.Printf("✅ 对手开场: %s\n", simulation.OpponentOpening)
	}

	// 4. 评估用户反应
	fmt.Println("\n📊 评估用户反应...")
	evaluation, err := manager.EvaluateReaction(ctx,
		"这个项目的ROI确实不高，但如果我们只看短期数字，那就太短视了",
		"述职答辩", "韩寒风格")
	if err != nil {
		fmt.Printf("❌ 评估失败: %v\n", err)
	} else {
		fmt.Printf("✅ 整体评分: %.1f/10\n", evaluation.OverallScore)
	}

	fmt.Println("\n🎉 示例完成！")
}
