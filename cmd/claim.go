package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ClaimRespStruct struct {
	Success        bool   `json:"success"`
	Message        string `json:"message"`
	Error          string `json:"error"`
	DrawsAdded     int    `json:"drawsAdded"`
	NewLevel       int    `json:"newLevel"`
	CurrentInvites int    `json:"currentInvites"`
}

func Claim(token string, proxyStr string) (bool, string, error) {
	ip, port, username, password := parseProxy(proxyStr)
	proxyAddress := "socks5://" + username + ":" + password + "@" + ip + ":" + port

	// 创建 HTTP 客户端
	client, err := newHTTPClientWithProxy(proxyAddress)
	if err != nil {
		fmt.Println(err)
		return false, "", err
	}

	// 创建请求
	req, err := http.NewRequest("POST", "https://api.x.ink/v1/claim-twitter-invite-reward", strings.NewReader("{}"))
	if err != nil {
		return false, "", err
	}

	req.Header = createHeaders(token)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return false, "", errors.New(strconv.Itoa(resp.StatusCode))
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return false, "", err
	}

	claimRespModel := &ClaimRespStruct{}
	json.Unmarshal(body, claimRespModel)

	if claimRespModel.Success {
		return true, claimRespModel.Message, nil
	} else {
		return false, claimRespModel.Error, nil
	}
}
