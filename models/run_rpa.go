package models

import (
	"kltRPA/utils"
	"log"
	"os"
	"time"
)

func RunRPA(selectedCampuses []string) {
	log.Println("开始执行RPA...")

	// 初始化 Chrome
	ctx, cancel, err := utils.InitChrome()
	if err != nil {
		log.Printf("初始化 Chrome 失败: %v", err)
		return
	}
	defer cancel()

	// 检查登录状态
	log.Println("正在检查登录状态...")
	isLoggedIn, err := TestLogin(ctx)
	if err != nil {
		log.Printf("检查登录状态失败: %v", err)
		return
	}

	if !isLoggedIn {
		log.Println("未登录智联招聘，请先登录")
		return
	}
	log.Println("已登录智联招聘")

	// 创建输出文件
	filePath := GenerateFilename("智联简历")
	if err := InitExcelFile(filePath); err != nil {
		log.Printf("初始化Excel文件失败: %v", err)
		return
	}
	log.Printf("创建输出文件: %s", filePath)

	// 开始下载简历
	err = DownloadResume(ctx, selectedCampuses, filePath)
	if err != nil {
		log.Printf("下载简历过程中出错: %v", err)
	}

	// 等待结果查看
	log.Println("处理完成，等待1秒...")
	time.Sleep(1 * time.Second)

	// 输出统计信息
	log.Printf("输出文件保存在: %s", filePath)
	if _, err := os.Stat(filePath); err == nil {
		log.Println("文件已成功保存")
	} else {
		log.Printf("警告：无法确认文件是否保存成功: %v", err)
	}

	log.Println("执行完成")
}
