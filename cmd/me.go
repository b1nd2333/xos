package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

type UserRespStruct struct {
	Data struct {
		Eth                              string      `json:"eth"`
		Code                             string      `json:"code"`
		Ens                              interface{} `json:"ens"`
		Points                           float64     `json:"points"`
		DepositAmount                    int         `json:"depositAmount"`
		ReferrerCode                     string      `json:"referrer_code"`
		CurrentSquadBalance              int         `json:"currentSquadBalance"`
		CurrentSquadBalanceRewardedDraws int         `json:"CurrentSquadBalanceRewardedDraws"`
		Shares                           int         `json:"shares"`
		LastDrawTime                     time.Time   `json:"LastDrawTime"`
		LastDrawTimeETH                  int         `json:"LastDrawTimeETH"`
		LastHarvestETH                   int         `json:"LastHarvestETH"`
		LastHarvestTime                  time.Time   `json:"LastHarvestTime"`
		CurrentDraws                     int         `json:"currentDraws"`
		LastCheckIn                      time.Time   `json:"lastCheckIn"`
		CheckInCount                     int         `json:"check_in_count"`
		DepositAmountSol                 int         `json:"depositAmountSol"`
		Sol                              interface{} `json:"sol"`
		Ticket                           int         `json:"ticket"`
		CurrentSquadBalanceSol           int         `json:"currentSquadBalanceSol"`
		TicketSquad                      int         `json:"ticketSquad"`
		LastHarvestTicket                int         `json:"LastHarvestTicket"`
		TwitterInviteRewardLevel         interface{} `json:"twitter_invite_reward_level"`
		InviteCount                      int         `json:"inviteCount"`
		VerifiedTwitterInviteCount       int         `json:"verifiedTwitterInviteCount"`
		MyInviteCount                    int         `json:"myInviteCount"`
		Twitter                          struct {
			Id        string `json:"id"`
			Username  string `json:"username"`
			Avatar    string `json:"avatar"`
			HasReward bool   `json:"hasReward"`
		} `json:"twitter"`
		Discord struct {
			Id        interface{} `json:"id"`
			Username  interface{} `json:"username"`
			Avatar    interface{} `json:"avatar"`
			HasReward bool        `json:"hasReward"`
		} `json:"discord"`
		TwitterInviteReward struct {
			Level                int  `json:"level"`
			VerifiedInvites      int  `json:"verifiedInvites"`
			NextLevelAt          int  `json:"nextLevelAt"`
			CanClaimReward       bool `json:"canClaimReward"`
			TotalPossibleRewards int  `json:"totalPossibleRewards"`
		} `json:"twitterInviteReward"`
		IsInGuardianWhitelist bool `json:"isInGuardianWhitelist"`
		IsOld                 bool `json:"isOld"`
		OldPoint              int  `json:"oldPoint"`
	} `json:"data"`
}

func GetMe(token string, proxyStr string) (bool, string, float64, int, error) {
	ip, port, username, password := parseProxy(proxyStr)
	proxyAddress := "socks5://" + username + ":" + password + "@" + ip + ":" + port

	// 创建 HTTP 客户端
	client, err := newHTTPClientWithProxy(proxyAddress)
	if err != nil {
		fmt.Println(err)
		return false, "", 0, 0, err
	}

	// 创建请求
	req, err := http.NewRequest("GET", "https://api.x.ink/v1/me", nil)
	if err != nil {
		return false, "", 0, 0, err
	}

	req.Header = createHeaders(token)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return false, "", 0, 0, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return false, "", 0, 0, err
	}

	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取响应体出错: %v\n", err)
		return false, "", 0, 0, err
	}

	userRespModel := UserRespStruct{}
	json.Unmarshal(body, &userRespModel)

	if userRespModel.Data.Eth == "" {
		return false, "", 0, 0, errors.New("获取失败")
	}

	// 是否能够领取
	canClaimReward := userRespModel.Data.TwitterInviteReward.CanClaimReward
	// 邀请码
	code := userRespModel.Data.Code
	// 当前分数
	points := userRespModel.Data.Points
	// 旋转次数
	currentDraws := userRespModel.Data.CurrentDraws
	return canClaimReward, code, points, currentDraws, nil
}
