package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net/http"
	"time"
	"xos/common"
)

type LoginStruct struct {
	WalletAddress string `json:"walletAddress"`
	SignMessage   string `json:"signMessage"`
	Signature     string `json:"signature"`
	Referrer      string `json:"referrer"`
}

type LoginRespStruct struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func Login(num int, publicKey, privateKey, proxyStr string) string {
	//// 创建红色打印函数
	//red := color.New(color.FgRed).SprintFunc()
	//// 创建绿色打印函数
	//green := color.New(color.FgGreen).SprintFunc()
	// 创建黄色打印函数
	yellow := color.New(color.FgYellow).SprintFunc()
	ip, port, username, password := parseProxy(proxyStr)
	proxyAddress := "socks5://" + username + ":" + password + "@" + ip + ":" + port
	//proxyAddress = "http://127.0.0.1:8080"

	msg := getMessage(publicKey, proxyAddress)
	for msg == "" {
		printColoredMessage(yellow, "INFO", fmt.Sprintf(yellow("账号%d获取签名失败，等待三秒重新获取"), num))
		time.Sleep(3 * time.Second)
		msg = getMessage(publicKey, proxyAddress)
	}

	signature, _ := common.SignMessage(msg, privateKey)

	loginModel := &LoginStruct{}
	loginModel.WalletAddress = publicKey
	loginModel.SignMessage = msg
	loginModel.Signature = signature
	loginModel.Referrer = "YNMOLP"

	marshal, _ := json.Marshal(loginModel)

	req, err := http.NewRequest("POST", "https://api.x.ink/v1/verify-signature2", bytes.NewReader(marshal))
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 创建 HTTP 客户端
	client, err := newHTTPClientWithProxy(proxyAddress)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	}
	req.Header.Set("User-Agent", userAgents[time.Now().UnixNano()%int64(len(userAgents))])
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	//req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	//req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Connection", "keep-alive")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		//roll(num, token, proxyStr)
		return ""
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return ""
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return ""
	}

	//fmt.Println(string(body))
	loginRespStruct := &LoginRespStruct{}
	err = json.Unmarshal(body, loginRespStruct)
	if err != nil {
		return ""
	}

	if loginRespStruct.Token != "" {
		return loginRespStruct.Token
	}
	return ""
}

type MessageRespStruct struct {
	Message string `json:"message"`
}

func getMessage(publicKey string, proxyAddress string) string {
	// 获取消息
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.x.ink/v1/get-sign-message2?walletAddress=%s", publicKey), nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:54.0) Gecko/20100101 Firefox/54.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36",
	}
	req.Header.Set("User-Agent", userAgents[time.Now().UnixNano()%int64(len(userAgents))])
	req.Header.Set("Connection", "keep-alive")
	// 创建 HTTP 客户端
	client, err := newHTTPClientWithProxy(proxyAddress)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		//roll(num, token, proxyStr)
		return ""
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return ""
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return ""
	}

	messageRespModel := &MessageRespStruct{}
	err = json.Unmarshal(body, messageRespModel)
	if err != nil {
		return ""
	}

	if messageRespModel.Message == "" {
		return ""
	}
	return messageRespModel.Message
}
