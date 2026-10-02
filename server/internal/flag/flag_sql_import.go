package flag

import "strings"

// splitSQL 将 SQL 脚本按分号拆分为可执行语句（忽略空语句）。
func splitSQL(data []byte) []string {
	sqlList := strings.Split(string(data), ";")
	out := make([]string, 0, len(sqlList))
	for _, sql := range sqlList {
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		out = append(out, sql)
	}
	return out
}
