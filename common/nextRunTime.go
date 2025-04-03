package common

import "time"

// NextScheduleTime 计算下一个目标时间
func NextScheduleTime(loc *time.Location) time.Time {
	now := time.Now().In(loc)

	// 构造今日 10:00
	todayTarget := time.Date(now.Year(), now.Month(), now.Day(),
		10, 0, 0, 0, loc)

	// 如果当前时间已过今日 10:00，则设置为明日 10:00
	if now.After(todayTarget) || now.Equal(todayTarget) {
		return todayTarget.Add(24 * time.Hour)
	}
	return todayTarget
}
