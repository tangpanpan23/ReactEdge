# ReactEdge AI服务集成说明

## 🎯 集成概述

ReactEdge项目已成功集成ExploraPal的AI服务架构，支持多种AI服务商和模型，为职场沟通训练提供强大的AI能力。

## 🏗️ 架构设计

### AI服务包结构 (`pkg/ai/`)
```
pkg/ai/
├── config.go          # 配置管理
├── client.go          # TAL AI客户端 (主客户端)
├── openai_client.go   # OpenAI客户端
├── other_clients.go   # Claude/Azure/Baidu客户端
├── manager.go         # AI服务管理器
├── types.go           # 数据结构定义
└── example_test.go    # 测试和示例
```

### 支持的AI服务商
- **TAL** (默认) - 内部AI服务，支持Qwen系列模型
- **OpenAI** - GPT-4, GPT-3.5-turbo等
- **Claude** - Anthropic Claude模型
- **Azure** - Azure OpenAI服务
- **Baidu** - 百度AI服务

## 🤖 支持的AI模型

### TAL内部服务 (推荐)
| 模型 | 用途 | 特点 |
|------|------|------|
| `qwen3-vl-plus` | 图像分析 | 视觉理解，支持思考模式 |
| `qwen-flash` | 文本生成 | 快速响应，适合对话 |
| `qwen3-max` | 复杂推理 | 高级推理，专业分析 |
| `qwen3-omni-flash` | 语音交互 | 多模态，语音处理 |
| `doubao-seedance-1.0-lite-i2v` | 视频生成 | 图像到视频转换 |

### OpenAI服务
- `gpt-4o` - 图像分析和复杂任务
- `gpt-4` - 高级文本生成
- `gpt-3.5-turbo` - 快速文本处理

### Claude服务
- `claude-3-opus-20240229` - 高级分析
- `claude-3-haiku-20240307` - 快速生成

## ⚙️ 配置方法

### 1. 配置文件 (`config/ai.yaml`)
```yaml
# 默认服务商
defaultProvider: "tal"  # tal, openai, claude, azure, baidu

# TAL配置 (推荐)
tal:
  talMLOpsAppId: "your-app-id"      # 或环境变量 TAL_MLOPS_APP_ID
  talMLOpsAppKey: "your-app-key"    # 或环境变量 TAL_MLOPS_APP_KEY
  baseURL: "http://ai-service.tal.com/openai-compatible/v1"
  timeout: 30
  maxTokens: 2000
  temperature: 0.7

# OpenAI配置
openai:
  apiKey: "sk-..."                   # 或环境变量 OPENAI_API_KEY
  baseURL: "https://api.openai.com/v1"
  timeout: 30
  maxTokens: 2000
  temperature: 0.7
```

### 2. 环境变量
```bash
# TAL服务
export TAL_MLOPS_APP_ID="your-app-id"
export TAL_MLOPS_APP_KEY="your-app-key"

# OpenAI
export OPENAI_API_KEY="sk-your-key"

# Claude
export ANTHROPIC_API_KEY="sk-ant-..."

# Azure
export AZURE_OPENAI_API_KEY="your-azure-key"
export AZURE_OPENAI_ENDPOINT="https://your-resource.openai.azure.com/"

# 百度
export BAIDU_API_KEY="your-baidu-key"
export BAIDU_SECRET_KEY="your-secret-key"
```

## 🔧 使用方法

### 初始化AI管理器
```go
import "reactedge/pkg/ai"

// 创建AI服务管理器
manager, err := ai.NewManager("config/ai.yaml")
if err != nil {
    log.Fatal("AI服务初始化失败:", err)
}

// 获取默认客户端
client := manager.GetClient()

// 切换服务商
err = manager.SwitchProvider(ai.ProviderOpenAI)
```

