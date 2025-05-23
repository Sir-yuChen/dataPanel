package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// https://juejin.cn/post/6844904114699108365 参考 解决时间相关问题
const TimeFormat = "2006-01-02 15:04:05"

type LocalTime time.Time

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalTime(time.Time{})
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now, err := time.ParseInLocation(`"`+TimeFormat+`"`, string(data), loc)
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(TimeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, TimeFormat)
	b = append(b, '"')
	return b, nil
}

func (t LocalTime) Value() (driver.Value, error) {
	if t.String() == "0001-01-01 00:00:00" {
		return nil, nil
	}
	return []byte(time.Time(t).Format(TimeFormat)), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if v == nil {
		*t = LocalTime(time.Time{}) // 处理 NULL 值
		return nil
	}
	// 类型安全转换
	switch vt := v.(type) {
	case time.Time:
		*t = LocalTime(vt)
		return nil
	case []byte:
		// 处理数据库驱动返回的字节流（如 MySQL 的 DATETIME）
		parsedTime, err := time.Parse(TimeFormat, string(vt))
		if err != nil {
			return fmt.Errorf("parsing time from bytes: %w", err)
		}
		*t = LocalTime(parsedTime)
		return nil
	case string:
		// 处理字符串类型输入
		parsedTime, err := time.Parse(TimeFormat, vt)
		if err != nil {
			return fmt.Errorf("parsing time from string: %w", err)
		}
		*t = LocalTime(parsedTime)
		return nil
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
}

func (t LocalTime) String() string {
	return time.Time(t).Format(TimeFormat)
}
func (t LocalTime) ToTime() time.Time {
	return time.Time(t)
}
