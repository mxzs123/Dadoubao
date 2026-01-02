package main

import (
	"embed"
	"fmt"
	"image/png"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"voicebridge/internal/clipboard"
	"voicebridge/internal/keyboard"
	"voicebridge/internal/network"
	"voicebridge/internal/qrcode"
	"voicebridge/internal/server"
)

//go:embed web/static/*
var staticFiles embed.FS

const port = 8765

func main() {
	log.SetFlags(log.Ltime)
	log.Println("🚀 VoiceBridge 启动中...")

	// 初始化剪贴板
	if err := clipboard.Init(); err != nil {
		log.Println("⚠️  剪贴板初始化失败:", err)
		log.Println("   程序将继续运行，但粘贴功能可能不可用")
	}

	// 获取本机 IP
	ip, err := network.GetLocalIP()
	if err != nil {
		log.Fatal("❌ 获取 IP 失败:", err)
	}
	url := fmt.Sprintf("http://%s:%d", ip, port)

	// WebSocket 处理器
	wsHandler := &server.Handler{
		OnTextReceived: func(text string, autoEnter bool) {
			log.Printf("📝 收到文字: %s (自动回车: %v)\n", text, autoEnter)
			clipboard.WriteText(text)
			time.Sleep(50 * time.Millisecond) // 等待剪贴板写入
			if err := keyboard.SimulatePaste(); err != nil {
				log.Println("⚠️  粘贴失败:", err)
			} else {
				log.Println("✅ 已粘贴到光标位置")
			}
			// 自动回车
			if autoEnter {
				time.Sleep(30 * time.Millisecond)
				if err := keyboard.SimulateEnter(); err != nil {
					log.Println("⚠️  回车失败:", err)
				} else {
					log.Println("↵ 已自动回车")
				}
			}
		},
		OnConnect: func() {
			log.Println("📱 手机已连接")
		},
		OnDisconnect: func() {
			log.Println("📱 手机已断开")
		},
	}

	// 静态文件服务
	staticFS, err := fs.Sub(staticFiles, "web/static")
	if err != nil {
		log.Fatal("❌ 静态文件加载失败:", err)
	}

	// 路由
	http.Handle("/", http.FileServer(http.FS(staticFS)))
	http.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// 生成并保存二维码
	qrImg, err := qrcode.Generate(url, 256)
	if err != nil {
		log.Println("⚠️  二维码生成失败:", err)
	} else {
		// 保存到临时文件
		qrFile := "voicebridge-qr.png"
		f, err := os.Create(qrFile)
		if err == nil {
			png.Encode(f, qrImg)
			f.Close()
			log.Printf("📷 二维码已保存: %s\n", qrFile)
		}
	}

	// 打印启动信息
	fmt.Println()
	fmt.Println("════════════════════════════════════════")
	fmt.Println("  VoiceBridge 语音桥 已启动!")
	fmt.Println("════════════════════════════════════════")
	fmt.Printf("  📱 手机扫码或访问: %s\n", url)
	fmt.Println("  📷 二维码文件: voicebridge-qr.png")
	fmt.Println("  💡 提示: 电脑光标放到输入位置，手机输入即可")
	fmt.Println("════════════════════════════════════════")
	fmt.Println()

	// 优雅退出
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		fmt.Println("\n👋 正在退出...")
		os.Exit(0)
	}()

	// 启动服务器
	log.Printf("🌐 HTTP 服务器运行在 :%d\n", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("❌ 服务器启动失败:", err)
	}
}