### 核心AI功能
```go
ctx := context.Background()

// 1. 生成反应模板
templates, err := client.GenerateReactionTemplates(ctx, "述职答辩", "韩寒风格")

// 2. 分析表达风格
analysis, err := client.AnalyzeExpressionStyle(ctx, "韩寒", "样本文本...")

// 3. 模拟辩论
simulation, err := client.SimulateDebate(ctx, "述职答辩", 2, "韩寒风格")

// 4. 评估反应
evaluation, err := client.EvaluateReaction(ctx, "用户反应...", "述职答辩", "韩寒风格")

// 5. 个性化训练
training, err := manager.GeneratePersonalizedTraining(ctx, userProfile, level)
```

## 🎯 ReactEdge集成

### HanStyleAI增强
原有的`HanStyleAI`已集成新的AI服务：

```go
hanAI := ai.NewHanStyleAI()

// 新的AI增强方法
templates, _ := hanAI.GenerateReactionTemplatesAI(ctx, "述职答辩", "韩寒风格")
analysis, _ := hanAI.AnalyzeExpressionStyleAI(ctx, "韩寒", "样本文本...")
simulation, _ := hanAI.SimulateDebateAI(ctx, "述职答辩", 2, "韩寒风格")
evaluation, _ := hanAI.EvaluateReactionAI(ctx, "用户反应...", "述职答辩", "韩寒风格")
```

### 自动降级
- 当AI服务不可用时，自动使用模拟数据
- 保证系统在任何情况下都能正常运行
- 提供详细的日志信息

## 📊 功能特性

### 职场沟通训练
- ✅ **述职答辩** - 应对领导质疑的专业反应
- ✅ **分享会刁难** - 处理公开场合的挑战性问题
- ✅ **争辩冲突** - 日常沟通中的立场维护

### AI能力
- 🔍 **图像分析** - 分析表情、姿势等非语言信号
- 📝 **文本生成** - 生成个性化反应模板
- 🎭 **风格分析** - 深度解析表达风格特征
- ⚔️ **辩论模拟** - 真实场景的AI对手
- 📊 **反应评估** - 多维度反应质量分析
- 🎯 **个性化训练** - 基于用户特征的定制计划

## 🚀 性能优化

### 智能服务商选择
- 根据任务类型自动选择最适合的模型
- 支持服务商切换和负载均衡
- 提供故障转移机制

### 缓存和优化
- 响应结果缓存
- 并发请求处理
- 超时和重试机制

## 🔍 监控和调试

### 日志系统
```
✅ AI服务状态: 已连接 TAL 服务商
✅ TAL配置: 已配置
✅ 可用模型: 7 个
```

### 调试信息
- AI请求参数记录
- 响应时间统计
- 错误详情输出

## 🧪 测试

运行AI服务测试：
```bash
cd pkg/ai
go test -v -run TestAIIntegration
```

运行使用示例：
```go
// 在 example_test.go 中
ExampleUsage()
```

## 📈 扩展计划

### 新增服务商
- [ ] 通义千问官方API
- [ ] 智谱GLM
- [ ] 月之暗面Kimi
- [ ] 腾讯混元

### 功能增强
- [ ] 语音识别和合成
- [ ] 视频内容分析
- [ ] 多模态交互
- [ ] 实时对话训练

## 🔗 相关链接

- [ExploraPal AI架构](https://github.com/explorapal) - 参考的AI服务架构
- [TAL AI服务](http://ai-service.tal.com) - 内部AI服务平台
- [ReactEdge主项目](https://github.com/reactedge) - 主项目仓库

---

## 🎊 总结

通过集成ExploraPal的AI服务架构，ReactEdge现在具备了：

1. **多服务商支持** - TAL、OpenAI、Claude等主流AI服务
2. **丰富模型选择** - 从轻量级到专业级的完整模型系列
3. **智能降级机制** - 确保服务稳定性和可用性
4. **职场专用优化** - 专门针对沟通训练场景优化
5. **易于扩展** - 模块化设计，易于添加新的AI能力

这为ReactEdge提供了强大的AI驱动能力，能够为用户提供更加智能、个性化的职场沟通训练体验！ 🚀✨
