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
	fmt.Println("🎤 AI酷表达实验室 · 言刃 ReactEdge 启动中...")
	fmt.Println("   职场沟通的\"防弹衣\" - 述职答辩、分享会刁难、争辩冲突的快速反应训练")

	// 初始化AI引擎
	hanAI := ai.NewHanStyleAI()
	fmt.Printf("✅ 三大职场危机应对引擎已加载，包含 %d 个反应模式\n", len(hanAI.GetExpressionPatterns()))
	fmt.Println("   支持康辉式专业防御、成铭式逻辑反击、韩寒式态度反制")

	// 初始化挑战管理器
	challengeManager := challenge.NewManager(hanAI)

	// 初始化Web服务器
	server := web.NewServer(challengeManager)

	// 启动服务器
	fmt.Println("🚀 服务器启动在 http://localhost:8080")
	fmt.Println("🎯 准备好开始你的3分钟表达挑战了吗？")

	log.Fatal(http.ListenAndServe(":8080", server.Router()))
}
