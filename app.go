package main

import (
	"archive/zip"
	"bufio"
	"context"
	"fmt"
	"io"
	"kltRPA/logs"
	"kltRPA/models"
	"kltRPA/utils"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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

func (a *App) SayHi() string {
	err := models.SayHi()
	if err != nil {
		return fmt.Sprintf("打招呼失败: %v", err)
	}
	return "打招呼成功"
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

func (a *App) ImportTableFromExcel(filePath string) ([][]string, error) {
	return utils.ImportTableFromExcel(filePath)
}

func (a *App) SaveFile(filePath string, data []byte) error {
	return utils.SaveFile(filePath, data)
}

func (a *App) GetDownloadPath() string {
	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("获取用户主目录失败: %v", err)
		return ""
	}

	// 根据操作系统获取默认下载路径
	var downloadDir string
	switch runtime.GOOS {
	case "windows":
		downloadDir = filepath.Join(homeDir, "Downloads")
	case "darwin": // macOS
		downloadDir = filepath.Join(homeDir, "Downloads")
	case "linux":
		downloadDir = filepath.Join(homeDir, "Downloads")
	default:
		downloadDir = filepath.Join(homeDir, "Downloads")
	}

	absPath, err := filepath.Abs(downloadDir)
	if err != nil {
		fmt.Printf("获取下载目录绝对路径失败: %v", err)
		return ""
	}

	fmt.Printf("Chrome 默认下载路径: %s\n", absPath)
	return absPath
}

// 检查文件是否存在
func (a *App) CheckResumeFile(name string, position string, campus string) (string, error) {
	downloadDir := a.GetDownloadPath()
	fmt.Printf("正在搜索目录: %s\n", downloadDir)
	fmt.Printf("搜索简历: 姓名=%s, 岗位=%s, 校区=%s\n", name, position, campus)

	files, err := os.ReadDir(downloadDir)
	if err != nil {
		fmt.Printf("读取目录失败: %v\n", err)
		return "", fmt.Errorf("读取目录失败: %v", err)
	}

	for _, file := range files {
		fileName := file.Name()
		// 只检查带有"简历"字样的PDF文件
		if !strings.HasSuffix(strings.ToLower(fileName), ".pdf") ||
			!strings.Contains(strings.ToLower(fileName), "简历") {
			continue
		}

		fmt.Printf("检查文件: %s\n", fileName)
		// 检查文件名是否同时包含姓名、岗位和校区
		if strings.Contains(strings.ToLower(fileName), strings.ToLower(name)) &&
			strings.Contains(strings.ToLower(fileName), strings.ToLower(position)) &&
			strings.Contains(strings.ToLower(fileName), strings.ToLower(campus)) {
			return filepath.Join(downloadDir, fileName), nil
		}
	}

	return "", nil
}

// PackResumes 打包简历
func (a *App) PackResumes(groupedData map[string]map[string][]string) (string, error) {
	// 创建临时目录用于存放打包文件
	tempDir := filepath.Join(a.GetDownloadPath(), "KLTRPA简历打包")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", fmt.Errorf("创建临时目录失败: %v", err)
	}

	// 创建zip文件
	timestamp := time.Now().Format("20060102_150405")
	zipPath := filepath.Join(tempDir, fmt.Sprintf("简历_%s.zip", timestamp))
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("创建zip文件失败: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 遍历校区
	for campus, positions := range groupedData {
		// 遍历岗位
		for position, files := range positions {
			// 创建校区/岗位目录
			dirPath := filepath.Join(campus, position)

			// 遍历文件
			for _, filePath := range files {
				// 读取源文件
				srcFile, err := os.Open(filePath)
				if err != nil {
					fmt.Printf("打开文件失败 %s: %v\n", filePath, err)
					continue
				}
				defer srcFile.Close()

				// 在zip中创建文件
				fileName := filepath.Base(filePath)
				zipPath := filepath.Join(dirPath, fileName)
				writer, err := zipWriter.Create(zipPath)
				if err != nil {
					fmt.Printf("创建zip条目失败 %s: %v\n", zipPath, err)
					continue
				}

				// 复制文件内容
				if _, err := io.Copy(writer, srcFile); err != nil {
					fmt.Printf("复制文件内容失败 %s: %v\n", filePath, err)
					continue
				}

				fmt.Printf("已添加文件: %s\n", zipPath)
			}
		}
	}

	return zipPath, nil
}
