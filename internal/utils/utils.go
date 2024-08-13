package utils

import (
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"strings"
	"time"
)

// GetUUIDStr 获取唯一字符串
func GetUUIDStr() string {
	u1 := uuid.NewV4()
	return strings.Replace(u1.String(), "-", "", -1)
}

// GetLocalDayStart 获取当天开始时间
func GetLocalDayStart() (time.Time, error) {
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstSh)
	nowStr := now.Format("2006-01-02")
	localDayStart, err := time.ParseInLocation("2006-01-02", nowStr, cstSh)
	if err != nil {
		return time.Time{}, err
	}
	return localDayStart, nil
}

// GetNowDateStart 获取当天开始时间
func GetNowDateStart() (time.Time, error) {
	now := time.Now()
	nowStr := now.Format("2006-01-02")
	localDayStart, err := time.ParseInLocation("2006-01-02", nowStr, time.Local)
	if err != nil {
		return time.Time{}, err
	}
	return localDayStart, nil
}

// GetLocalWeekStart 获取周天开始时间
func GetLocalWeekStart() (time.Time, error) {
	// 计算起始时间
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstSh)
	nowStr := now.Format("2006-01-02")
	localDayStart, err := time.ParseInLocation("2006-01-02", nowStr, cstSh)
	if err != nil {
		return time.Time{}, err
	}
	weekDay := localDayStart.Weekday()
	if weekDay == 0 {
		weekDay = 6
	} else {
		weekDay -= 1
	}
	dailyStart := localDayStart.AddDate(0, 0, -int(weekDay))
	return dailyStart, nil
}

// GetLocalDayByStartHour 获取当天日期开始时间
func GetLocalDayByStartHour(hour int) (time.Time, error) {
	dayStart, err := GetLocalDayStart()
	if err != nil {
		return time.Time{}, err
	}
	dayStart = dayStart.Add(time.Hour * time.Duration(hour))
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstSh)

	if now.Hour() < hour {
		// 计入上一天
		return dayStart.AddDate(0, 0, -1), nil

	}
	return dayStart, nil
}

// GetLocalDayStr 获取当天日期字符串
func GetLocalDayStr() string {
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstSh)
	nowStr := now.Format("2006-01-02")
	return nowStr
}

// GetNowDateStr 获取当天日期字符串
func GetNowDateStr() string {
	now := time.Now()
	nowStr := now.Format("2006-01-02")
	return nowStr
}

// GetLocalYesterdayStr 获取昨天天日期字符串
func GetLocalYesterdayStr() string {
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstSh).Add(time.Hour * -24)
	nowStr := now.Format("2006-01-02")
	return nowStr
}

// GetLocalDayStrByStartHour 获取当天日期字符串
func GetLocalDayStrByStartHour(hour int) (string, error) {
	todayStr := GetLocalDayStr()
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstSh)

	dataStart, err := GetLocalDayStart()
	if err != nil {
		return "", err
	}
	if now.Hour() < hour {
		// 计入上一天
		todayStr = dataStart.AddDate(0, 0, -1).Format("2006-01-02")
	}
	return todayStr, nil
}

// GetLocalWeekStrByStartHour 获取当周日期字符串
func GetLocalWeekStrByStartHour(hour int) (string, error) {
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstSh)
	nowStr := now.Format("2006-01-02")

	weekStartTime, err := GetLocalWeekStart()
	if err != nil {
		return "", err
	}
	weekStr := weekStartTime.Format("2006-01-02")
	if now.Hour() < hour && nowStr == weekStr {
		// 计入上一周
		weekStartTime = weekStartTime.AddDate(0, 0, -7)
	}
	return weekStartTime.Format("2006-01-02"), nil
}

// GetLocalWeekByStartHour 获取当周日期开始时间
func GetLocalWeekByStartHour(hour int) (time.Time, error) {
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Now().In(cstSh)
	nowStr := now.Format("2006-01-02")

	weekStartTime, err := GetLocalWeekStart()
	if err != nil {
		return time.Time{}, err
	}
	weekStr := weekStartTime.Format("2006-01-02")
	if now.Hour() < hour && nowStr == weekStr {
		// 计入上一周
		weekStartTime = weekStartTime.AddDate(0, 0, -7)
	}
	return weekStartTime.Add(time.Hour * time.Duration(hour)), nil
}

// GetUnixLocalWeekStart 获取周天开始时间
func GetUnixLocalWeekStart(unixAt int64) (time.Time, error) {
	// 计算起始时间
	cstSh := time.FixedZone("CST", 8*3600)
	now := time.Unix(unixAt, 0).In(cstSh)
	nowStr := now.Format("2006-01-02")
	localDayStart, err := time.ParseInLocation("2006-01-02", nowStr, cstSh)
	if err != nil {
		return time.Time{}, err
	}
	weekDay := localDayStart.Weekday()
	if weekDay == 0 {
		weekDay = 6
	} else {
		weekDay -= 1
	}
	dailyStart := localDayStart.AddDate(0, 0, -int(weekDay))
	return dailyStart, nil
}

func IndexOfInt64(l []int64, v int64) int {
	var index = -1
	for i, j := range l {
		if j == v {
			index = i
		}
	}
	return index
}

func CountInt64(l []int64, v int64) int64 {
	var count int64
	for _, j := range l {
		if j == v {
			count++
		}
	}
	return count
}

func MaxInt64(l ...int64) int64 {
	var max int64
	for _, j := range l {
		if j > max {
			max = j
		}
	}
	return max
}

// GenerateRandomAvatar 生成随机头像
func GenerateRandomAvatar() string {
	return fmt.Sprintf("https://www.gravatar.com/avatar/%d?d=monsterid", rand.Int63n(10000))
}

type Choice struct {
	Key    string
	Weight int
}

func WeightedChoice(choices []Choice) (Choice, error) {
	var ret Choice
	sum := 0
	for _, c := range choices {
		sum += c.Weight
	}
	r := rand.Intn(sum)
	for _, c := range choices {
		r -= c.Weight
		if r < 0 {
			return c, nil
		}
	}
	err := errors.New("Internal error - code should not reach this point")
	return ret, err
}
