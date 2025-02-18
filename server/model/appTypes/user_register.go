package appTypes

import "encoding/json"

type Register int

const (
	Email Register = iota
	QQ
)

func (r Register) MarShalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Register) UnMarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*r = ToRegister(str)
	return nil
}

func (r Register) String() string {
	var str string
	switch r {
	case Email:
		str = "邮箱"
	case QQ:
		str = "qq"
	default:
		str = "未知"
	}
	return str
}

func ToRegister(str string) Register {
	switch str {
	case "邮箱":
		return Email
	case "qq":
		return QQ
	default:
		return Email
	}
}
