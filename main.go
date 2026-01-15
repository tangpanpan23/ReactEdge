package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"reactedge/config"
	"reactedge/internal/ai"
	aiPkg "reactedge/pkg/ai"
	"reactedge/web"
)

func main() {
	fmt.Println("🎭 职场沟通风格演示系统 · 言刃 ReactEdge 启动中...")
	fmt.Println("   看康辉、董卿、韩寒、成铭如何回答你的职场问题！")

	// 加载应用配置
	appConfig, err := config.Load()
	if err != nil {
		fmt.Printf("❌ 配置加载失败: %v\n", err)
		fmt.Println("⚠️ 将使用默认配置继续运行")
		// 不设置appConfig，使用nil，让各个组件使用默认值
		appConfig = nil
	} else {
		fmt.Printf("✅ 应用配置加载成功，端口: %s，AI模式: %s\n", appConfig.Server.Port, appConfig.AI.Mode)
	}

	// 初始化AI引擎
	hanAI := ai.NewHanStyleAI()
	fmt.Printf("✅ AI风格模仿引擎已加载，包含 %d 个表达模式\n", len(hanAI.GetExpressionPatterns()))
	fmt.Println("   支持康辉、董卿、韩寒、成铭四人风格")

	// 初始化AI管理器
	aiManager, err := aiPkg.NewManager("config/ai.yaml")
	if err != nil {
		fmt.Printf("❌ AI服务初始化失败: %v\n", err)
		fmt.Println("⚠️ 将使用本地模拟回答")
		aiManager = nil
	} else {
		fmt.Println("✅ AI服务管理器初始化成功")
	}

	// 初始化Web服务器
	server := web.NewServer(hanAI, aiManager, appConfig)

	// 如果配置为空，使用默认配置
	if appConfig == nil {
		appConfig = config.GetDefaultConfig()
	}

	// 创建HTTP服务器
	addr := fmt.Sprintf("%s:%s", appConfig.Server.Host, appConfig.Server.Port)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      server.Router(),
		ReadTimeout:  time.Duration(appConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(appConfig.Server.WriteTimeout) * time.Second,
	}

	// 启动服务器
	fmt.Printf("🚀 服务器启动在 http://%s\n", addr)
	fmt.Println("🎯 开始你的职场沟通风格探索之旅！")
	fmt.Printf("   AI模式: %s，使用qwen3-max模型进行推理\n", appConfig.AI.Mode)

	if appConfig.Server.TLSEnabled && appConfig.Server.TLSCertFile != "" && appConfig.Server.TLSKeyFile != "" {
		fmt.Println("🔒 HTTPS模式已启用")
		log.Fatal(httpServer.ListenAndServeTLS(appConfig.Server.TLSCertFile, appConfig.Server.TLSKeyFile))
	} else {
		log.Fatal(httpServer.ListenAndServe())
	}
}
