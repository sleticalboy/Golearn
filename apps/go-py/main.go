package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-python/gpython/py"
	"github.com/go-python/gpython/repl"
	"github.com/go-python/gpython/repl/cli"
	_ "github.com/go-python/gpython/stdlib"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"sync"
)

// ExecuteRequest 执行请求结构体
type ExecuteRequest struct {
	Code string `json:"code"`
}

// ExecuteResponse 响应结构体
type ExecuteResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 生产环境中应更严格
	},
}

type pytermminal struct {
	prompt string
	out    string
}

func (p *pytermminal) SetPrompt(prompt string) {
	p.prompt = prompt
}

func (p *pytermminal) Print(msg string) {
	log.Printf("msg -> %s\n", msg)
	p.out = msg
}

type PythonSession struct {
	mu      sync.Mutex
	output  chan string
	context context.Context
	cancel  context.CancelFunc
	pr      *repl.REPL
	term    *pytermminal
}

func NewPythonSession() *PythonSession {
	ctx, cancel := context.WithCancel(context.Background())
	pyrepl := repl.New(nil)
	term := &pytermminal{}
	pyrepl.SetUI(term)
	return &PythonSession{
		term:    term,
		pr:      pyrepl,
		output:  make(chan string),
		context: ctx,
		cancel:  cancel,
	}
}

func (s *PythonSession) Stop() {
	s.cancel()
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade failed:", err)
		return
	}
	defer func() { _ = conn.Close() }()
	session := NewPythonSession()
	defer session.Stop()
	cli.RunREPL(nil)
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}
		// 处理输入
		if len(msg) > 0 {
			select {
			default:
				session.pr.Run(string(msg))
				_ = conn.WriteMessage(websocket.TextMessage, []byte(session.term.out))
			case <-session.context.Done():
				return
			}
		}
	}
}

func xmain() {
	ctx := py.NewContext(py.DefaultContextOpts())
	defer func() { _ = ctx.Close() }()
	src, err := py.RunSrc(ctx, `print('Hello from embedded Python!')`, "", nil)
	log.Printf("src: %+v, err: %v\n", src, err)

	// 设置静态文件服务
	fs := http.FileServer(http.Dir("D:\\code\\Golearn\\apps\\go-py\\static"))
	http.Handle("/", fs)
	// API 端点 - 执行 Python 代码
	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		var req ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// 创建响应
		response := ExecuteResponse{}
		if !isCodeSafe(req.Code) {
			response.Error = "danger code!"
		} else {
			// 捕获输出
			// output := ""
			// py.SetStdout(py.NewFileFromString(&output))
			// py.SetStderr(py.NewFileFromString(&output))
			// 执行代码
			if m, err := py.RunSrc(ctx, req.Code, "", nil); err != nil {
				response.Error = err.Error()
			} else {
				log.Printf("result: %+v\n", m)
				buf, _ := json.Marshal(m.GetDict())
				response.Output = string(buf)
			}
		}
		// 返回 JSON 响应
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	})
	// WebSocket 端点
	http.HandleFunc("/ws", handleWebSocket)
	// 添加保存功能端点
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		// 实现代码保存逻辑
	})
	// 启动服务器
	port := "8080"
	fmt.Printf("Server running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

// 在 main.go 中添加安全检查
func isCodeSafe(code string) bool {
	// 禁止导入危险模块
	blacklist := []string{"os", "sys", "subprocess"}
	for _, module := range blacklist {
		if strings.Contains(code, "import "+module) ||
			strings.Contains(code, "from "+module) {
			return false
		}
	}
	return true
}

func main() {
	xmain()
}
