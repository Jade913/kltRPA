package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"kltRPA/logs"
	"kltRPA/models"
	"kltRPA/utils"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// App struct
type App struct {
	ctx         context.Context
	omoInstance *utils.OmoIntegrate
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
	if err != nil {
		return fmt.Sprintf("登录失败: %v", err)
	}

	if !success {
		return "用户名或密码错误"
	}

	a.omoInstance = omo // 保存登录成功的实例
	return "登录成功！"
}

func (a *App) UpdateOmo(data []map[string]interface{}) ([]map[string]interface{}, error) {
	if a.omoInstance == nil {
		return nil, fmt.Errorf("未登录！")
	}
	return a.omoInstance.UpdateOmo(data)
}

func (a *App) RunRPA(selectedCampuses []string) {
	models.RunRPA(selectedCampuses)
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

// GetLatestTable returns the path of the latest table file
func (a *App) GetLatestTable() (string, error) {
	dir := "export_data/table-data"
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("无法读取目录: %v", err)
	}

	// Convert []os.DirEntry to []os.FileInfo
	fileInfos := make([]os.FileInfo, 0, len(files))
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			return "", fmt.Errorf("无法获取文件信息: %v", err)
		}
		fileInfos = append(fileInfos, info)
	}

	// 按修改时间排序文件
	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ModTime().After(fileInfos[j].ModTime())
	})

	// 找到最新的表格文件
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".xlsx" {
			return filepath.Join(dir, fileInfo.Name()), nil
		}
	}

	return "", fmt.Errorf("没有找到表格文件")
}

// ServeFile serves a file from the given path
func (a *App) ServeFile(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		http.Error(w, "File path is required", http.StatusBadRequest)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("无法打开文件: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 使用 io.Copy 直接将文件内容写入响应
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	if _, err := io.Copy(w, file); err != nil {
		http.Error(w, fmt.Sprintf("无法读取文件: %v", err), http.StatusInternalServerError)
	}
}

func (a *App) ImportTableFromExcel(filePath string) ([][]string, error) {
	return utils.ImportTableFromExcel(filePath)
}

func (a *App) SaveFile(filePath string, data []byte) error {
	return utils.SaveFile(filePath, data)
}
