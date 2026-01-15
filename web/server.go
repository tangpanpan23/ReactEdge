package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"reactedge/config"
	"reactedge/internal/ai"
	aiPkg "reactedge/pkg/ai"
	"github.com/gorilla/websocket"
)

// Server Web服务器
type Server struct {
	aiEngine *ai.HanStyleAI
	aiManager *aiPkg.Manager
	config   *config.Config
	router   *http.ServeMux
	upgrader websocket.Upgrader
}

// NewServer 创建Web服务器
func NewServer(aiEngine *ai.HanStyleAI, aiManager *aiPkg.Manager, config *config.Config) *Server {
	server := &Server{
		aiEngine: aiEngine,
		aiManager: aiManager,
		config:   config,
		router:   http.NewServeMux(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// 允许所有来源，生产环境中应该更严格
				return true
			},
		},
	}

	server.setupRoutes()

	return server
}

// Router 获取路由器
func (s *Server) Router() *http.ServeMux {
	return s.router
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	s.router.HandleFunc("/", s.handleHome)
	s.router.HandleFunc("/demo", s.handleDemo)
	s.router.HandleFunc("/generate", s.handleGenerate)
	s.router.HandleFunc("/ws", s.handleWebSocket)
}

// handleHome 首页
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>职场沟通风格演示系统 · 言刃 ReactEdge</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .container { text-align: center; }
        .button { background: #007bff; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; margin: 10px; }
        .button:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎭 职场沟通风格演示系统</h1>
        <h2>言刃 ReactEdge</h2>
        <p>看康辉、董卿、韩寒、成铭如何回答你的职场问题！</p>
        <a href="/demo"><button class="button">开始演示</button></a>
    </div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// handleDemo 演示页面
func (s *Server) handleDemo(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>职场沟通演示</title>
    <style>
        body {
            font-family: 'Microsoft YaHei', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            color: #333;
        }

        .container {
            background: rgba(255, 255, 255, 0.95);
            border-radius: 15px;
            padding: 30px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            backdrop-filter: blur(10px);
        }

        h1 {
            text-align: center;
            color: #2c3e50;
            margin-bottom: 30px;
            font-size: 2.5em;
            text-shadow: 2px 2px 4px rgba(0,0,0,0.1);
        }

        .step {
            margin: 25px 0;
            padding: 25px;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            background: #fff;
            transition: all 0.3s ease;
        }

        .step:hover {
            border-color: #667eea;
            box-shadow: 0 5px 15px rgba(102, 126, 234, 0.2);
        }

        .step h3 {
            color: #495057;
            margin-bottom: 15px;
            font-size: 1.3em;
        }

        .form-group { margin: 15px 0; }
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #495057;
        }

        select, textarea {
            width: 100%;
            padding: 12px;
            border: 2px solid #e9ecef;
            border-radius: 8px;
            font-size: 14px;
            transition: border-color 0.3s ease;
        }

        select:focus, textarea:focus {
            outline: none;
            border-color: #667eea;
            box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
        }

        textarea {
            resize: vertical;
            min-height: 80px;
            font-family: inherit;
        }

        .button {
            background: linear-gradient(135deg, #667eea, #764ba2);
            color: white;
            padding: 12px 30px;
            border: none;
            border-radius: 25px;
            cursor: pointer;
            font-size: 16px;
            font-weight: 600;
            transition: all 0.3s ease;
            box-shadow: 0 4px 15px rgba(102, 126, 234, 0.4);
        }

        .button:hover:not(:disabled) {
            transform: translateY(-2px);
            box-shadow: 0 6px 20px rgba(102, 126, 234, 0.6);
        }

        .button:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none;
        }

        .result {
            margin-top: 30px;
            padding: 25px;
            background: linear-gradient(135deg, #f8f9fa, #e9ecef);
            border-radius: 10px;
            border-left: 5px solid #667eea;
            animation: fadeIn 0.5s ease-in;
        }

        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(20px); }
            to { opacity: 1; transform: translateY(0); }
        }

        .result h3 {
            color: #495057;
            margin-bottom: 15px;
            display: flex;
            align-items: center;
            gap: 10px;
        }

        #response {
            background: #fff;
            padding: 20px;
            border-radius: 8px;
            border: 1px solid #e9ecef;
            margin: 15px 0;
            line-height: 1.8;
            font-size: 16px;
            color: #2c3e50;
        }

        #status {
            font-size: 14px;
            font-weight: 500;
            padding: 8px 0;
            border-radius: 4px;
        }

        .loading {
            color: #007bff !important;
            animation: pulse 2s infinite;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.7; }
        }

        .success { color: #28a745 !important; }
        .error { color: #dc3545 !important; }

        .response-content p {
            margin: 12px 0;
            line-height: 1.8;
            text-align: justify;
        }

        .response-content p:first-child {
            text-indent: 0;
            font-weight: 500;
            color: #2c3e50;
        }

        .connection-status {
            position: fixed;
            top: 10px;
            right: 10px;
            padding: 8px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 500;
            z-index: 1000;
            display: flex;
            align-items: center;
            gap: 5px;
        }

        .connection-status.connected {
            background: #d4edda;
            color: #155724;
            border: 1px solid #c3e6cb;
        }

        .connection-status.connecting {
            background: #fff3cd;
            color: #856404;
            border: 1px solid #ffeaa7;
        }

        .connection-status.disconnected {
            background: #f8d7da;
            color: #721c24;
            border: 1px solid #f5c6cb;
        }

        .help-text {
            font-size: 12px;
            color: #6c757d;
            margin-top: 5px;
            font-style: italic;
        }

        @media (max-width: 768px) {
            body { padding: 10px; }
            .container { padding: 20px; }
            h1 { font-size: 2em; }
            .step { padding: 20px; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div id="connectionStatus" class="connection-status">
            <span id="statusIndicator">🔴</span>
            <span id="statusText">连接中...</span>
        </div>
        <h1>🎯 职场沟通风格演示</h1>

    <div class="step">
        <h3>第一步：选择名人风格</h3>
        <div class="form-group">
            <label for="style">选择风格：</label>
            <select id="style">
                <option value="kanghui">康辉（专业得体）- 沉稳权威，适合正式场合</option>
                <option value="dongqing">董卿（温婉大气）- 情感共鸣，适合沟通交流</option>
                <option value="hanhan">韩寒（犀利风格）- 反常规视角，适合辩论表达</option>
                <option value="chengming">成铭（逻辑严谨）- 理性分析，适合策略破局</option>
            </select>
        </div>
    </div>

    <div class="step">
        <h3>第二步：选择经典讲话内容</h3>
        <div class="form-group">
            <label for="content">选择经典内容：</label>
            <select id="content">
                <option value="news">《新闻联播》疫情报道（康辉）</option>
                <option value="poetry">《中国诗词大会》总决赛主持词（董卿）</option>
                <option value="blog">博客文章《一座城池》（韩寒）</option>
                <option value="debate">《奇葩说》经典辩论回合（成铭）</option>
            </select>
        </div>
    </div>

    <div class="step">
        <h3>第三步：输入职场问题</h3>
        <div class="form-group">
            <label for="question">输入你的职场问题：</label>
            <textarea id="question" rows="3" placeholder="例如：领导问我这个项目的ROI为什么这么低？如何处理团队冲突？项目延期了怎么汇报？"></textarea>
            <div class="help-text">💡 提示：按 Enter 键快速生成回答，Shift+Enter 换行</div>
        </div>
        <div style="display: flex; gap: 10px; align-items: center;">
            <button class="button" id="generateBtn" onclick="generateResponse()">🤖 生成AI回答</button>
            <button class="button" id="cancelBtn" onclick="cancelRequest()" style="display: none; background: #dc3545;" disabled>⏹️ 取消请求</button>
        </div>
    </div>

    <div id="result" class="result" style="display: none;">
        <h3>🤖 AI生成回答</h3>
        <div id="response" class="response-content" style="font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;"></div>
        <div id="status" style="margin-top: 10px; font-size: 14px; color: #666;"></div>
    </div>

    <script>
        let websocket = null;
        let isConnected = false;
        let currentRequestId = null; // 跟踪当前请求
        let reconnectAttempts = 0;
        const maxReconnectAttempts = 5;

        function connectWebSocket() {
            if (websocket && websocket.readyState === WebSocket.OPEN) {
                return; // 已经连接
            }

            if (reconnectAttempts >= maxReconnectAttempts) {
                updateConnectionStatus('disconnected', '🔴', '连接失败');
                console.log('达到最大重连次数，停止重连');
                return;
            }

            updateConnectionStatus('connecting', '🟡', '连接中...');

            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws';

            try {
                websocket = new WebSocket(wsUrl);
                reconnectAttempts++;

                websocket.onopen = function(event) {
                    isConnected = true;
                    reconnectAttempts = 0; // 重置重连计数
                    updateConnectionStatus('connected', '🟢', '已连接');
                    console.log('WebSocket连接已建立');
                };

                // 处理pong响应
                websocket.onpong = function(event) {
                    console.log('收到pong响应');
                };

                websocket.onmessage = function(event) {
                    try {
                        const message = JSON.parse(event.data);
                        handleWebSocketMessage(message);
                    } catch (e) {
                        console.error('解析WebSocket消息失败:', e);
                    }
                };

                websocket.onclose = function(event) {
                    isConnected = false;
                    const wasClean = event.wasClean;
                    console.log('WebSocket连接已关闭, clean:', wasClean, 'code:', event.code);

                    // 如果有正在进行的请求，标记为失败
                    if (currentRequestId) {
                        const statusDiv = document.getElementById('status');
                        const button = document.getElementById('generateBtn');
                        const cancelBtn = document.getElementById('cancelBtn');

                        statusDiv.textContent = '连接断开，请求失败';
                        statusDiv.style.color = '#dc3545';
                        statusDiv.className = 'error';

                        // 恢复按钮状态
                        button.textContent = '🤖 生成AI回答';
                        button.disabled = false;
                        cancelBtn.style.display = 'none';
                        cancelBtn.disabled = true;

                        currentRequestId = null;
                    }

                    updateConnectionStatus('disconnected', '🔴', '已断开');

                    // 自动重连
                    if (!wasClean && reconnectAttempts < maxReconnectAttempts) {
                        setTimeout(connectWebSocket, 3000);
                    }
                };

                websocket.onerror = function(error) {
                    isConnected = false;
                    console.error('WebSocket错误:', error);
                    updateConnectionStatus('disconnected', '🔴', '连接错误');
                };

            } catch (e) {
                console.error('创建WebSocket连接失败:', e);
                updateConnectionStatus('disconnected', '🔴', '连接失败');

                if (reconnectAttempts < maxReconnectAttempts) {
                    setTimeout(connectWebSocket, 5000);
                }
            }
        }

        function updateConnectionStatus(statusClass, indicator, text) {
            const statusDiv = document.getElementById('connectionStatus');
            const indicatorSpan = document.getElementById('statusIndicator');
            const textSpan = document.getElementById('statusText');

            statusDiv.className = 'connection-status ' + statusClass;
            indicatorSpan.textContent = indicator;
            textSpan.textContent = text;
        }

        function handleWebSocketMessage(message) {
            const statusDiv = document.getElementById('status');
            const resultDiv = document.getElementById('result');
            const responseDiv = document.getElementById('response');
            const button = document.getElementById('generateBtn');
            const cancelBtn = document.getElementById('cancelBtn');

            switch (message.type) {
                case 'status':
                    const data = message.data;
                    statusDiv.textContent = data.message || '处理中...';
                    statusDiv.style.color = '#007bff';
                    statusDiv.className = 'loading';

                    if (data.stage === 'started') {
                        resultDiv.style.display = 'block';
                        responseDiv.textContent = 'AI正在思考中...';
                    }
                    break;

                case 'result':
                    const result = message.data;
                    const formattedResponse = formatResponse(result.response);
                    responseDiv.innerHTML = formattedResponse;

                    statusDiv.textContent = '回答生成完成 (' + result.length + ' 字符)';
                    statusDiv.style.color = '#28a745';
                    statusDiv.className = 'success';

                    // 清理请求状态
                    currentRequestId = null;

                    // 恢复按钮状态
                    button.textContent = '🤖 生成AI回答';
                    button.disabled = false;
                    cancelBtn.style.display = 'none';
                    cancelBtn.disabled = true;

                    // 滚动到结果区域
                    resultDiv.scrollIntoView({ behavior: 'smooth', block: 'start' });
                    break;

                case 'error':
                    statusDiv.textContent = '生成失败: ' + message.data.message;
                    statusDiv.style.color = '#dc3545';
                    statusDiv.className = 'error';

                    responseDiv.innerHTML = '<div style="color: #dc3545; padding: 15px; background: #f8d7da; border-radius: 5px; border: 1px solid #f5c6cb;"><strong>很抱歉，暂时无法生成回答</strong><br><small>' + message.data.message + '</small></div>';
                    resultDiv.style.display = 'block';

                    // 清理请求状态
                    currentRequestId = null;

                    // 恢复按钮状态
                    button.textContent = '🤖 生成AI回答';
                    button.disabled = false;
                    cancelBtn.style.display = 'none';
                    cancelBtn.disabled = true;
                    break;
            }
        }

        function cancelRequest() {
            if (currentRequestId) {
                // 取消当前请求
                currentRequestId = null;

                const statusDiv = document.getElementById('status');
                const button = document.getElementById('generateBtn');
                const cancelBtn = document.getElementById('cancelBtn');

                statusDiv.textContent = '请求已取消';
                statusDiv.style.color = '#ffc107';
                statusDiv.className = '';

                // 恢复按钮状态
                button.textContent = '🤖 生成AI回答';
                button.disabled = false;
                cancelBtn.style.display = 'none';
                cancelBtn.disabled = true;

                console.log('用户取消了当前请求');
            } else {
                // 如果没有正在进行的请求，关闭连接
                if (websocket && websocket.readyState === WebSocket.OPEN) {
                    websocket.close();
                    updateConnectionStatus('disconnected', '🔴', '已断开');
                }
            }
        }

        async function generateResponse() {
            const style = document.getElementById('style').value;
            const content = document.getElementById('content').value;
            const question = document.getElementById('question').value;

            if (!question.trim()) {
                alert('请输入问题！');
                return;
            }

            // 检查WebSocket连接
            if (!isConnected || !websocket || websocket.readyState !== WebSocket.OPEN) {
                alert('WebSocket连接未建立，请稍后重试或刷新页面');
                connectWebSocket(); // 尝试重新连接
                return;
            }

            // 如果有正在进行的请求，提示用户等待
            if (currentRequestId) {
                alert('有请求正在进行中，请等待完成后再试');
                return;
            }

            // 生成请求ID
            currentRequestId = Date.now().toString();

            // 显示加载状态
            const button = document.getElementById('generateBtn');
            const cancelBtn = document.getElementById('cancelBtn');
            const originalText = button.textContent;
            button.textContent = '🤖 AI正在深度思考中...';
            button.disabled = true;
            cancelBtn.style.display = 'inline-block';
            cancelBtn.disabled = false;

            try {
                // 发送WebSocket请求
                const requestData = {
                    action: 'generate',
                    style: style,
                    content: content,
                    question: question,
                    requestId: currentRequestId
                };

                websocket.send(JSON.stringify(requestData));
                console.log('发送WebSocket请求:', requestData);

            } catch (error) {
                console.error('WebSocket请求失败:', error);
                currentRequestId = null;

                const statusDiv = document.getElementById('status');
                const resultDiv = document.getElementById('result');
                const responseDiv = document.getElementById('response');

                statusDiv.textContent = 'WebSocket请求失败: ' + error.message;
                statusDiv.style.color = '#dc3545';
                statusDiv.className = 'error';

                responseDiv.innerHTML = '<div style="color: #dc3545; padding: 15px; background: #f8d7da; border-radius: 5px; border: 1px solid #f5c6cb;"><strong>很抱歉，暂时无法生成回答</strong><br><small>WebSocket连接问题，请刷新页面重试</small></div>';
                resultDiv.style.display = 'block';

                // 恢复按钮状态
                button.textContent = originalText;
                button.disabled = false;
                cancelBtn.style.display = 'none';
                cancelBtn.disabled = true;
            }
        }

        function formatResponse(text) {
            // 简单的文本格式化：保留换行，添加段落样式
            return text
                .split('\n')
                .map(line => line.trim() ? '<p style="margin: 8px 0; text-indent: 2em;">' + line + '</p>' : '<br>')
                .join('')
                .replace(/<p[^>]*><\/p>/g, '<br>'); // 清理空段落
        }

        // 支持回车键快速提交
        document.getElementById('question').addEventListener('keypress', function(e) {
            if (e.key === 'Enter' && !e.shiftKey) {
                e.preventDefault();
                generateResponse();
            }
        });

        // 页面加载时初始化WebSocket连接
        document.addEventListener('DOMContentLoaded', function() {
            connectWebSocket();
        });

        // 页面卸载时关闭WebSocket连接
        window.addEventListener('beforeunload', function() {
            if (websocket && websocket.readyState === WebSocket.OPEN) {
                websocket.close();
            }
        });
    </script>
    </div>
</body>
</html>`
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

// handleGenerate 生成回答
func (s *Server) handleGenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Style    string `json:"style"`
		Content  string `json:"content"`
		Question string `json:"question"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 记录AI请求详情
	fmt.Println("📥 AI请求详情:")
	fmt.Printf("   风格: %s\n", req.Style)
	fmt.Printf("   经典内容: %s\n", req.Content)
	fmt.Printf("   职场问题: %s\n", req.Question)
	fmt.Printf("   客户端IP: %s\n", getClientIP(r))

	// 使用AI服务生成风格化回答
	var response string
	var err error
	if s.aiManager != nil {
		// 使用配置的AI交互超时时间
		timeoutSeconds := 100 // 默认100秒
		if s.config != nil && s.config.AI.InteractionTimeout > 0 {
			timeoutSeconds = s.config.AI.InteractionTimeout
		}
		fmt.Printf("⏰ AI交互超时设置: %d秒\n", timeoutSeconds)
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
		defer cancel()

		response, err = s.generateAIResponse(ctx, req.Style, req.Question, req.Content)
		if err != nil {
			log.Printf("AI生成回答失败: %v", err)

			// 检查是否是配额错误，为用户提供友好的提示
			errMsg := err.Error()
			if strings.Contains(errMsg, "429") || strings.Contains(errMsg, "quota") {
				log.Println("⚠️ AI服务配额超限，已切换到本地模拟回答")
				response = fmt.Sprintf("🤖 AI服务暂时不可用（配额限制），为您提供%s风格的本地模拟回答：\n\n%s",
					req.Style, s.aiEngine.GenerateStyleResponse(req.Style, req.Question, req.Content))
			} else {
				// 其他错误也降级到本地模拟回答
				response = s.aiEngine.GenerateStyleResponse(req.Style, req.Question, req.Content)
			}
		}
	} else {
		// AI服务不可用，直接使用本地模拟回答
		response = s.aiEngine.GenerateStyleResponse(req.Style, req.Question, req.Content)
	}

	// 记录AI响应详情
	fmt.Printf("📤 AI响应详情:\n")
	fmt.Printf("   响应长度: %d 字符\n", len(response))
	if len(response) > 500 {
		fmt.Printf("   响应内容预览: %s...\n", response[:500])
	} else {
		fmt.Printf("   完整响应内容: %s\n", response)
	}
	fmt.Printf("   是否使用AI: %t\n", s.aiManager != nil)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

// handleWebSocket 处理WebSocket连接
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("新的WebSocket连接建立: %s", r.RemoteAddr)

	// 设置读写超时 - 增加超时时间以适应长AI请求
	conn.SetReadDeadline(time.Now().Add(120 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))

	// 启动心跳goroutine
	go s.handleWebSocketHeartbeat(conn)

	for {
		// 读取客户端消息
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		// 处理消息
		action, ok := msg["action"].(string)
		if !ok {
			s.sendWebSocketError(conn, "缺少action字段")
			continue
		}

		switch action {
		case "generate":
			s.handleWebSocketGenerate(conn, msg)
		default:
			s.sendWebSocketError(conn, "未知的action: "+action)
		}
	}
}

// handleWebSocketGenerate 处理WebSocket生成请求
func (s *Server) handleWebSocketGenerate(conn *websocket.Conn, msg map[string]interface{}) {
	// 解析请求参数
	style, ok := msg["style"].(string)
	if !ok {
		s.sendWebSocketError(conn, "缺少style字段")
		return
	}

	content, ok := msg["content"].(string)
	if !ok {
		s.sendWebSocketError(conn, "缺少content字段")
		return
	}

	question, ok := msg["question"].(string)
	if !ok {
		s.sendWebSocketError(conn, "缺少question字段")
		return
	}

	if question == "" {
		s.sendWebSocketError(conn, "问题不能为空")
		return
	}

	// 发送开始状态
	s.sendWebSocketMessage(conn, "status", map[string]interface{}{
		"stage":   "started",
		"message": "AI开始分析问题...",
	})

	// 记录请求详情
	fmt.Println("📥 WebSocket AI请求详情:")
	fmt.Printf("   风格: %s\n", style)
	fmt.Printf("   经典内容: %s\n", content)
	fmt.Printf("   职场问题: %s\n", question)
	fmt.Printf("   客户端: %s\n", conn.RemoteAddr())

	// 使用goroutine异步处理AI请求
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("WebSocket处理panic: %v", r)
				// 由于连接可能已断开，这里只记录日志，不发送错误消息
				log.Printf("由于panic，跳过向客户端发送错误消息")
			}
		}()

		s.processWebSocketAIRequest(conn, style, question, content)
	}()
}

