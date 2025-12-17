package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"reactedge/internal/challenge"
)

// Server Web服务器
type Server struct {
	challengeManager *challenge.ChallengeManager
	router          *http.ServeMux
	templates       *template.Template
}

// NewServer 创建Web服务器
func NewServer(cm *challenge.ChallengeManager) *Server {
	server := &Server{
		challengeManager: cm,
		router:          http.NewServeMux(),
	}

	server.setupRoutes()
	server.loadTemplates()

	return server
}

// Router 获取路由器
func (s *Server) Router() *http.ServeMux {
	return s.router
}

// setupRoutes 设置路由
func (s *Server) setupRoutes() {
	s.router.HandleFunc("/", s.handleHome)
	s.router.HandleFunc("/challenge/start", s.handleStartChallenge)
	s.router.HandleFunc("/challenge/state", s.handleGetState)
	s.router.HandleFunc("/challenge/next", s.handleNextPhase)
	s.router.HandleFunc("/challenge/speech", s.handleSubmitSpeech)
	s.router.HandleFunc("/challenge/profile", s.handleUpdateProfile)
	s.router.HandleFunc("/static/", s.handleStatic)
}

// loadTemplates 加载模板
func (s *Server) loadTemplates() {
	s.templates = template.Must(template.New("main").Parse(`
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>AI酷表达实验室 · 韩寒特训版</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #333;
            min-height: 100vh;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            background: white;
            border-radius: 10px;
            box-shadow: 0 10px 30px rgba(0,0,0,0.1);
            margin-top: 20px;
            margin-bottom: 20px;
        }
        .header {
            text-align: center;
            margin-bottom: 30px;
        }
        .title {
            font-size: 2.5em;
            margin-bottom: 10px;
            background: linear-gradient(45deg, #FF6B6B, #4ECDC4);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
        }
        .subtitle {
            font-size: 1.2em;
            color: #666;
            margin-bottom: 20px;
        }
        .phase-content {
            margin: 20px 0;
            padding: 20px;
            border-radius: 8px;
            background: #f8f9fa;
        }
        .timer {
            position: fixed;
            top: 20px;
            right: 20px;
            background: rgba(255,255,255,0.9);
            padding: 10px 20px;
            border-radius: 25px;
            font-weight: bold;
            box-shadow: 0 4px 12px rgba(0,0,0,0.1);
        }
        .button {
            background: linear-gradient(45deg, #FF6B6B, #4ECDC4);
            color: white;
            border: none;
            padding: 12px 24px;
            border-radius: 25px;
            cursor: pointer;
            font-size: 16px;
            margin: 10px;
            transition: transform 0.2s;
        }
        .button:hover {
            transform: translateY(-2px);
        }
        .speech-input {
            width: 100%;
            min-height: 150px;
            padding: 15px;
            border: 2px solid #ddd;
            border-radius: 8px;
            font-size: 16px;
            margin: 10px 0;
            resize: vertical;
        }
        .analysis-result {
            background: linear-gradient(45deg, #f093fb 0%, #f5576c 100%);
            color: white;
            padding: 20px;
            border-radius: 8px;
            margin: 10px 0;
        }
        .score {
            font-size: 2em;
            font-weight: bold;
            text-align: center;
        }
        .tags {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            margin: 10px 0;
        }
        .tag {
            background: rgba(255,255,255,0.2);
            padding: 5px 10px;
            border-radius: 15px;
            font-size: 14px;
        }
        .weapon {
            background: white;
            margin: 10px 0;
            padding: 15px;
            border-radius: 8px;
            border-left: 4px solid #FF6B6B;
        }
        .template {
            background: #e8f5e8;
            padding: 15px;
            border-radius: 8px;
            font-family: monospace;
            white-space: pre-line;
            margin: 10px 0;
        }
        .framework {
            background: #fff3cd;
            padding: 10px;
            border-radius: 5px;
            margin: 5px 0;
        }
        @media (max-width: 768px) {
            .container { margin: 10px; }
            .title { font-size: 2em; }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1 class="title">🎤 AI酷表达实验室</h1>
            <h2 class="subtitle">韩寒特训版 · 3分钟挑战</h2>
        </div>

        <div id="timer" class="timer">03:00</div>

        <div id="content" class="phase-content">
            <!-- 动态内容将在这里显示 -->
        </div>
    </div>

    <script>
        let currentState = null;
        let timerInterval = null;

        // 启动挑战
        function startChallenge() {
            fetch('/challenge/start', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    currentState = data;
                    updateUI();
                    startTimer();
                });
        }

        // 获取状态
        function getState() {
            fetch('/challenge/state')
                .then(response => response.json())
                .then(data => {
                    currentState = data;
                    updateUI();
                });
        }

        // 下一阶段
        function nextPhase() {
            fetch('/challenge/next', { method: 'POST' })
                .then(response => response.json())
                .then(data => {
                    currentState = data;
                    updateUI();
                });
        }

        // 提交语音
        function submitSpeech() {
            const speech = document.getElementById('speech-input').value;
            fetch('/challenge/speech', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ speech: speech })
            })
                .then(response => response.json())
                .then(data => {
                    currentState = data;
                    updateUI();
                });
        }

        // 更新UI
        function updateUI() {
            if (!currentState) return;

            updateTimer();
            updateContent();
        }

        // 更新计时器
        function updateTimer() {
            const timerEl = document.getElementById('timer');
            const remaining = currentState.time_remaining || 0;
            const minutes = Math.floor(remaining / 60);
            const seconds = remaining % 60;
            timerEl.textContent = String(minutes).padStart(2, '0') + ':' + String(seconds).padStart(2, '0');
        }

        // 启动计时器
        function startTimer() {
            if (timerInterval) clearInterval(timerInterval);
            timerInterval = setInterval(() => {
                if (currentState && currentState.time_remaining > 0) {
                    currentState.time_remaining--;
                    updateTimer();
                }
            }, 1000);
        }

        // 更新内容
        function updateContent() {
            const contentEl = document.getElementById('content');
            const content = currentState.content || {};

            let html = '';

            switch (currentState.current_phase) {
                case 0: // PhaseWelcome
                    html = `
                        <h3>${content.title}</h3>
                        <p style="white-space: pre-line; margin: 20px 0;">${content.description}</p>
                        <button class="button" onclick="nextPhase()">开始挑战 🚀</button>
                    `;
                    break;

                case 1: // PhaseAIDeconstruction
                    html = `<h3>${content.title}</h3>`;
                    if (content.weapons) {
                        content.weapons.forEach(weapon => {
                            html += `
                                <div class="weapon">
                                    <h4>${weapon.name}</h4>
                                    <p style="white-space: pre-line;">${weapon.description}</p>
                                </div>
                            `;
                        });
                    }
                    if (content.tools) {
                        html += '<h4>🛠️ 你的专属工具箱：</h4><ul>';
                        content.tools.forEach(tool => {
                            html += `<li>${tool}</li>`;
                        });
                        html += '</ul>';
                    }
                    html += '<button class="button" onclick="nextPhase()">继续 →</button>';
                    break;

                case 2: // PhasePersonalizedTemplate
                    html = `
                        <h3>${content.title}</h3>
                        <p>${content.profile_detection}</p>
                        <h4>${content.template_title}</h4>
                        <div class="template">${content.template}</div>
                        <h4>💡 你的表达框架：</h4>
                    `;
                    if (content.framework) {
                        content.framework.forEach(item => {
                            html += `<div class="framework">${item}</div>`;
                        });
                    }
                    html += '<button class="button" onclick="nextPhase()">开始录音 🎤</button>';
                    break;

                case 3: // PhaseRecording
                    html = `
                        <h3>${content.title}</h3>
                        <p>${content.instruction}</p>
                        <p>${content.tips}</p>
                        <p><strong>题目：${content.topic}</strong></p>
                        <textarea id="speech-input" class="speech-input" placeholder="在这里输入你的回答..."></textarea>
                        <button class="button" onclick="submitSpeech()">提交回答 📤</button>
                    `;
                    break;

                case 4: // PhaseDNAAnalysis
                    html = `<h3>${content.title}</h3>`;
                    if (content.sharpeness_score !== undefined) {
                        html += `
                            <div class="analysis-result">
                                <div class="score">🔥 犀利指数：${content.sharpeness_score}/100</div>
                            </div>
                            <h4>🎯 个性标签：</h4>
                            <div class="tags">
                        `;
                        if (content.personality_tags) {
                            content.personality_tags.forEach(tag => {
                                html += `<span class="tag">${tag}</span>`;
                            });
                        }
                        html += `
                            </div>
                            <h4>💎 发现你的独家表达模式：</h4>
                            <ul>
                        `;
                        if (content.unique_patterns) {
                            content.unique_patterns.forEach(pattern => {
                                html += `<li>${pattern}</li>`;
                            });
                        }
                        html += `
                            </ul>
                            <h4>🆙 AI优化建议：</h4>
                            <ul>
                        `;
                        if (content.recommendations) {
                            content.recommendations.forEach(rec => {
                                html += `<li>${rec}</li>`;
                            });
                        }
                        html += `
                            </ul>
                            <h4>🎮 明日挑战预告：</h4>
                            <p>${content.next_challenge}</p>
                        `;
                    }
                    html += '<button class="button" onclick="nextPhase()">完成挑战 🎉</button>';
                    break;

                case 5: // PhaseComplete
                    html = `
                        <h3>${content.title}</h3>
                        <p>${content.message}</p>
                        <button class="button" onclick="window.location.reload()">再来一次挑战 🔄</button>
                    `;
                    break;
            }

            contentEl.innerHTML = html;
        }

        // 页面加载完成后自动获取状态
        document.addEventListener('DOMContentLoaded', function() {
            getState();
        });
    </script>
</body>
</html>
`))
}

