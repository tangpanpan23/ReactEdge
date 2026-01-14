#!/bin/bash

# AI酷表达实验室 · 韩寒特训版 - 运行脚本

echo "🎤 AI酷表达实验室 · 言刃 ReactEdge"
echo "=================================="

# 检查Go是否安装
if ! command -v go &> /dev/null; then
    echo "❌ 错误：Go未安装。请先安装Go 1.19+"
    echo "   下载地址：https://golang.org/dl/"
    exit 1
fi

echo "✅ Go版本：$(go version)"

# 设置环境变量
export GO111MODULE=on
export GOPROXY=https://goproxy.cn,direct

# 清理之前的构建
echo "🧹 清理之前的构建..."
rm -f reactedge demo

# 选择运行模式
echo ""
echo "请选择运行模式："
echo "1) 🌐 Web服务器模式 (完整功能)"
echo "2) 💻 命令行演示模式 (快速体验)"
echo "3) 🧪 运行测试"
echo ""
read -p "请输入选择 (1-3): " choice

case $choice in
    1)
        echo "🚀 启动Web服务器模式..."
        echo "   服务器将在 http://localhost:8080 启动"
        echo "   按 Ctrl+C 停止服务器"
        echo ""

        # 尝试运行Web服务器
        if go run . 2>/dev/null; then
            echo "✅ 服务器启动成功"
        else
            echo "⚠️  Web服务器启动失败，可能由于权限问题"
            echo "   建议运行命令行演示模式体验功能"
            echo ""
            read -p "是否运行演示模式？(y/N): " demo_choice
            if [[ $demo_choice =~ ^[Yy]$ ]]; then
                choice=2
            else
                exit 1
            fi
        fi
        ;;
    2)
        echo "💻 启动命令行演示模式..."
        echo "   这是一个完整的3分钟表达挑战体验"
        echo ""

        # 运行演示
        if go run cmd/demo.go; then
            echo ""
            echo "✅ 演示完成！"
        else
            echo "❌ 演示运行失败"
            exit 1
        fi
        ;;
    3)
        echo "🧪 运行测试..."
        echo ""

        # 运行测试
        go test ./... -v
        ;;
    *)
        echo "❌ 无效选择"
        exit 1
        ;;
esac

echo ""
echo "感谢使用 AI酷表达实验室！"
echo "有任何问题或建议欢迎反馈。"