// processWebSocketAIRequest 处理WebSocket AI请求
func (s *Server) processWebSocketAIRequest(conn *websocket.Conn, style, question, content string) {
	// 检查连接是否仍然有效
	if conn == nil {
		log.Printf("WebSocket连接已断开，跳过AI处理")
		return
	}

	// 获取AI交互超时时间
	timeoutSeconds := 600 // 默认600秒
	if s.config != nil && s.config.AI.InteractionTimeout > 0 {
		timeoutSeconds = s.config.AI.InteractionTimeout
	}

	fmt.Printf("⏰ WebSocket AI交互超时设置: %d秒\n", timeoutSeconds)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// 发送处理状态
	s.sendWebSocketMessage(conn, "status", map[string]interface{}{
		"stage":   "processing",
		"message": "AI正在生成风格化回答...",
	})

	var response string
	var err error

	if s.aiManager != nil {
		response, err = s.generateAIResponse(ctx, style, question, content)
		if err != nil {
			log.Printf("WebSocket AI生成回答失败: %v", err)

			// 检查是否是配额错误
			errMsg := err.Error()
			if strings.Contains(errMsg, "429") || strings.Contains(errMsg, "quota") {
				log.Println("⚠️ WebSocket AI服务配额超限，已切换到本地模拟回答")
				response = fmt.Sprintf("🤖 AI服务暂时不可用（配额限制），为您提供%s风格的本地模拟回答：\n\n%s",
					style, s.aiEngine.GenerateStyleResponse(style, question, content))

				s.sendWebSocketMessage(conn, "status", map[string]interface{}{
					"stage":   "fallback",
					"message": "AI服务配额限制，使用本地模拟回答",
				})
			} else {
				s.sendWebSocketError(conn, "AI生成失败: "+err.Error())
				return
			}
		}
	} else {
		// AI服务不可用，使用本地模拟回答
		response = s.aiEngine.GenerateStyleResponse(style, question, content)

		s.sendWebSocketMessage(conn, "status", map[string]interface{}{
			"stage":   "local",
			"message": "使用本地AI引擎生成回答",
		})
	}

	// 发送完成状态和结果
	s.sendWebSocketMessage(conn, "result", map[string]interface{}{
		"response": response,
		"length":   len(response),
	})

	fmt.Printf("📤 WebSocket AI响应详情:\n")
	fmt.Printf("   响应长度: %d 字符\n", len(response))
	if len(response) > 200 {
		fmt.Printf("   响应预览: %s...\n", response[:200])
	}
}

