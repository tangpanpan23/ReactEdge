package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"reactedge/internal/ai"
	"reactedge/internal/challenge"
	"reactedge/web"
)

func main() {
	fmt.Println("🎤 AI酷表达实验室 · 韩寒特训版 启动中...")

	// 初始化AI引擎
	hanAI := ai.NewHanStyleAI()
	fmt.Printf("✅ 韩寒表达引擎已加载，包含 %d 个表达模式\n", len(hanAI.GetExpressionPatterns()))

	// 初始化挑战管理器
	challengeManager := challenge.NewManager(hanAI)

	// 初始化Web服务器
	server := web.NewServer(challengeManager)

	// 启动服务器
	fmt.Println("🚀 服务器启动在 http://localhost:8080")
	fmt.Println("🎯 准备好开始你的3分钟表达挑战了吗？")

	log.Fatal(http.ListenAndServe(":8080", server.Router()))
}
