package models

import (
	"context"
	"encoding/json"
	"fmt"
	"kltRPA/utils"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/xuri/excelize/v2"
)

// ResumeInfo 简历信息结构体
type ResumeInfo struct {
	Number    int    `json:"序号"`
	Course    string `json:"意向课程"`
	Phone     string `json:"手机"`
	Campus    string `json:"校区"`
	Name      string `json:"姓名"`
	Gender    string `json:"性别"`
	Email     string `json:"邮箱"`
	Education string `json:"学历"`
	WorkYears int    `json:"工作年限"`
	JobTitle  string `json:"应聘职位"`
	Location  string `json:"居住地"`
	Status    string `json:"在职情况"`
	ResumeID  string `json:"简历编号"`
	Source    string `json:"来源"`
}

// GenerateFilename 生成文件名
func GenerateFilename(prefix string) string {
	// 确保 export_data/table-data 目录存在
	exportDir := "export_data/table-data"
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		log.Printf("创建目录失败: %v", err)
		// 如果创建目录失败，使用当前目录作为备选
		exportDir = "."
	}

	// 生成文件名
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.xlsx", prefix, timestamp)

	// 返回完整的文件路径
	return filepath.Join(exportDir, filename)
}

// InitExcelFile 初始化Excel文件
func InitExcelFile(filePath string) error {
	f := excelize.NewFile()
	headers := []string{
		"序号", "意向课程", "手机", "校区", "姓名", "性别",
		"邮箱", "学历", "工作年限", "应聘职位", "居住地",
		"在职情况", "简历编号", "来源",
	}

	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	return f.SaveAs(filePath)
}

// AppendRowToExcel 添加行到Excel
func AppendRowToExcel(info ResumeInfo, filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return err
	}
	rowNum := len(rows) + 1

	values := []interface{}{
		info.Number, info.Course, info.Phone, info.Campus,
		info.Name, info.Gender, info.Email, info.Education,
		info.WorkYears, info.JobTitle, info.Location,
		info.Status, info.ResumeID, info.Source,
	}

	for i, value := range values {
		cell := fmt.Sprintf("%c%d", 'A'+i, rowNum)
		f.SetCellValue("Sheet1", cell, value)
	}

	return f.Save()
}