// sendWebSocketMessage 发送WebSocket消息
func (s *Server) sendWebSocketMessage(conn *websocket.Conn, msgType string, data interface{}) {
	// 检查连接是否仍然有效
	if conn == nil {
		log.Printf("WebSocket连接为空，跳过消息发送")
		return
	}

	message := map[string]interface{}{
		"type": msgType,
		"data": data,
		"time": time.Now().Unix(),
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err := conn.WriteJSON(message); err != nil {
		log.Printf("WebSocket发送消息失败: %v", err)
		// 可以考虑在这里关闭连接或标记连接为无效
		// 但由于这是异步处理的，不在这里处理连接管理
	}
}

// sendWebSocketError 发送WebSocket错误消息
func (s *Server) sendWebSocketError(conn *websocket.Conn, errorMsg string) {
	s.sendWebSocketMessage(conn, "error", map[string]interface{}{
		"message": errorMsg,
	})
}

// handleWebSocketHeartbeat 处理WebSocket心跳
func (s *Server) handleWebSocketHeartbeat(conn *websocket.Conn) {
	ticker := time.NewTicker(30 * time.Second) // 每30秒发送一次心跳
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 发送ping消息
			conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Printf("发送心跳失败: %v", err)
				return
			}
		}
	}
}

// generateAIResponse 使用AI服务生成风格化回答
func (s *Server) generateAIResponse(ctx context.Context, style, question, content string) (string, error) {
	// 构建风格描述
	styleDesc := getStyleDescription(style)

	// 构建提示词
	prompt := fmt.Sprintf(`你是一个职场沟通风格模仿专家，请模仿%s的沟通风格回答以下职场问题。

风格特点：%s

经典讲话内容参考：%s

职场问题：%s

请用%s的风格给出专业的回答。回答要体现该风格的核心特点，自然流畅，有说服力。

回答：`, styleDesc["name"], styleDesc["description"], content, question, styleDesc["name"])

	// 获取AI客户端
	client := s.aiManager.GetClient()

	// 根据AI模式和客户端类型选择合适的模型
	switch c := client.(type) {
	case *aiPkg.TALClient:
		// TAL客户端：根据AI模式选择模型
		var modelName string
		if s.config != nil && s.config.AI.Mode == "internal" {
			// 内部模式：使用advancedReasoning模型
			modelName = "deepseek-reasoner"
		} else {
			// 其他模式：使用textGeneration模型
			modelName = "deepseek-chat"
		}
		return c.GenerateResponseWithModel(ctx, prompt, modelName)
	case *aiPkg.SparkClient:
		// 星火客户端：使用spark-x模型
		return c.GenerateResponseWithModel(ctx, prompt, "spark-x")
	case *aiPkg.OpenAIClient:
		// OpenAI客户端：使用gpt-4
		return c.GenerateResponseWithModel(ctx, prompt, "gpt-4")
	default:
		// 其他客户端尝试通用方法
		return "", fmt.Errorf("不支持的AI客户端类型: %T", client)
	}
}

// getClientIP 获取客户端IP地址
func getClientIP(r *http.Request) string {
	// 尝试从X-Forwarded-For头获取
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// 尝试从X-Real-IP头获取
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// 从RemoteAddr获取
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}
	return ip
}

// getStyleDescription 获取风格描述
func getStyleDescription(style string) map[string]string {
	descriptions := map[string]map[string]string{
		"kanghui": {
			"name": "康辉",
			"description": "专业得体，逻辑严谨，数据支撑，权威感强，结构清晰，适合正式场合和汇报答辩",
		},
		"dongqing": {
			"name": "董卿",
			"description": "温婉大气，情感共鸣，优雅从容，善解人意，注重倾听，创造和谐沟通氛围",
		},
		"hanhan": {
			"name": "韩寒",
			"description": "犀利穿透，直言不讳，敢于挑战常规，反问拆解，态度鲜明，真诚表达",
		},
		"chengming": {
			"name": "成铭",
			"description": "逻辑严谨，层层递进，策略性强，归谬反驳，理性分析，掌控局面",
		},
	}

	if desc, exists := descriptions[style]; exists {
		return desc
	}

	return descriptions["kanghui"] // 默认康辉风格
}

