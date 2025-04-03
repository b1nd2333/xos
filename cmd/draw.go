package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type DrawRespStruct struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	PointsEarned int    `json:"pointsEarned"`
}

func Draw(token string, proxyStr string) (string, error) {
	ip, port, username, password := parseProxy(proxyStr)
	proxyAddress := "socks5://" + username + ":" + password + "@" + ip + ":" + port

	// 创建 HTTP 客户端
	client, err := newHTTPClientWithProxy(proxyAddress)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", "https://api.x.ink/v1/draw", strings.NewReader("{}"))
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header = createHeaders(token)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
		return "", errors.New(resp.Status)
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return "", err
	}

	drawRespModel := DrawRespStruct{}
	json.Unmarshal(body, &drawRespModel)
	if drawRespModel.Success {
		return drawRespModel.Message, nil
	} else {
		return drawRespModel.Message, errors.New("转圈失败")
	}
}