// CheckResumeConditions 检查简历是否满足年龄、统招、工作经历和电话号码要求
func CheckResumeConditions(ctx context.Context, jobTitle string) (bool, string, string, error) {
	// 1. 检查年龄
	var ageText string
	err := runWithTimeout(ctx,
		chromedp.Text(".resume-basic-new__meta-item:nth-child(1)", &ageText),
	)
	if err != nil {
		if err == context.DeadlineExceeded {
			return false, "获取年龄信息超时（5秒）", "", err
		}
		return false, "获取年龄信息失败", "", err
	}

	ageMatch := regexp.MustCompile(`(\d+)岁`).FindStringSubmatch(ageText)
	if len(ageMatch) < 2 {
		return false, "无法获取年龄信息", "", nil
	}
	age, _ := strconv.Atoi(ageMatch[1])
	if age < 23 || age > 39 {
		return false, fmt.Sprintf("年龄%d岁不在23-39岁范围内", age), "", nil
	}

	// 2. 检查电话号码
	var phoneNumber string

	// 尝试直接获取电话号码
	err = runWithTimeout(ctx,
		chromedp.Sleep(500*time.Millisecond),
		chromedp.WaitVisible(".resume-basic-new__contacts--phone .is-ml-20 > div:last-child"),
		chromedp.Text(".resume-basic-new__contacts--phone .is-ml-20 > div:last-child", &phoneNumber),
	)
	if err != nil {
		return false, "无法获取电话号码", "", err
	}

	// 清理并验证电话号码
	phoneNumber = strings.TrimSpace(phoneNumber)
	phoneNumber = regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")

	// 验证是否为11位数字
	if regexp.MustCompile(`^\d{11}$`).MatchString(phoneNumber) {
		return true, "", phoneNumber, nil
	}

	// 如果格式不正确，检查是否存在且可见的"查看详情"按钮
	var hasVisibleDetailButton bool
	err = runWithTimeout(ctx,
		chromedp.Evaluate(`document.querySelector(".resume-button.get-phone button") !== null`, &hasVisibleDetailButton),
	)
	if err != nil {
		log.Printf("检查查看详情按钮时出错: %v", err)
		// 继续执行，因为这不是致命错误
	}

	if hasVisibleDetailButton {
		log.Printf("发现可见的查看详情按钮，准备点击")
		// 如果存在可见的"查看详情"按钮，点击获取完整号码
		err = runWithTimeout(ctx,
			chromedp.WaitVisible(`button.km-button.km-control.km-ripple-off.km-button--primary.km-button--text.is-ml-4`, chromedp.ByQuery),
			chromedp.Click(`button.km-button.km-control.km-ripple-off.km-button--primary.km-button--text.is-ml-4`, chromedp.ByQuery),
			chromedp.Sleep(800*time.Millisecond),
		)
		if err != nil {
			log.Printf("点击查看详情按钮失败: %v", err)
			return false, "无法获取电话号码", "", nil
		}

		// 点击确定查看按钮
		err = runWithTimeout(ctx,
			chromedp.WaitVisible("//div[contains(@class, 'get-phone-popover__content--btn')][contains(text(), '确定查看')]", chromedp.BySearch),
			chromedp.Sleep(300*time.Millisecond),
			chromedp.Click("//div[contains(@class, 'get-phone-popover__content--btn')][contains(text(), '确定查看')]", chromedp.BySearch),
		)
		if err != nil {
			log.Printf("点击确定查看按钮失败: %v", err)
			return false, "无法获取电话号码", "", nil
		}

		// 再次尝试获取电话号码
		err = runWithTimeout(ctx,
			chromedp.Sleep(800*time.Millisecond),
			chromedp.WaitVisible(".resume-basic-new__contacts--phone .is-ml-20 > div:last-child"),
			chromedp.Text(".resume-basic-new__contacts--phone .is-ml-20 > div:last-child", &phoneNumber),
		)
		if err != nil {
			return false, "无法获取电话号码", "", err
		}

		// 清理并验证电话号码
		phoneNumber = strings.TrimSpace(phoneNumber)
		phoneNumber = regexp.MustCompile(`\D`).ReplaceAllString(phoneNumber, "")

		// 验证是否为11位数字
		if !regexp.MustCompile(`^\d{11}$`).MatchString(phoneNumber) {
			return false, fmt.Sprintf("电话号码格式不正确: %s", phoneNumber), "", nil
		}
	}

	// 3. 检查是否统招
	var educationItems []string
	err = runWithTimeout(ctx,
		chromedp.Evaluate(`Array.from(document.querySelectorAll('.new-education-experiences__item')).map(item => {
			const timeText = item.querySelector('.new-education-experiences__item-time').textContent;
			const eduText = item.textContent;
			return JSON.stringify({timeText, eduText});
		})`, &educationItems),
	)
	if err != nil {
		return false, "获取教育经历失败", "", err
	}

	isRec := false
	for _, itemJSON := range educationItems {
		var item struct {
			TimeText string `json:"timeText"`
			EduText  string `json:"eduText"`
		}
		if err := json.Unmarshal([]byte(itemJSON), &item); err != nil {
			continue
		}

		if strings.Contains(item.EduText, "非统招") {
			return false, "非统招", "", nil
		}

		// 检查入学和毕业时间
		if strings.Contains(item.TimeText, " - ") {
			dates := strings.Split(item.TimeText, " - ")
			if len(dates) == 2 {
				startParts := strings.Split(dates[0], ".")
				endParts := strings.Split(dates[1], ".")
				if len(startParts) == 2 && len(endParts) == 2 {
					startYear, _ := strconv.Atoi(startParts[0])
					startMonth, _ := strconv.Atoi(startParts[1])
					endYear, _ := strconv.Atoi(endParts[0])
					endMonth, _ := strconv.Atoi(endParts[1])

					if endYear-startYear >= 3 && startMonth == 9 && (endMonth == 6 || endMonth == 7) {
						isRec = true
						break
					}
				}
			}
		}
	}

	if !isRec {
		return false, "非统招生", "", nil
	}

	// 4. 检查工作经历
	var workExperience struct {
		Titles []string `json:"titles"`
		Descs  []string `json:"descs"`
		Skills []string `json:"skills"`
	}
	err = runWithTimeout(ctx,
		chromedp.Evaluate(`(() => {
			const titles = Array.from(document.querySelectorAll('.new-work-experiences__work-job-title')).map(el => el.getAttribute('title'));
			const descs = Array.from(document.querySelectorAll('.new-work-experiences__desc .is-pre-text')).map(el => el.textContent);
			const skills = Array.from(document.querySelectorAll('.new-resume-detail-skill-tag span')).map(el => el.textContent);
			return {titles, descs, skills};
		})()`, &workExperience),
	)
	if err != nil {
		return false, "获取工作经历失败", "", err
	}

	workExperienceText := strings.Join(append(append(workExperience.Titles, workExperience.Descs...), workExperience.Skills...), " ")
	if !utils.CheckWorkExperience(jobTitle, workExperienceText) {
		return false, "工作经历不符合要求", "", nil
	}

	return true, "", phoneNumber, nil
}

