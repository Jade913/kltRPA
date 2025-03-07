package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func DealNewGreet(ctx context.Context) error {
	log.Println("开始处理新招呼")

	// 访问目标页面
	err := runWithTimeout(ctx,
		chromedp.Navigate("https://rd6.zhaopin.com/app/recommend?tab=recommend#sortType=recommend"),
		chromedp.WaitVisible("//span[contains(text(),'人才管理')]"),
	)
	if err != nil {
		return fmt.Errorf("页面加载失败: %v", err)
	}

	// 点击人才管理
	err = runWithTimeout(ctx,
		chromedp.Click("//span[contains(text(),'人才管理')]"),
	)
	if err != nil {
		return fmt.Errorf("点击人才管理失败: %v", err)
	}
	log.Println("成功点击【人才管理】")

	// 点击年龄筛选
	err = runWithTimeout(ctx,
		chromedp.Click(".age-selector"),
		chromedp.WaitVisible(".km-select__dropdown"),
	)
	if err != nil {
		return fmt.Errorf("点击年龄筛选失败: %v", err)
	}

	// 选择年龄范围
	targetAges := []string{"16-25岁", "26-30岁", "31-35岁", "36-40岁"}
	for _, age := range targetAges {
		err = runWithTimeout(ctx,
			chromedp.Click(fmt.Sprintf(`.km-select__dropdown .km-option .km-option__label div[title="%s"]`, age)),
		)
		if err != nil {
			log.Printf("选择年龄 %s 失败: %v", age, err)
			continue
		}
		log.Printf("已选择年龄：%s", age)
		time.Sleep(1 * time.Second)
	}

	// 点击确定按钮
	err = runWithTimeout(ctx,
		chromedp.Click(".km-select__dropdown-footer button.km-button--filled"),
	)

	if err != nil {
		return fmt.Errorf("点击年龄确定按钮失败: %v", err)
	}
	log.Println("筛选年龄成功")
	time.Sleep(1 * time.Second)

	// 点击学历筛选
	err = runWithTimeout(ctx,
		chromedp.Click(".edu-selector"),
		// chromedp.WaitVisible(".km-select__dropdown"),
	)
	if err != nil {
		return fmt.Errorf("点击学历筛选失败: %v", err)
	}

	// 选择学历
	targetEducations := []string{"大专", "本科", "硕士"}
	for _, edu := range targetEducations {
		err = runWithTimeout(ctx,
			chromedp.Click(fmt.Sprintf(`.km-select__dropdown .km-option .km-option__label div[title="%s"]`, edu)),
		)
		if err != nil {
			log.Printf("选择学历 %s 失败: %v", edu, err)
			continue
		}
		log.Printf("已选择学历：%s", edu)
	}

	// 点击确定按钮
	err = runWithTimeout(ctx,
		// chromedp.Click(".km-select__dropdown-footer button.km-button--primary"),
		chromedp.Evaluate(`document.querySelectorAll('.km-select__dropdown-footer button.km-button--filled')[1].click()`, nil), // 点击第二个确定按钮
	)
	if err != nil {
		return fmt.Errorf("点击确定按钮失败: %v", err)
	}
	log.Println("筛选学历成功")

	// 点击未读复选框
	err = runWithTimeout(ctx,
		chromedp.WaitVisible(".not-read"),
		chromedp.Click(".not-read .km-checkbox__icon"),
	)
	if err != nil {
		log.Println("点击未读复选框失败")
	}
	time.Sleep(1 * time.Second)

	// 获取总结果数
	var totalResults int
	err = runWithTimeout(ctx,
		chromedp.Evaluate(`
			const countEl = document.querySelector('.search-result-count .number');
			countEl ? parseInt(countEl.textContent) : 0;
		`, &totalResults),
	)
	if err != nil {
		log.Printf("获取总结果数失败: %v", err)
		totalResults = 0
	}

	// 计算需要循环的次数（每次处理20个，向上取整）
	loopCount := (totalResults + 19) / 20
	log.Printf("总共有 %d 个结果，需要循环 %d 次", totalResults, loopCount)

	// 循环处理全选和打招呼
	for i := 0; i < loopCount; i++ {
		log.Printf("开始第 %d/%d 轮处理", i+1, loopCount)

		// 等待页面加载
		err = runWithTimeout(ctx,
			chromedp.WaitVisible(".page-action-bar"),
			chromedp.Sleep(4*time.Second),
		)
		if err != nil {
			log.Printf("找不到页面操作栏，跳过当前轮次")
			continue
		}

		// 点击全选
		err = runWithTimeout(ctx,
			chromedp.Click(".footer-action-bar__inner .km-checkbox__icon"),
		)
		if err != nil {
			log.Println("点击全选失败，退出循环")
			break
		}

		// 点击批量回复按钮
		err = runWithTimeout(ctx,
			chromedp.Click(".footer-action-bar__end .is-ml-12 button.km-button--filled"),
		)
		if err != nil {
			log.Println("点击【批量回复】失败")
			continue
		}

		// 处理可能出现的打招呼模态框
		err = runWithTimeout(ctx,
			chromedp.WaitVisible(".setting-greet", chromedp.ByQuery),
		)
		if err == nil {
			// 如果模态框出现，点击发送按钮
			err = runWithTimeout(ctx,
				chromedp.Click(".setting-greet__footer button.km-button--filled"),
				chromedp.Sleep(1*time.Second),
			)
			if err != nil {
				log.Println("点击发送按钮失败，尝试使用JavaScript点击")
				// 使用JavaScript作为备选方案
				err = runWithTimeout(ctx,
					chromedp.Evaluate(`
						const sendBtn = document.querySelector('.setting-greet__footer button.km-button--filled');
						if (sendBtn) {
							sendBtn.click();
							return true;
						}
						return false;
					`, nil),
				)
				if err != nil {
					log.Println("JavaScript点击发送按钮也失败")
				} else {
					log.Println("通过JavaScript成功点击发送按钮")
				}
			} else {
				log.Println("点击了发送按钮")
			}
		}

		log.Printf("完成第 %d/%d 轮处理", i+1, loopCount)
	}

	log.Printf("所有轮次处理完成，共循环 %d 次", loopCount)
	return nil
}
