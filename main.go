package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
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

	// 创建HTTP服务器，并自动处理端口冲突
	addr, httpServer := createHTTPServer(appConfig, server)

	// 启动服务器
	fmt.Printf("🚀 服务器启动在 http://%s\n", addr)
	fmt.Println("🎯 开始你的职场沟通风格探索之旅！")
	if appConfig.AI.Mode == "internal" {
		fmt.Printf("   AI模式: %s，使用TAL(deepseek-reasoner)进行推理\n", appConfig.AI.Mode)
	} else {
		fmt.Printf("   AI模式: %s，使用星火AI(spark-x)进行推理\n", appConfig.AI.Mode)
	}
	fmt.Printf("   AI交互超时: %d秒\n", appConfig.AI.InteractionTimeout)

	if appConfig.Server.TLSEnabled && appConfig.Server.TLSCertFile != "" && appConfig.Server.TLSKeyFile != "" {
		fmt.Println("🔒 HTTPS模式已启用")
		fmt.Printf("📡 尝试启动HTTPS服务器在: %s\n", addr)
		log.Fatal(httpServer.ListenAndServeTLS(appConfig.Server.TLSCertFile, appConfig.Server.TLSKeyFile))
	} else {
		fmt.Printf("📡 尝试启动HTTP服务器在: %s\n", addr)
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Printf("❌ 服务器启动失败: %v\n", err)
			log.Fatal(err)
		}
	}
}

// createHTTPServer 创建HTTP服务器，自动处理端口冲突
func createHTTPServer(appConfig *config.Config, server *web.Server) (string, *http.Server) {
	basePort, _ := strconv.Atoi(appConfig.Server.Port)

	// 尝试从基础端口开始，逐步增加直到成功绑定
	for port := basePort; port < basePort+100; port++ {
		addr := fmt.Sprintf("%s:%d", appConfig.Server.Host, port)

		// 直接尝试创建HTTP服务器并监听，如果成功则返回
		httpServer := &http.Server{
			Addr:         addr,
			Handler:      server.Router(),
			ReadTimeout:  time.Duration(appConfig.Server.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(appConfig.Server.WriteTimeout) * time.Second,
		}

		// 尝试监听端口
		listener, err := net.Listen("tcp", addr)
		if err == nil {
			// 成功绑定，关闭临时监听器让http.Server使用
			listener.Close()
			return addr, httpServer
		}

		if port == basePort {
			fmt.Printf("⚠️ 端口 %d 被占用，尝试查找可用端口...\n", port)
		}
	}

	// 如果没有找到可用端口，使用系统分配的随机端口
	listener, err := net.Listen("tcp", appConfig.Server.Host+":0")
	if err != nil {
		log.Fatalf("无法创建监听器: %v", err)
	}

	actualAddr := listener.Addr().String()
	listener.Close() // 关闭临时监听器，http.Server会重新创建

	fmt.Printf("✅ 使用随机可用端口: %s\n", actualAddr)

	httpServer := &http.Server{
		Addr:         actualAddr,
		Handler:      server.Router(),
		ReadTimeout:  time.Duration(appConfig.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(appConfig.Server.WriteTimeout) * time.Second,
	}

	return actualAddr, httpServer
}
