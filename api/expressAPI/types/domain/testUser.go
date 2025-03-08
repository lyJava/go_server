package domain

import (
	"database/sql/driver"
	"fmt"
	"log"
	"time"
)

// TestUser 测试用户信息
//
// 如果Birthday使用字符串，需要返回yyyy-MM-dd直接截取前十位即可或者启用func (u *TestUser) MarshalJSON()，前端传入yyyy-MM-dd就行，
// 如果是time.Time类型，前端传入需要格式化为RFC3339，这里是apifox的mock设置{{$date.anytime|format('yyyy-MM-dd')|formatRFC3339}}
// 如果使用自定义时间LocalDate，则需要实现UnmarshalJSON，MarshalJSON，Value与Scan四个方法，前端只需要传入yyyy-MM-dd字符串格式即可
type TestUser struct {
	Id       int64      `json:"id,omitempty"`        // 主键ID
	Username string     `json:"username,omitempty"`  // 用户名
	Password string     `json:"password,omitempty"`  // 密码
	Email    string     `json:"email,omitempty"`     // 邮箱
	Birthday *LocalDate `json:"birthday,omitempty" ` // 出生日期
	Phone    string     `json:"phone,omitempty"`     // 手机
	Address  string     `json:"address,omitempty"`   // 地址
}

// MarshalJSON 1自定义序列号json时候转换日期"birthday": "2018-05-15T00:00:00Z",为yyyy-MM-dd
/*func (u *TestUser) MarshalJSON() ([]byte, error) {
	type Alias TestUser
	return json.Marshal(&struct {
		*Alias
		Birthday string `json:"birthday,omitempty"`
	}{
		Alias:    (*Alias)(u),
		Birthday: u.Birthday.Format("2006-01-02"),
	})
}*/

// 日期格式化
const dateFormat = "2006-01-02"

// LocalDate 自定义日期 下边四个函数保证序列化与反序列化的返回格式，这里会覆盖默认的格式
type LocalDate struct {
	time.Time
}

func (t *LocalDate) UnmarshalJSON(data []byte) (err error) {
	if len(data) == 2 {
		*t = LocalDate{Time: time.Time{}}
		return
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now, err := time.ParseInLocation(`"`+dateFormat+`"`, string(data), loc)
	*t = LocalDate{Time: now}
	return
}

func (t LocalDate) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	formatted := fmt.Sprintf("\"%s\"", t.Format(dateFormat))
	return []byte(formatted), nil
}

func (t LocalDate) Value() (driver.Value, error) {

	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	log.Printf("测试用户生日Value()===%v", t.Time)
	return t.Time, nil
}

// Scan value of time.Time
func (t *LocalDate) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	log.Printf("测试用户生日Scan()===%v", value)
	if ok {
		*t = LocalDate{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to date", v)
}
