package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
	"time"
	"xos/cmd"
)

func main() {
	green := color.New(color.FgGreen).SprintFunc()
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
	for i, v := range tokens {
		proxyStr := proxies[i%proxyCount]
		keys := strings.Split(v, ":")
		if len(keys) != 2 {
			fmt.Printf("账号%dtoken格式错误，跳过\n", i+1)
			continue
		}
		token := cmd.Login(i+1, keys[0], keys[1], proxyStr)
		for token == "" {
			token = cmd.Login(i+1, keys[0], keys[1], proxyStr)
		}

		currentTime := time.Now().Format("2006-01-02 15:04:05")
		fmt.Printf("%s [%s] %s\n", green(currentTime), green("SUCCESS"), fmt.Sprintf(green("账号%d登录成功，token为%s"), i+1, token[:4]+"..."+token[len(token)-4:]))

		cmd.CheckIn(i+1, token, proxyStr)
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
