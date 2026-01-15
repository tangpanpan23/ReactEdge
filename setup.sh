#!/bin/bash

# ReactEdge 项目配置设置脚本
# 帮助用户快速配置项目

echo "🎭 ReactEdge 项目配置设置"
echo "=========================="

# 检查配置文件是否存在
if [ ! -f "config/app.yaml" ]; then
    echo "📋 复制应用配置文件..."
    cp config/app.yaml.example config/app.yaml
    echo "✅ config/app.yaml 已创建"
else
    echo "ℹ️  config/app.yaml 已存在"
fi

if [ ! -f "config/ai.yaml" ]; then
    echo "🤖 复制AI配置文件..."
    cp config/ai.yaml.example config/ai.yaml
    echo "✅ config/ai.yaml 已创建"
else
    echo "ℹ️  config/ai.yaml 已存在"
fi

echo ""
echo "📝 配置文件说明："
echo "=================="
echo ""
echo "1. 应用配置 (config/app.yaml):"
echo "   - 服务器端口、超时设置"
echo "   - AI服务参数、缓存配置"
echo "   - 日志、监控、环境设置"
echo ""
echo "2. AI配置 (config/ai.yaml):"
echo "   - TAL AI服务配置 (需要设置AppID和AppKey)"
echo "   - OpenAI、Claude等第三方AI配置"
echo "   - 模型映射和参数设置"
echo ""
echo "⚠️  重要提醒："
echo "=============="
echo "1. ⚠️  配置文件已被添加到 .gitignore，不会提交到版本控制"
echo "2. 🔒 请勿在配置文件中填写真实API密钥"
echo "3. 🌍 建议使用环境变量设置敏感信息："
echo "   - TAL_MLOPS_APP_ID=your-app-id"
echo "   - TAL_MLOPS_APP_KEY=your-app-key"
echo "   - OPENAI_API_KEY=your-openai-key"
echo ""
echo "4. 📝 示例配置已设置合理的默认值，生产环境请根据需要调整"
echo "5. 🔄 如果需要重新配置，请删除现有配置文件后重新运行此脚本"
echo ""

# 检查是否需要编辑配置文件
read -p "是否现在编辑配置文件？(y/N): " edit_config
if [[ $edit_config =~ ^[Yy]$ ]]; then
    echo "选择要编辑的配置文件："
    echo "1) 应用配置 (config/app.yaml)"
    echo "2) AI配置 (config/ai.yaml)"
    echo "3) 两个都编辑"
    read -p "请选择 (1-3): " choice

    case $choice in
        1)
            ${EDITOR:-nano} config/app.yaml
            ;;
        2)
            ${EDITOR:-nano} config/ai.yaml
            ;;
        3)
            ${EDITOR:-nano} config/app.yaml
            ${EDITOR:-nano} config/ai.yaml
            ;;
        *)
            echo "无效选择，跳过编辑"
            ;;
    esac
fi

echo ""
echo "✅ 配置设置完成！"
echo ""
echo "🚀 运行项目："
echo "   ./run.sh              # 使用运行脚本"
echo "   go run .             # 直接运行"
echo "   SERVER_PORT=6001 go run .  # 指定端口"
echo ""
echo "📖 访问地址："
echo "   http://localhost:6000  # 默认端口"
echo "   http://localhost:6000/demo  # 演示页面"
