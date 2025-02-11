package utils

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/chromedp/chromedp"
)

type ChromeManager struct {
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

var (
	manager *ChromeManager
	once    sync.Once
)

// GetChromeManager 获取单例的 ChromeManager
func GetChromeManager() *ChromeManager {
	once.Do(func() {
		manager = &ChromeManager{}
	})
	return manager
}

// InitChrome 初始化 Chrome，如果已经初始化则返回现有的 context
func (cm *ChromeManager) InitChrome() (context.Context, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	// 如果已经初始化且仍然有效，直接返回
	if cm.ctx != nil {
		select {
		case <-cm.ctx.Done():
			// context 已关闭，需要重新初始化
		default:
			return cm.ctx, nil
		}
	}

	// 获取当前工作目录
	projectDir, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前工作目录失败: %v", err)
		return nil, err
	}

	// 设置 Chrome 用户数据目录
	userDataDir := filepath.Join(projectDir, "chrome-data")
	println(userDataDir)

	// 设置下载目录
	downloadDir := filepath.Join(projectDir, "export_data", "resumeRPA-data")
	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		log.Printf("创建下载目录失败: %v", err)
		return nil, err
	}
	absDownloadDir, err := filepath.Abs(downloadDir)
	if err != nil {
		log.Printf("获取下载目录绝对路径失败: %v", err)
		return nil, err
	}

	// 清理可能存在的锁文件
	lockFile := filepath.Join(userDataDir, "SingletonLock")
	if err := os.Remove(lockFile); err != nil && !os.IsNotExist(err) {
		log.Printf("清理锁文件失败: %v", err)
	}

	// 配置 Chrome 启动选项
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-site-isolation-trials", true),
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
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx)

	cm.ctx = ctx
	cm.cancel = cancel

	return ctx, nil
}

// GetContext 获取当前的 context，如果没有初始化则返回 nil
func (cm *ChromeManager) GetContext() context.Context {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	return cm.ctx
}

// CloseChrome 关闭 Chrome
func (cm *ChromeManager) CloseChrome() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.cancel != nil {
		cm.cancel()
		cm.ctx = nil
		cm.cancel = nil
	}
}

// IsContextValid 检查当前 context 是否有效
func (cm *ChromeManager) IsContextValid() bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if cm.ctx == nil {
		return false
	}

	select {
	case <-cm.ctx.Done():
		return false
	default:
		return true
	}
}

// EnsureChrome 确保 Chrome 已初始化，如果没有则初始化
func (cm *ChromeManager) EnsureChrome() (context.Context, error) {
	if !cm.IsContextValid() {
		return cm.InitChrome()
	}
	return cm.GetContext(), nil
}
