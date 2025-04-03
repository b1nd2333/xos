package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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
