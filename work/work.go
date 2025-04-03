package work

import (
	"fmt"
	"github.com/fatih/color"
	"time"
	"xos/cmd"
)

func Work(num int, keys []string, proxyStr string) {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// 登录账号
	token := cmd.Login(num, keys[0], keys[1], proxyStr)
	for token == "" {
		token = cmd.Login(num, keys[0], keys[1], proxyStr)
	}

	// 获取个人信息
	canClaimReward, code, points, currentDraws, err := cmd.GetMe(token, proxyStr)
	for err != nil {
		canClaimReward, code, points, currentDraws, err = cmd.GetMe(token, proxyStr)
	}

	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s [%s] %s\n", green(currentTime), green("SUCCESS"), fmt.Sprintf(green("账号%d登录成功，token为%s，邀请码：%s，当前分数：%f"), num, token[:4]+"..."+token[len(token)-4:], code, points))

	// 查看是否可以领取旋转
	if canClaimReward {
		fmt.Printf("%s [%s] %s\n", yellow(currentTime), green("INFO"), fmt.Sprintf(yellow("账号%d存在可领取旋转次数"), num))
		success, msg, err := cmd.Claim(token, proxyStr)
		for err != nil {
			success, msg, err = cmd.Claim(token, proxyStr)
		}
		if success {
			fmt.Printf("%s [%s] %s\n", green(currentTime), green("SUCCESS"), fmt.Sprintf(green("账号%d领取旋转次数成功,%s"), num, msg))
		} else {
			fmt.Printf("%s [%s] %s\n", yellow(currentTime), green("INFO"), fmt.Sprintf(yellow("账号%d领取旋转次数失败,%s"), num, msg))
		}
	}

	_, _, _, currentDraws, err = cmd.GetMe(token, proxyStr)
	for err != nil {
		_, _, _, currentDraws, err = cmd.GetMe(token, proxyStr)
	}

	fmt.Printf("%s [%s] %s\n", green(currentTime), green("SUCCESS"), fmt.Sprintf(green("账号%d当前可用旋转次数%d"), num, currentDraws))

	// 开始旋转
	for k := 0; k < currentDraws; k++ {
		drawMsg, err := cmd.Draw(token, proxyStr)
		for err != nil && err.Error() != "转圈失败" {
			drawMsg, err = cmd.Draw(token, proxyStr)
		}
		if err != nil {
			fmt.Printf("%s [%s] %s\n", yellow(currentTime), yellow("INFO"), fmt.Sprintf(yellow("账号%d第%d次旋转失败，%s"), num, k+1, drawMsg))
		} else {
			fmt.Printf("%s [%s] %s\n", green(currentTime), green("SUCCESS"), fmt.Sprintf(green("账号%d第%d次旋转完成，%s"), num, k+1, drawMsg))
		}
	}

	cmd.CheckIn(num, token, proxyStr)
	_, _, points, _, err = cmd.GetMe(token, proxyStr)
	for err != nil {
		_, _, points, _, err = cmd.GetMe(token, proxyStr)
	}
	fmt.Printf("%s [%s] %s\n", green(currentTime), green("SUCCESS"), fmt.Sprintf(green("账号%d本次领取完毕，当前分数%f"), num, points))
}
