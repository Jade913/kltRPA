package utils

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

// InitChrome 初始化 Chrome 实例
func InitChrome() (context.Context, context.CancelFunc, error) {
	// 获取当前工作目录
	projectDir, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前工作目录失败: %v", err)
		return nil, nil, err
	}

	// 设置 Chrome 用户数据目录
	userDataDir := filepath.Join(projectDir, "chrome-data")
	println(userDataDir)

	// 设置下载目录
	downloadDir := filepath.Join(projectDir, "export_data", "resumeRPA-data")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		log.Printf("创建下载目录失败: %v", err)
		return nil, nil, err
	}
	absDownloadDir, err := filepath.Abs(downloadDir)
	if err != nil {
		log.Printf("获取下载目录绝对路径失败: %v", err)
		return nil, nil, err
	}

	// 清理可能存在的锁文件
	lockFile := filepath.Join(userDataDir, "SingletonLock")
	if err := os.Remove(lockFile); err != nil && !os.IsNotExist(err) {
		log.Printf("清理锁文件失败: %v", err)
	}

	// 配置 Chrome 启动选项
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.WindowSize(1280, 800),
		chromedp.UserDataDir(userDataDir),
		chromedp.Flag("no-first-run", true),
		chromedp.Flag("no-default-browser-check", true),
		// 下载相关配置
		chromedp.Flag("download.default_directory", absDownloadDir),
		chromedp.Flag("download.prompt_for_download", false),
		chromedp.Flag("download.directory_upgrade", true),
		chromedp.Flag("safebrowsing.enabled", true),
	}
	opts = append(chromedp.DefaultExecAllocatorOptions[:], opts...)

	// 创建 Chrome 实例
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)

	// 创建新的 Chrome 上下文
	ctx, cancel := chromedp.NewContext(allocCtx)

	// 设置全局超时
	ctx, cancel = context.WithTimeout(ctx, 30*time.Minute)

	return ctx, cancel, nil
}