// handleHome 处理首页
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	s.templates.Execute(w, nil)
}

// handleStartChallenge 处理开始挑战
func (s *Server) handleStartChallenge(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := s.getUserID(r)
	state := s.challengeManager.StartChallenge(userID)
	content := s.challengeManager.GetPhaseContent(state)
	state.Content = content

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// handleGetState 处理获取状态
func (s *Server) handleGetState(w http.ResponseWriter, r *http.Request) {
	userID := s.getUserID(r)
	state := s.challengeManager.GetChallengeState(userID)

	if state == nil {
		// 如果没有状态，返回欢迎页面
		state = &challenge.ChallengeState{
			CurrentPhase: challenge.PhaseWelcome,
			TimeRemaining: 180,
		}
	}

	content := s.challengeManager.GetPhaseContent(state)
	state.Content = content

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// handleNextPhase 处理下一阶段
func (s *Server) handleNextPhase(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := s.getUserID(r)
	state := s.challengeManager.AdvancePhase(userID)

	if state == nil {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	content := s.challengeManager.GetPhaseContent(state)
	state.Content = content

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// handleSubmitSpeech 处理提交语音
func (s *Server) handleSubmitSpeech(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Speech string `json:"speech"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := s.getUserID(r)
	state := s.challengeManager.SubmitSpeech(userID, req.Speech)

	if state == nil {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	// 自动推进到下一阶段（分析阶段）
	state = s.challengeManager.AdvancePhase(userID)
	content := s.challengeManager.GetPhaseContent(state)
	state.Content = content

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// handleUpdateProfile 处理更新用户画像
func (s *Server) handleUpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var profile ai.UserProfile
	if err := json.NewDecoder(r.Body).Decode(&profile); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userID := s.getUserID(r)
	state := s.challengeManager.UpdateProfile(userID, profile)

	if state == nil {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	content := s.challengeManager.GetPhaseContent(state)
	state.Content = content

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// handleStatic 处理静态文件
func (s *Server) handleStatic(w http.ResponseWriter, r *http.Request) {
	// 简单的静态文件处理
	http.NotFound(w, r)
}

// getUserID 获取用户ID（简化版，使用session或IP）
func (s *Server) getUserID(r *http.Request) string {
	// 简化版：使用IP地址作为用户ID
	ip := r.RemoteAddr
	if colon := strings.LastIndex(ip, ":"); colon != -1 {
		ip = ip[:colon]
	}
	return ip
}
