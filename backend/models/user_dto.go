package models

type UserDTO struct {
	User
	Solved       int64                       `json:"solved"`                          // 用户已经解决的题目数量
	Heatmaps     map[string]map[string]int   `json:"heatmaps"`                        // 热力图 year -> {"MM-DD": count}
	PastYearHeatmap map[string]int           `json:"past_year_heatmap"`               // 近一年（滚动365天）"YYYY-MM-DD": count
	Streak       int                         `json:"streak"`                          // 连续打卡天数
	LastActive   string                      `json:"last_active"`                     // 最后活跃日期
}
