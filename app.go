package main

import (
	"bufio"
	"context"
	"fmt"
	"kltRPA/logs"
	"kltRPA/models"
	"kltRPA/utils"
	"os"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	logs.InitLog()
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// omo登录
func (a *App) Login(username, password string) string {
	omo := utils.NewOmoIntegrate("omo.kelote.com", "klt_omo", username, password)

	success, err := omo.Login()
	if success {
		return "Login successful!"
	} else {
		if err != nil {
			return fmt.Sprintf("Login failed: %v", err)
		} else {
			return "Invalid username or password."
		}
	}
}

func (a *App) RunRPA() {
	models.RunRPA()
}

// GetLogs returns the last 30 lines of the log file
func (a *App) GetLogs() (string, error) {
	file, err := os.Open("logs/app.log")
	if err != nil {
		return "", fmt.Errorf("无法读取日志文件: %v", err)
	}
	defer file.Close()

	// 使用一个切片来存储最后30行
	var lines []string
	scanner := bufio.NewScanner(file)

	// 使用一个循环来读取文件的每一行
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > 30 {
			lines = lines[1:] // 保持切片中只有最后30行
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取日志文件失败: %v", err)
	}

	// 使用换行符连接每一行
	return fmt.Sprintf("%s", strings.Join(lines, "\n")), nil
}