// ProcessAllCampuses 处理所有校区
func ProcessAllCampuses(ctx context.Context, campuses []string, resumeNum *int, filePath string) error {
	// 最外层循环：处理每个校区
	for _, campus := range campuses {
		log.Printf("\n开始处理校区：%s", campus)
		if err := processCampus(ctx, campus, resumeNum, filePath); err != nil {
			log.Printf("处理校区 %s 失败: %v", campus, err)
			continue
		}
	}
	return nil
}

// processCampus 处理单个校区的简历
func processCampus(ctx context.Context, campus string, resumeNum *int, filePath string) error {
	// 1. 点击职位选择器
	err := runWithTimeout(ctx,
		chromedp.WaitVisible(".job-selector"),
		chromedp.Click(".job-selector"),
	)
	if err != nil {
		return fmt.Errorf("点击职位选择器失败: %v", err)
	}

	// 2. 输入校区
	err = runWithTimeout(ctx,
		chromedp.WaitVisible("//div[contains(@class, 'km-select__search-input')]//input[@class='km-input__original is-normal']"),
		chromedp.Focus("//div[contains(@class, 'km-select__search-input')]//input[@class='km-input__original is-normal']"),
		chromedp.SetValue("//div[contains(@class, 'km-select__search-input')]//input[@class='km-input__original is-normal']", ""),
		chromedp.Sleep(200*time.Millisecond),
		chromedp.SendKeys("//div[contains(@class, 'km-select__search-input')]//input[@class='km-input__original is-normal']", campus),
		chromedp.Sleep(200*time.Millisecond),
	)
	if err != nil {
		return fmt.Errorf("输入校区失败: %v", err)
	}

	// 3. 获取并遍历该校区的所有岗位
	var jobs []struct {
		Title    string `json:"title"`
		City     string `json:"city"`
		IsOnline bool   `json:"isOnline"`
	}
	time.Sleep(2 * time.Second)

	err = runWithTimeout(ctx,
		chromedp.WaitVisible(".km-select__options .job-selector__item"),
		chromedp.Evaluate(`
			Array.from(document.querySelectorAll('.km-select__options .job-selector__item')).map(item => ({
				title: item.querySelector('.job-selector__item-title').textContent,
				city: item.querySelector('.job-selector__item-city').textContent.trim(),
				isOnline: true
			}))
		`, &jobs),
	)

	if err != nil {
		return fmt.Errorf("获取岗位列表失败: %v", err)
	}

	log.Printf("\n获取到%s的岗位列表，共%d个岗位", campus, len(jobs))

	// 打印岗位列表
	for i, job := range jobs {
		status := "在线"
		if !job.IsOnline {
			status = "已下线"
		}
		log.Printf("[%d] %s - %s (%s)", i+1, job.Title, job.City, status)
	}

	// 处理每个职位
	for i, job := range jobs {
		log.Printf("\n===== 开始处理第 %d/%d 个岗位 =====", i+1, len(jobs))

		// 跳过"不限"选项
		if job.Title == "不限" {
			log.Printf("跳过'不限'选项")
			continue
		}

		// 检查是否已下线
		if !job.IsOnline {
			log.Printf("跳过已下线岗位 [%d/%d]: %s", i+1, len(jobs), job.Title)
			continue
		}

		// 验证是否匹配当前校区
		if !strings.Contains(job.City, campus) {
			log.Printf("跳过不匹配校区的岗位 [%d/%d]: %s - %s", i+1, len(jobs), job.Title, job.City)
			continue
		}

		log.Printf("\n正在处理岗位 [%d/%d]: %s, 地点: %s", i+1, len(jobs), job.Title, job.City)

		// 点击选择该职位
		err := runWithTimeout(ctx,
			chromedp.Click(fmt.Sprintf(".job-selector__item .job-selector__item-title[title='%s']", job.Title)),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			log.Printf("选择职位失败: %v", err)
			continue
		}

		// 处理该职位下的所有简历
		if err := processJob(ctx, job.Title, job.City, resumeNum, filePath); err != nil {
			log.Printf("处理职位失败: %v", err)
		}

		log.Printf("\n===== 完成处理第 %d/%d 个岗位 =====", i+1, len(jobs))

		// 重新打开职位选择器，准备处理下一个岗位
		err = runWithTimeout(ctx,
			chromedp.Click(".job-selector"),
			chromedp.Sleep(1*time.Second),
		)
		if err != nil {
			log.Printf("重新打开职位选择器失败: %v", err)
			continue
		}
	}

	return nil
}

