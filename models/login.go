package models

import (
	"context"
	"log"

	"github.com/chromedp/chromedp"
)

// TestLogin 测试登录状态
func TestLogin(ctx context.Context) (bool, error) {
	var elementText string

	err := chromedp.Run(ctx,
		// 访问智联招聘页面
		chromedp.Navigate("https://rd6.zhaopin.com/app/recommend?tab=recommend#sortType=recommend"),

		// 等待页面标题元素出现并获取文本
		chromedp.WaitVisible(".app-header__title", chromedp.ByQuery),
		chromedp.Text(".app-header__title", &elementText, chromedp.ByQuery),
	)

	if err != nil {
		log.Printf("检查登录状态时出错: %v", err)
		return false, err
	}

	isLoggedIn := elementText == "推荐人才" || elementText == "首页"
	if isLoggedIn {
		log.Printf("找到有效的标题文本: %s", elementText)
	} else {
		log.Printf("未找到预期的标题文本，当前标题为: %s", elementText)
	}

	return isLoggedIn, nil
}
