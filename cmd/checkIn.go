package cmd

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/proxy"
	"io"
	"net/http"
	"net/url"
	"strings"

	"time"
)

type RespStruct struct {
	Success         bool   `json:"success"`
	Message         string `json:"message"`
	PointsEarned    int    `json:"pointsEarned"`
	CheckInCount    int    `json:"check_in_count"`
	AdditionalDraws int    `json:"additionalDraws"`
}

type RespFailStruct struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// 封装颜色输出函数
func printColoredMessage(colorFunc func(a ...interface{}) string, level, message string) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s [%s] %s\n", colorFunc(currentTime), colorFunc(level), message)
}

func CheckIn(num int, token string, proxyStr string) {
	// 创建红色打印函数
	red := color.New(color.FgRed).SprintFunc()
	// 创建绿色打印函数
	green := color.New(color.FgGreen).SprintFunc()
	// 创建黄色打印函数
	yellow := color.New(color.FgYellow).SprintFunc()

	ip, port, username, password := parseProxy(proxyStr)
	proxyAddress := "socks5://" + username + ":" + password + "@" + ip + ":" + port

	// 创建 HTTP 客户端
	client, err := newHTTPClientWithProxy(proxyAddress)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonData := "{}"
	// 创建请求
	req, err := http.NewRequest("POST", "https://api.x.ink/v1/check-in", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return
	}

	req.Header = createHeaders(token)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return
	}

	if strings.Contains(string(body), "true") {
		respModel := RespStruct{}
		json.Unmarshal(body, &respModel)
		printColoredMessage(green, "SUCCESS", fmt.Sprintf(green("账号%d签到成功"), num))
		return
	} else {
		respModel := RespFailStruct{}
		json.Unmarshal(body, &respModel)
		if respModel.Error == "Already checked in today" {
			printColoredMessage(yellow, "INFO", fmt.Sprintf(yellow("账号%d今日已签到"), num))
			return
		} else if respModel.Error == "Please follow Twitter or join Discord first" {
			printColoredMessage(red, "ERROR", fmt.Sprintf(red("账号%d请先绑定推特或DC"), num))
			return

		} else {
			printColoredMessage(red, "ERROR", fmt.Sprintf(red("账号%d签到失败，%s"), num, respModel.Error))
			return
		}
	}

	//// 定义结构体变量用于存储解析结果
	//respModel := RespStruct{}
	//
	//// 使用解码器将响应体直接解析到结构体上
	//err = decoder.Decode(&respModel)
	//if err != nil {
	//	return "", fmt.Errorf("解析 JSON 数据出错: %v\n", err)
	//}

}

func parseProxy(account string) (ip, port, username, password string) {
	// 假设 proxy 格式为 "ip:port:username:password"
	parts := strings.Split(account, ":")
	ip, port = parts[0], parts[1]
	if len(parts) > 2 {
		username = parts[2]
	}
	if len(parts) > 3 {
		password = parts[3]
	}

	return
}

func newHTTPClientWithProxy(proxyAddress string) (*http.Client, error) {
	// 解析代理地址
	proxyURL, err := url.Parse(proxyAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proxy address: %v", err)
	}
	transport := &http.Transport{}

	if proxyURL.Scheme == "socks5" {
		// 设置 SOCKS5 代理并进行身份验证
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return nil, fmt.Errorf("failed to create SOCKS5 dialer: %v", err)
		}
		// 创建 HTTP Transport 使用 SOCKS5 代理
		transport = &http.Transport{
			Dial: dialer.Dial,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略 HTTPS 错误
			},
		}
	} else {
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 忽略 HTTPS 错误
			},
		}
	}

	// 创建 HTTP 客户端
	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}

	return client, nil
}

func createHeaders(token string) http.Header {
	headers := http.Header{}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	}
	headers.Set("User-Agent", userAgents[time.Now().UnixNano()%int64(len(userAgents))])
	headers.Set("Authorization", "Bearer "+token)
	return headers
}
