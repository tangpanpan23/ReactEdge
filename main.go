package main

import (
	"fmt"
	"log"
	"net/http"

	"reactedge/internal/ai"
	"reactedge/web"
)

func main() {
	fmt.Println("🎭 职场沟通风格演示系统 · 言刃 ReactEdge 启动中...")
	fmt.Println("   看康辉、董卿、韩寒、成铭如何回答你的职场问题！")

	// 初始化AI引擎
	hanAI := ai.NewHanStyleAI()
	fmt.Printf("✅ AI风格模仿引擎已加载，包含 %d 个表达模式\n", len(hanAI.GetExpressionPatterns()))
	fmt.Println("   支持康辉、董卿、韩寒、成铭四人风格")

	// 初始化Web服务器
	server := web.NewServer(hanAI)

	// 启动服务器
	fmt.Println("🚀 服务器启动在 http://localhost:8080")
	fmt.Println("🎯 开始你的职场沟通风格探索之旅！")

	log.Fatal(http.ListenAndServe(":8080", server.Router()))
}
