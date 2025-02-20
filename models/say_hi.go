package models

import (
	"context"
	"fmt"
	"log"
	"time"

	"kltRPA/utils"

	"github.com/chromedp/chromedp"
)

func SayHi() error {
	log.Println("开始处理打招呼...")

	// 获取 Chrome 实例
	chromeManager := utils.GetChromeManager()
	ctx, err := chromeManager.EnsureChrome()
	if err != nil {
		log.Printf("初始化 Chrome 失败: %v", err)
		return err
	}

	// 检查登录状态
	isLoggedIn, err := TestLogin(ctx)
	if err != nil {
		log.Printf("检查登录状态失败: %v", err)
		return err
	}
	if !isLoggedIn {
		log.Println("未登录智联招聘，请先登录")
		return fmt.Errorf("未登录")
	}

	// 执行打招呼逻辑
	// if err := doSayHi(ctx); err != nil {
	// 	log.Printf("打招呼失败: %v", err)
	// 	return err
	// }

	if err := DealNewGreet(ctx); err != nil {
		log.Printf("处理新招呼失败: %v", err)
		return err
	}

	chromeManager.CloseChrome()

	return nil
}

// 执行具体的打招呼操作

func doSayHi(ctx context.Context) error {
	totalGreetings := 0
	scrollCount := 0

	// 访问目标页面
	err := runWithTimeout(ctx,
		chromedp.Navigate("https://rd6.zhaopin.com/app/recommend?tab=recommend#sortType=recommend"),
		chromedp.WaitVisible(".tr-filter-trigger"),
	)
	if err != nil {
		return fmt.Errorf("页面加载失败: %v", err)
	}
	log.Println("页面加载成功")

	// 点击筛选按钮
	err = runWithTimeout(ctx,
		chromedp.Click(".tr-filter-trigger"),
	)
	if err != nil {
		return fmt.Errorf("点击筛选按钮失败: %v", err)
	}
	log.Println("筛选按钮已点击")

	// 点击最近筛选按钮
	err = runWithTimeout(ctx,
		chromedp.Click("//div[@class='tr-talent-filter__history-con']//div[@class='tr-talent-filter__history-item'][1]"),
	)
	if err != nil {
		return fmt.Errorf("点击最近筛选按钮失败: %v", err)
	}
	log.Println("最近筛选按钮已点击")

	// 点击确定按钮
	err = runWithTimeout(ctx,
		chromedp.Click("//div[@class='km-modal__footer']//button[@zp-stat-id='rsmlist-confirm']"),
	)
	if err != nil {
		return fmt.Errorf("点击确定按钮失败: %v", err)
	}
	log.Println("确定按钮已点击")

	for totalGreetings < 20 {
		// 获取候选人信息
		var candidates []struct {
			Name     string
			CanGreet bool
		}

		err = runWithTimeout(ctx,
			chromedp.Evaluate(`
                Array.from(document.querySelectorAll('div[role="listitem"]')).map(item => {
                    const nameEl = item.querySelector('.talent-basic-info__name--inner');
                    // 检查打电话按钮是否不包含 is-disabled 类
                    const callButton = item.querySelector('.resume-button.download-cv-test .large-screen-btn:not(.is-disabled)');
                    // 如果打电话按钮可点击，则查找打招呼按钮
                    const greetButton = callButton ? 
                        item.querySelector('.cv-test-recommend-chat .large-screen-btn.km-button--filled') : null;
                    
                    const name = nameEl ? nameEl.textContent.trim() : '';
                    const canGreet = !!callButton && !!greetButton;
                    
                    console.log('候选人信息:', {
                        name,
                        hasCallButton: !!callButton,
                        hasGreetButton: !!greetButton,
                        canGreet
                    });
                    
                    return {
                        name: name,
                        canGreet: canGreet
                    };
                })
            `, &candidates),
		)
		if err != nil {
			return fmt.Errorf("获取候选人信息失败: %v", err)
		}

		// 处理当前可见的候选人
		processedCount := 0
		for i, candidate := range candidates {
			if totalGreetings >= 20 {
				break
			}

			log.Printf("开始处理第 %d 个候选人: %s (可打招呼: %v)", i+1, candidate.Name, candidate.CanGreet)

			// 每处理两个候选人后滚动页面
			if processedCount > 0 && processedCount%2 == 0 {
				err = runWithTimeout(ctx,
					chromedp.Evaluate(`
                        window.scrollTo({
                            top: window.scrollY + 300,
                            behavior: 'smooth'
                        });
                    `, nil),
					chromedp.Sleep(1500*time.Millisecond),
				)
				if err != nil {
					log.Printf("滚动页面失败: %v", err)
					continue
				}
				scrollCount++
				log.Printf("已滚动页面 %d 次", scrollCount)
			}

			// 先确保元素可见
			err = runWithTimeout(ctx,
				chromedp.ScrollIntoView(fmt.Sprintf(`div[role="listitem"]:nth-of-type(%d)`, i+1)),
				chromedp.Sleep(1*time.Second),
			)
			if err != nil {
				log.Printf("滚动到候选人失败 - %s: %v", candidate.Name, err)
				continue
			}

			if candidate.CanGreet {
				// 检查元素是否存在和可见
				var isVisible bool
				err = runWithTimeout(ctx,
					chromedp.Evaluate(fmt.Sprintf(`
                        (function() {
                            const button = document.querySelector('div[role="listitem"]:nth-of-type(%d) .cv-test-recommend-chat .large-screen-btn.km-button--filled');
                            if (!button) {
                                console.log('打招呼按钮不存在');
                                return false;
                            }
                            const rect = button.getBoundingClientRect();
                            const isVisible = rect.top >= 0 && rect.bottom <= window.innerHeight;
                            console.log('按钮位置信息:', {
                                top: rect.top,
                                bottom: rect.bottom,
                                visible: isVisible
                            });
                            return isVisible;
                        })()
                    `, i+1), &isVisible),
				)
				if err != nil {
					log.Printf("检查按钮可见性失败 - %s: %v", candidate.Name, err)
					continue
				}
				log.Printf("打招呼按钮可见性: %v", isVisible)

				// 尝试点击打招呼按钮
				err = runWithTimeout(ctx,
					chromedp.Click(fmt.Sprintf(`div[role="listitem"]:nth-of-type(%d) .cv-test-recommend-chat .large-screen-btn.km-button--filled`, i+1)),
					chromedp.Sleep(500*time.Millisecond),
				)
				if err != nil {
					log.Printf("点击打招呼失败 - %s: %v", candidate.Name, err)

					// 尝试使用JavaScript点击
					err = runWithTimeout(ctx,
						chromedp.Evaluate(fmt.Sprintf(`
                            (function() {
                                const button = document.querySelector('div[role="listitem"]:nth-of-type(%d) .cv-test-recommend-chat .large-screen-btn.km-button--filled');
                                if (button) {
                                    console.log('尝试通过JavaScript点击按钮');
                                    button.click();
                                    return true;
                                }
                                return false;
                            })()
                        `, i+1), nil),
					)
					if err != nil {
						log.Printf("JavaScript点击也失败 - %s: %v", candidate.Name, err)
					} else {
						log.Printf("通过JavaScript成功点击打招呼 - %s", candidate.Name)
						totalGreetings++
					}
				} else {
					log.Printf("成功点击打招呼 - %s", candidate.Name)
					totalGreetings++
				}
			} else {
				log.Printf("跳过 %s - 不可打招呼", candidate.Name)
			}

			processedCount++
		}
	}

	log.Printf("打招呼任务完成！共打了%d个招呼\n", totalGreetings)
	return nil
}