// processJob 处理单个职位下的简历
func processJob(ctx context.Context, jobTitle, jobLocation string, resumeNum *int, filePath string) error {
	err := runWithTimeout(ctx,
		chromedp.WaitVisible(".resume-item__content"),
		chromedp.Click(".resume-item__content"),
	)
	if err != nil {
		return fmt.Errorf("没有找到简历: %v", err)
	}

	oldResumeNumber := ""
	resumeCount := 0
	for {
		resumeCount++

		// 检查简历条件
		isValid, reason, phoneNumber, err := CheckResumeConditions(ctx, jobTitle)
		if err != nil {
			log.Printf("检查简历条件失败: %v", err)
			break
		}

		if !isValid {
			log.Printf("简历不符合条件: %s", reason)
			if !clickNextResume(ctx) {
				break
			}
			continue
		}

		// 获取简历编号
		var currentURL string
		if err := runWithTimeout(ctx, chromedp.Location(&currentURL)); err != nil {
			log.Printf("获取URL失败: %v", err)
			break
		}

		re := regexp.MustCompile(`resumeNumber=([^&]+)`)
		matches := re.FindStringSubmatch(currentURL)
		if len(matches) < 2 {
			log.Printf("无法获取简历编号")
			break
		}

		resumeNumber := matches[1]
		if resumeNumber == oldResumeNumber {
			log.Printf("当前简历已经是最后一份，退出处理")
			// 关闭简历
			err := runWithTimeout(ctx,
				chromedp.WaitVisible("//div[contains(@class, 'new-shortcut-resume__close')]"),
				chromedp.Click("//div[contains(@class, 'new-shortcut-resume__close')]"),
				chromedp.Sleep(500*time.Millisecond),
			)
			if err != nil {
				log.Printf("关闭简历时出错: %v", err)
			}
			log.Printf("成功关闭简历窗口")
			break
		}

		// 提取简历信息并保存
		if err := extractAndSaveResume(ctx, jobTitle, jobLocation, phoneNumber, resumeNumber, resumeNum, filePath); err != nil {
			log.Printf("处理简历失败: %v", err)
		}

		oldResumeNumber = resumeNumber
		*resumeNum++

		if !clickNextResume(ctx) {
			break
		}
	}

	log.Printf("\n=== %s岗位共处理了 %d 份简历 ===", jobTitle, resumeCount)
	return nil
}

