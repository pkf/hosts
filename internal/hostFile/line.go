package hostfile

import (
	"fmt"
	"net"
	"strings"
)

// 根据新的域名IP映射来转换原来的host行，如果原来是多域名单行，可能会返回多行
func TransLine(raw string, hosts map[string]string) []string {
	result := make([]string, 0)
	if strings.HasPrefix(strings.TrimSpace(raw), "#") { // 注释行直接返回
		return append(result, raw)
	}

	var line, comment string
	if strings.Contains(raw, "#") { // 行包含注释
		commentSplit := strings.Split(raw, "#")
		line = commentSplit[0]
		comment = "#" + commentSplit[1]
	} else {
		line = raw
	}

	fields := strings.Fields(line)
	if len(fields) == 0 { // 应该是空白行
		return append(result, raw)
	}

	rawIP := fields[0]
	if net.ParseIP(rawIP) == nil { // 格式非法行，不处理
		return append(result, raw)
	}

	fields = fields[1:]
	if len(fields) == 0 { // 只有IP
		return append(result, raw)
	}

	oldFileds := make([]string, 0)
	for _, domain := range fields {
		newIp, ok := hosts[domain]
		if ok {
			if newIp != "" {
				result = append(result, fmt.Sprintf("%-20s %s", newIp, domain))
				hosts[domain] = ""
			}

		} else {
			oldFileds = append(oldFileds, domain)
		}
	}
	if len(oldFileds) > 0 {
		if len(comment) > 0 {
			result = append(result, fmt.Sprintf("%-20s %s %s", rawIP, strings.Join(oldFileds, " "), comment))
		} else {
			result = append(result, fmt.Sprintf("%-20s %s", rawIP, strings.Join(oldFileds, " ")))
		}

	}

	return result
}
