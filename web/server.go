package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"reactedge/internal/ai"
)

// Server Web服务器
type Server struct {
	aiEngine *ai.HanStyleAI
	router   *http.ServeMux
}

// NewServer 创建Web服务器
func NewServer(aiEngine *ai.HanStyleAI) *Server {
	server := &Server{
		aiEngine: aiEngine,
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

	// 生成回答
	response := s.aiEngine.GenerateStyleResponse(req.Style, req.Question, req.Content)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"response": response})
}

