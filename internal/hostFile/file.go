package hostfile

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

var HostsPath = os.Getenv("SystemRoot") + `\System32\drivers\etc\hosts`

func UpdateHostFile(hosts map[string]string) error {
	tmpFilePath := HostsPath + ".tmp"
	err := update(HostsPath, tmpFilePath, hosts)
	if err != nil {
		return err
	}
	err = os.Rename(tmpFilePath, HostsPath)
	if err != nil {
		return fmt.Errorf("Rename file raed failed! err: %v\n", err)
	}
	return nil
}

var regStarMe = regexp.MustCompile(`^\s*#\s*(Star me|Update time|Update url):`)

func update(oldFile, newFile string, hosts map[string]string) error {
	//读写方式打开文件
	file, err := os.OpenFile(oldFile, os.O_RDWR, 0666)
	if err != nil {
		return fmt.Errorf("Open file failed! err: %v\n", err)
	}
	//defer关闭文件
	defer file.Close()

	//读取文件内容到io中
	scanner := bufio.NewScanner(file)

	// 新建临时文件
	tempFile, err := os.OpenFile(newFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Temp create failed! err: %v\n", err)
	}
	defer tempFile.Close()
	writer := bufio.NewWriter(tempFile)
	_ = writer.Flush()

	for scanner.Scan() {
		// 读取当前行内容
		line := scanner.Text()
		if regStarMe.MatchString(line) {
			continue
		}
		newLines := TransLine(line, hosts)
		for _, newLine := range newLines {
			_, _ = writer.WriteString(newLine + "\r\n")
		}
	}

	for domain, ip := range hosts {
		if ip != "" {
			newLine := fmt.Sprintf("%-20s %s\r\n", ip, domain)
			_, _ = writer.WriteString(newLine)
			hosts[domain] = ""
		}
	}
	_, _ = writer.WriteString(fmt.Sprintf("# Update time: %s\r\n", time.Now().Format(time.RFC3339)))
	_, _ = writer.WriteString("# Star me: http://github.com/pkf/hosts\r\n")

	_ = writer.Flush()

	return nil
}
