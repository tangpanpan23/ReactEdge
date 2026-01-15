package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"reactedge/config"
	"reactedge/internal/ai"
	aiPkg "reactedge/pkg/ai"
)

// Server Web服务器
type Server struct {
	aiEngine *ai.HanStyleAI
	aiManager *aiPkg.Manager
	config   *config.Config
	router   *http.ServeMux
}

// NewServer 创建Web服务器
func NewServer(aiEngine *ai.HanStyleAI, aiManager *aiPkg.Manager, config *config.Config) *Server {
	server := &Server{
		aiEngine: aiEngine,
		aiManager: aiManager,
		config:   config,
		router:   http.NewServeMux(),
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
        body { font-family: Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; }
        .step { margin: 20px 0; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        .form-group { margin: 10px 0; }
        label { display: block; margin-bottom: 5px; }
        select, textarea { width: 100%; padding: 8px; border: 1px solid #ddd; border-radius: 3px; }
        .button { background: #28a745; color: white; padding: 10px 20px; border: none; border-radius: 5px; cursor: pointer; }
        .button:hover { background: #218838; }
        .result { margin-top: 20px; padding: 15px; background: #f8f9fa; border-radius: 5px; }
    </style>
</head>
<body>
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
            <textarea id="question" rows="3" placeholder="例如：领导问我这个项目的ROI为什么这么低？"></textarea>
        </div>
        <button class="button" onclick="generateResponse()">生成回答</button>
    </div>

    <div id="result" class="result" style="display: none;">
        <h3>🤖 生成的回答</h3>
        <div id="response"></div>
    </div>

    <script>
        async function generateResponse() {
            const style = document.getElementById('style').value;
            const content = document.getElementById('content').value;
            const question = document.getElementById('question').value;

            if (!question.trim()) {
                alert('请输入问题！');
                return;
            }

            const response = await fetch('/generate', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ style, content, question })
            });

            const data = await response.json();

            document.getElementById('response').innerHTML = '<pre>' + data.response + '</pre>';
            document.getElementById('result').style.display = 'block';
        }
    </script>
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

	// 使用AI服务生成风格化回答
	var response string
	var err error
	if s.aiManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		response, err = s.generateAIResponse(ctx, req.Style, req.Question, req.Content)
		if err != nil {
			log.Printf("AI生成回答失败: %v", err)
			// 降级到本地模拟回答
			response = s.aiEngine.GenerateStyleResponse(req.Style, req.Question, req.Content)
		}
	} else {
		// AI服务不可用，直接使用本地模拟回答
		response = s.aiEngine.GenerateStyleResponse(req.Style, req.Question, req.Content)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
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

	// 使用qwen3-max模型进行推理
	client := s.aiManager.GetClient()

	// 直接使用TAL客户端的底层API调用
	if talClient, ok := client.(*aiPkg.TALClient); ok {
		return talClient.GenerateResponseWithModel(ctx, prompt, "qwen3-max")
	}

	// 如果不是TAL客户端，使用通用方法
	// 这里暂时返回错误，后续可以扩展
	return "", fmt.Errorf("不支持的AI客户端类型")
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