// clickNextResume 点击下一份简历
func clickNextResume(ctx context.Context) bool {
	err := runWithTimeout(ctx,
		chromedp.Click(".new-shortcut-resume__right"),
		chromedp.Sleep(500*time.Millisecond),
	)
	return err == nil
}

// extractAndSaveResume 提取并保存简历信息
func extractAndSaveResume(ctx context.Context, jobTitle, jobLocation, phoneNumber, resumeNumber string, resumeNum *int, filePath string) error {
	var name, workYears, education, status, email string

	// 先获取基本信息
	err := runWithTimeout(ctx,
		chromedp.Text(".resume-basic-new__name", &name),
		chromedp.Text(".resume-basic-new__meta-item:nth-child(2)", &workYears),
		chromedp.Text(".resume-basic-new__meta-item:nth-child(3)", &education),
		chromedp.Text(".resume-basic-new__meta-item:nth-child(4)", &status),
	)
	if err != nil {
		log.Printf("提取基本简历信息失败: %v", err)
	}

	// 尝试获取email
	err = runWithTimeout(ctx,
		chromedp.Text(".resume-basic-new__email .is-ml-4", &email),
	)
	if err != nil || email == "" {
		// 如果获取email失败，使用手机号代替
		email = phoneNumber
	}

	// 处理工作年限
	re := regexp.MustCompile(`\d+`)
	workYearsInt := 0
	if matches := re.FindString(workYears); matches != "" {
		workYearsInt, _ = strconv.Atoi(matches)
	}

	// 创建简历信息
	info := ResumeInfo{
		Number:    *resumeNum + 1,
		Course:    utils.GetCourse(jobTitle),
		Phone:     phoneNumber,
		Campus:    jobLocation,
		Name:      name,
		Gender:    "", // 性别暂时留空
		Email:     email,
		Education: education,
		WorkYears: workYearsInt,
		JobTitle:  jobTitle,
		Location:  "", // 居住地暂时留空
		Status:    status,
		ResumeID:  resumeNumber,
		Source:    "智联",
	}

	// 保存到Excel
	if err := AppendRowToExcel(info, filePath); err != nil {
		return fmt.Errorf("保存简历信息失败: %v", err)
	}

	// 点击"存至本地"按钮
	err = runWithTimeout(ctx,
		chromedp.WaitVisible("//div[contains(@class, 'new-resume-sidebar__actions-operate')]//div[contains(@class, 'resume-button')]//span[text()='存至本地']/parent::div/parent::div"),
		chromedp.Click("//div[contains(@class, 'new-resume-sidebar__actions-operate')]//div[contains(@class, 'resume-button')]//span[text()='存至本地']/parent::div/parent::div"),
	)
	if err != nil {
		log.Printf("点击存至本地按钮失败: %v", err)
	}

	// 等待模态框出现并点击确认
	err = runWithTimeout(ctx,
		chromedp.WaitVisible("//div[contains(@class, 'km-modal km-modal--open km-modal--v-centered km-modal--normal km-modal--no-icon km-modal--scrollable')]"),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Printf("等待模态框出现失败: %v", err)
		return err
	}

	// 点击保存按钮
	log.Printf("准备点击保存按钮...")
	err = runWithTimeout(ctx,
		chromedp.Click("//body/div[contains(@class, 'km-modal__wrapper save-resume')]/div[contains(@class, 'km-modal--open')]//div[@class='km-modal__footer']//button[contains(@class, 'km-button--primary')]"),
	)
	if err != nil {
		log.Printf("点击保存按钮失败: %v", err)
	}

	//等待模态框消失
	err = runWithTimeout(ctx,
		chromedp.WaitNotVisible("//div[contains(@class, 'km-modal km-modal--open km-modal--v-centered km-modal--normal km-modal--no-icon km-modal--scrollable')]"),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		log.Printf("等待模态框消失失败: %v", err)
		return err
	}

	return nil
}

