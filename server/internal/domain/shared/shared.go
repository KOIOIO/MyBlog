// Package shared 提供跨限界上下文共享的领域概念（枚举与基础类型）。
package shared

import "encoding/json"

// RoleID 用户角色。
type RoleID int

const (
	Guest RoleID = iota // 游客
	User                // 普通用户
	Admin               // 管理员
)

// Register 用户注册来源。
type Register int

const (
	Email Register = iota // 邮箱验证码注册
	QQ                    // QQ登录注册
)

// MarshalJSON 实现 json.Marshaler 接口（序列化为中文描述，保持既有行为）。
func (r Register) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

// UnmarshalJSON 实现 json.Unmarshaler 接口。
func (r *Register) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*r = ToRegister(str)
	return nil
}

// String 返回 Register 的字符串表示。
func (r Register) String() string {
	switch r {
	case Email:
		return "邮箱"
	case QQ:
		return "QQ"
	default:
		return "未知"
	}
}

// ToRegister 将字符串转换为 Register。
func ToRegister(str string) Register {
	switch str {
	case "邮箱":
		return Email
	case "QQ":
		return QQ
	default:
		return -1
	}
}
