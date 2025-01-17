package logs

import (
	"log"
	"os"
)

func InitLog() {
	// 确保 logs 目录存在
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}

	// 打开或创建 app.log 文件
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}

	// 设置日志输出到文件
	log.SetOutput(file)
	log.Println("日志保存完成")
}
