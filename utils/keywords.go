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
		"订单采购", "出入库", "供应链数据库", "采购管理", "物流配送",
	},
	"SD": {
		// "销售", "市场", "商务", "数据分析", "招商", "课程顾问",
		// "运营", "店长", "订单管理", "客户管理", "销售管理",
		// "供应商", "渠道销售", "销售内勤", "ERP", "系统运维", "招标",

		"销售", "市场", "商务", "数据分析", "招商", "课程顾问", "店长", "订单管理", "客户管理", "销售管理", "供应商", "渠道销售", "销售内勤", "ERP", "系统运维", "招标", "销售分析", "销售助理", "数据", "市场分析", "市场营销", "统计", "商务专员", "渠道", "销售订单管理", "大客户经理", "项目运营", "销售支持", "Java", "java", "贸易",
		"讲师", "培训", "内训", "外贸", "跟单", "运营",
	},
	"PP": {
		"生产计划", "车间", "生产", "生产质量", "生产管理",
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
	course := GetCourse(jobTitle)
	location := strings.TrimSpace(jobLocation)

	// 所有模块通用的校区映射
	allModulesCampusMap := map[string][]string{
		"深圳":  {"深圳"},
		"武汉":  {"武汉"},
		"东莞":  {"东莞"},
		"惠州":  {"惠州"},
		"南昌":  {"南昌"},
		"广西":  {"南宁"},
		"佛山":  {"佛山"},
		"珠海":  {"珠海"},
		"广东":  {"中山"},
		"沈阳":  {"沈阳"},
		"哈尔滨": {"哈尔滨"},
		"北京":  {"北京"},
		"天津":  {"天津"},
		"山西":  {"太原", "长治", "大同"},
		"郑州":  {"郑州", "洛阳"},
		"河北":  {"石家庄", "唐山", "廊坊", "邯郸", "保定", "承德", "秦皇岛"},
	}

	// 检查所有模块通用的校区映射
	for campus, cities := range allModulesCampusMap {
		for _, city := range cities {
			if strings.Contains(location, city) {
				return campus
			}
		}
	}

	// FICO和MM共用的校区映射
	ficoMMCampusMap := map[string][]string{
		"青岛": {"青岛", "临沂", "烟台"},
		"济南": {"济南", "潍坊", "淄博"},
		"苏州": {"苏州", "常州", "无锡"},
		"广州": {"广州", "佛山"},
		"长沙": {"长沙"},
		"上海": {"上海"},
		"厦门": {"厦门"},
	}

	if course == "FICO" || course == "MM" {
		for campus, cities := range ficoMMCampusMap {
			for _, city := range cities {
				if strings.Contains(location, city) {
					return campus
				}
			}
		}

		// FICO专属校区
		if course == "FICO" {
			ficoOnlyMap := map[string][]string{
				"合肥": {"合肥", "芜湖"},
				"重庆": {"重庆"},
				"杭州": {"杭州"},
			}
			for campus, cities := range ficoOnlyMap {
				for _, city := range cities {
					if strings.Contains(location, city) {
						return campus
					}
				}
			}
		}

		// MM专属校区
		if course == "MM" {
			mmOnlyMap := map[string][]string{
				"合肥": {"合肥", "芜湖", "阜阳"},
				"重庆": {"重庆", "成都"},
				"杭州": {"杭州"},
				"南京": {"南京"},
			}
			for campus, cities := range mmOnlyMap {
				for _, city := range cities {
					if strings.Contains(location, city) {
						return campus
					}
				}
			}
		}
	}

	// SD课程的城市直接返回
	if course == "SD" {
		sdCities := []string{"苏州", "青岛", "广州", "上海", "合肥", "重庆",
			"南京", "杭州", "厦门", "济南", "成都"}
		for _, city := range sdCities {
			if strings.Contains(location, city) {
				return city
			}
		}
	}

	// PP课程的城市直接返回
	if course == "PP" {
		if strings.Contains(location, "深圳") {
			return "深圳"
		}
		ppCities := []string{"苏州", "天津", "青岛", "广州", "武汉", "上海",
			"西安", "长沙", "北京", "长春", "合肥", "重庆",
			"南京", "成都", "杭州"}
		for _, city := range ppCities {
			if strings.Contains(location, city) {
				return city
			}
		}
	}

	// 如果没有找到匹配的校区，返回原始地点
	return location
}
