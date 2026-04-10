package utils

import "strings"

// SanitizeFTSSearch 清理 MySQL 全文检索 BOOLEAN 模式的特殊字符
// + - * " @ > < ( ) ~ 在 BOOLEAN 模式下是运算符，单独使用会导致语法错误
func SanitizeFTSSearch(keyword string) string {
	replacer := strings.NewReplacer(
		"+", " ",
		"-", " ",
		"*", " ",
		"\"", " ",
		"@", " ",
		">", " ",
		"<", " ",
		"(", " ",
		")", " ",
		"~", " ",
	)
	return strings.TrimSpace(replacer.Replace(keyword))
}
