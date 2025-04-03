package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"xos/common"
	"xos/work"
)

func main() {
	tokens, err := readFileLines("token.txt")
	if err != nil {
		fmt.Printf("未找到token.txt文件，请确保文件存在\n")
		return
	}

	proxies, err := readFileLines("proxy.txt")
	if err != nil {
		fmt.Printf("未找到proxy.txt文件，请确保文件存在。\n")
		return
	}

	proxyCount := len(proxies)

	// 加载中国时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Fatalf("时区加载失败: %v", err)
	}

	for {
		// 获取下一个执行时间
		nextRun := common.NextScheduleTime(loc)
		waitDuration := nextRun.Sub(time.Now().In(loc))

		// 创建停止通道
		stopChan := make(chan struct{})

		// 启动倒计时显示器
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					remaining := nextRun.Sub(time.Now().In(loc)).Round(time.Second)
					if remaining <= 0 {
						return
					}
					fmt.Printf("\r下一次执行时间: %s (剩余等待: %v)   ",
						nextRun.Format("2006-01-02 15:04:05"),
						remaining)
				case <-stopChan:
					return
				}
			}
		}()

		// 等待到目标时间
		<-time.After(waitDuration)

		// 执行任务
		fmt.Println("--------------------------------------------------------------------------------------------------------")
		for i, v := range tokens {
			proxyStr := proxies[i%proxyCount]

			keys := strings.Split(v, ":")
			if len(keys) != 2 {
				fmt.Printf("账号%dtoken格式错误，跳过\n", i+1)
				continue
			}

			work.Work(i+1, keys, proxyStr)
			fmt.Println("--------------------------------------------------------------------------------------------------------")

		}
	}

}

// 读取文件内容，返回每行内容的切片
func readFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
