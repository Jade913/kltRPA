package utils

import (
	"log"
	"strings"
)

// CourseKeywords 课程关键词映射
var CourseKeywords = map[string][]string{
	"FICO": {
		"财务", "会计", "审计", "核算", "应收应付", "账务", "财务分析",
		"费用管理", "成本会计", "记账", "票据管理", "对账", "用友软件",
		"结算", "报税", "初级会计师", "会计从业资格证", "现金管理",
		"资金收付", "报销管理", "往来会计", "财务报表", "成本分析",
		"成本管理", "成本控制", "成本计划", "成本决策", "fico",
	},
	"MM": {
		"采购", "物料", "物流", "供应链专员", "仓储", "仓库管理",
		"供应链", "库存", "单证", "外贸分析", "物流调度",
		"订单采购", "出入库", "供应商数据库", "采购管理", "物流配送",
	},
	"SD": {
		"销售", "市场", "商务", "数据分析", "招商", "课程顾问",
		"运营", "店长", "订单管理", "客户管理", "销售管理",
		"供应商", "渠道销售", "销售内勤", "ERP", "系统运维", "招标",
	},
	"PP": {
		"生产计划", "车间", "生产运营", "生产质量", "生产管理",
		"产线", "生产物料", "生产制造", "生产工艺", "工艺流程",
		"生产统计", "生产跟单", "工艺制造", "生产技术", "物料控制",
		"生产产品", "生产设备", "工厂", "生产主管", "生产组长",
		"生产督导",
	},
}

// GetCourse 判断岗位对应什么模块
func GetCourse(jobTitle string) string {
	jobTitle = strings.ToLower(jobTitle)

	for course, keywords := range CourseKeywords {
		for _, keyword := range keywords {
			if strings.Contains(jobTitle, strings.ToLower(keyword)) {
				return course
			}
		}
	}

	return "待填写"
}

// WorkKeywordsMap 工作经历关键词映射
var WorkKeywordsMap = map[string][]string{
	"FICO": {
		"财务", "会计", "审计", "核算", "应收应付", "账务", "财务分析", "费用管理",
		"成本会计", "记账", "票据管理", "对账", "用友软件", "结算", "报税",
		"初级会计师", "会计从业资格证", "现金管理", "资金收付", "报销管理",
		"往来会计", "财务报表", "成本分析", "成本管理", "成本控制", "成本计划", "成本决策",
	},
	"MM": {
		"采购", "物料", "物流", "供应链专员", "仓储", "仓库管理", "供应链管理",
		"库存", "单证", "外贸分析", "物流调度", "订单采购", "出入库",
		"供应商数据库", "采购管理", "物流配送", "买手", "跟单", "仓库",
		"调度", "供应商",
	},
	"SD": {
		"销售分析", "销售助理", "销售运营", "销售数据", "销售内勤", "市场分析",
		"市场营销", "销售统计", "销售招投标", "商务专员", "渠道销售", "渠道运营",
		"销售订单管理", "大客户经理", "项目运营", "销售支持", "销售", "市场",
		"商务", "数据分析", "招商", "课程顾问", "运营", "店长", "订单管理",
		"客户管理", "销售管理", "供应商", "供应链", "ERP", "系统运维", "订单",
		"售前售后", "贸易", "数据", "统计", "商务", "报价", "产品", "业务",
	},
	"PP": {
		"生产计划", "生产管理", "生产统计", "pmc管理", "车间主任", "质量管理",
		"供应商管理", "车间", "生产运营", "生产质量", "产线", "生产物料",
		"物料采购", "采购", "仓库", "仓库物料", "物料计划", "生产制造",
		"生产工艺", "工艺流程", "生产跟单", "工艺制造",
	},
}

// CheckWorkExperience 检查工作经历是否包含相关关键词
func CheckWorkExperience(jobTitle, workExperienceText string) bool {
	// 确定应聘职位属于哪个模块
	module := GetCourse(jobTitle)
	if module == "待填写" {
		return false
	}

	// 获取该模块的关键词列表
	keywords, exists := WorkKeywordsMap[module]
	if !exists {
		return false
	}

	workText := strings.ToLower(workExperienceText)

	// 检查是否包含任何关键词
	for _, keyword := range keywords {
		if strings.Contains(workText, strings.ToLower(keyword)) {
			return true
		}
	}

	log.Printf("未找到任何匹配的%s模块关键词", module)
	return false
}

// 根据岗位和职位发布地获取校区
func GetCampus(jobTitle, jobLocation string) string {
	// 获取课程类型
	course := GetCourse(jobTitle)
	location := strings.TrimSpace(jobLocation)

	// SD课程
	if course == "SD" {
		// SD课程直接返回城市作为校区
		return location
	}

	// MM课程
	if course == "MM" {
		mmCampusMap := map[string][]string{
			"合肥": {"合肥", "芜湖", "阜阳"},
			"重庆": {"重庆", "成都"},
			"杭州": {"杭州"},
		}

		// 先检查MM特有的映射
		for campus, cities := range mmCampusMap {
			for _, city := range cities {
				if strings.Contains(location, city) {
					return campus
				}
			}
		}
	}

	// FICO和MM共用的校区映射
	ficoCampusMap := map[string][]string{
		"青岛": {"青岛", "临沂", "烟台"},
		"济南": {"济南", "潍坊", "淄博"},
		"长沙": {"长沙"},
		"广州": {"广州", "佛山"},
		"上海": {"上海"},
		"苏州": {"苏州", "常州", "无锡"},
	}

	// FICO专属校区
	if course == "FICO" {
		ficoOnlyMap := map[string][]string{
			"合肥": {"合肥", "芜湖"},
			"重庆": {"重庆"},
			"杭州": {"杭州"},
		}
		// 合并FICO专属校区
		for k, v := range ficoOnlyMap {
			ficoCampusMap[k] = v
		}
	}

	// 检查FICO和MM共用的校区映射
	for campus, cities := range ficoCampusMap {
		for _, city := range cities {
			if strings.Contains(location, city) {
				return campus
			}
		}
	}

	// 其他校区映射（适用于所有课程）
	otherCampusMap := map[string][]string{
		"南京":  {"南京", "无锡"},
		"广西":  {"南宁"},
		"广东":  {"深圳", "佛山", "中山", "珠海"},
		"厦门":  {"厦门"},
		"沈阳":  {"沈阳"},
		"哈尔滨": {"哈尔滨"},
		"北京":  {"北京"},
		"天津":  {"天津"},
		"山西":  {"太原", "长治", "大同"},
		"郑州":  {"郑州", "洛阳"},
		"河北":  {"石家庄", "唐山", "廊坊", "邯郸", "保定", "承德", "秦皇岛"},
	}

	// 检查通用校区映射
	for campus, cities := range otherCampusMap {
		for _, city := range cities {
			if strings.Contains(location, city) {
				return campus
			}
		}
	}

	// 如果没有找到匹配的校区，返回原始地点
	return location
}
