package iplookup

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type result struct {
	Domain string
	Ip     string
	Err    error
}

func LookupDomains(domains []string) map[string]string {
	num := len(domains)
	hosts := make(map[string]string)
	// 需要2个管道
	// 1.job管道
	jobChan := make(chan string, num)
	// 2.结果管道
	resultChan := make(chan *result, num)

	// 循环创建job，输入到管道
	for _, domain := range domains {
		jobChan <- domain
	}

	// 3.创建工作池
	createPool(5, jobChan, resultChan)

	// 遍历结果管道
	for result := range resultChan {
		num--
		if result.Err == nil {
			hosts[result.Domain] = result.Ip
		}
		fmt.Printf("domain:%s ip:%s\n", result.Domain, result.Ip)
		if 0 == num {
			close(jobChan)
			close(resultChan)
		}
	}
	return hosts
}

// 创建工作池
// 参数1：开几个协程
func createPool(num int, jobChan chan string, resultChan chan *result) {
	for i := 0; i < num; i++ {
		go func(jobChan chan string, resultChan chan *result) {

			client := &http.Client{}

			for domain := range jobChan {
				for i := 0; i < 3; i++ {
					ip, err := lookup(client, domain)
					if err == notResolved && i < 2 {
						fmt.Println("Retry", i+1, domain)
						continue
					}
					r := &result{
						Domain: domain,
						Ip:     ip,
						Err:    err,
					}
					resultChan <- r
					break
				}
			}
		}(jobChan, resultChan)
	}
}

var ipReg = regexp.MustCompile(`<h1>IP Lookup : ([0-9.]+) \((.+?)\)</h1>`)
var ipReg1 = regexp.MustCompile(`<input name="ips" type="hidden" value="([0-9.]+?)\s`)
var notResolved = errors.New("Hostname could not be resolved to an IP address")

func lookup(client *http.Client, domain string) (string, error) {
	req, _ := http.NewRequest("POST", "https://www.ipaddress.com/ip-lookup", strings.NewReader("host="+domain))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("UserAgent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	strBody := string(body)

	rst := ipReg.FindStringSubmatch(strBody)
	if len(rst) == 0 {
		rst = ipReg1.FindStringSubmatch(strBody)
		if len(rst) == 0 {
			if strings.Contains(strBody, `Hostname could not be resolved to an IP address`) {
				return "", notResolved
			}
			return "", errors.New("not found")
		}
	}

	return rst[1], nil
}