// ClickNextResume 点击下一份简历
func ClickNextResume(ctx context.Context, reason string) (bool, error) {
	if reason != "" {
		log.Printf("%s，跳过处理", reason)
	}

	err := chromedp.Run(ctx,
		chromedp.WaitVisible(".new-shortcut-resume__right"),
		chromedp.Click(".new-shortcut-resume__right"),
		chromedp.Sleep(500*time.Millisecond),
	)
	if err != nil {
		log.Printf("点击下一份简历失败: %v", err)
		return false, err
	}
	return true, nil
}

// runWithTimeout 执行带超时的 chromedp 操作
func runWithTimeout(ctx context.Context, actions ...chromedp.Action) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := chromedp.Run(timeoutCtx, actions...)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("操作超时（5秒）: %v", err)
		}
		return err
	}
	return nil
}

// DownloadResume 下载简历主函数
func DownloadResume(ctx context.Context, selectedCampuses []string, filePath string) error {
	// 1. 访问智联招聘页面
	waitCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := runWithTimeout(waitCtx,
		chromedp.Navigate("https://rd6.zhaopin.com/app/candidate?tab=pending&jobNumber=-1&jobTitle=%E4%B8%8D%E9%99%90"),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("页面导航超时（10秒）: %v", err)
		}
		return fmt.Errorf("页面导航失败: %v", err)
	}

	// 2. 切换到"沟通中"标签页
	err = runWithTimeout(ctx,
		chromedp.WaitVisible("//div[@class='candidate-tabs']/div[@class='candidate-tabs--left']//span[text()='沟通中']"),
		chromedp.Click("//div[@class='candidate-tabs']/div[@class='candidate-tabs--left']//span[text()='沟通中']"),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		return fmt.Errorf("切换到沟通中标签页失败: %v", err)
	}

	// 3. 设置筛选条件
	if err := setFilters(ctx); err != nil {
		return fmt.Errorf("设置筛选条件失败: %v", err)
	}

	// 4. 处理所有选中的校区
	log.Printf("开始处理以下校区的简历: %v", selectedCampuses)
	resumeNum := 0
	for _, campus := range selectedCampuses {
		if err := processCampus(ctx, campus, &resumeNum, filePath); err != nil {
			log.Printf("处理校区 %s 时出错: %v", campus, err)
			continue
		}
	}

	return nil
}

// setFilters 设置筛选条件（有电话、未加标签）
func setFilters(ctx context.Context) error {
	// 筛选联系方式
	err := runWithTimeout(ctx,
		chromedp.WaitVisible(".contact-selector"),
		chromedp.Click(".contact-selector"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.WaitVisible("//div[@title='有电话']"),
		chromedp.Click("//div[@title='有电话']"),
		chromedp.Sleep(500*time.Millisecond),
	)
	if err != nil {
		return fmt.Errorf("设置电话筛选失败: %v", err)
	}

	// 筛选标签
	err = runWithTimeout(ctx,
		chromedp.WaitVisible(".resume-tag-selector"),
		chromedp.Click(".resume-tag-selector"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.WaitVisible("//div[@title='未加标签的']"),
		chromedp.Click("//div[@title='未加标签的']"),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.WaitVisible("//div[@class='candidate-filter-selector__footer']//button[@class='km-button km-control km-ripple-off km-button--primary km-button--filled is-mini']"),
		chromedp.Click("//div[@class='candidate-filter-selector__footer']//button[@class='km-button km-control km-ripple-off km-button--primary km-button--filled is-mini']"),
		chromedp.Sleep(1*time.Second),
	)
	if err != nil {
		return fmt.Errorf("设置标签筛选失败: %v", err)
	}

	return nil
}
