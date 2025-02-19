package response

type Captcha struct {
	CaptchaID string `json:"captcha_Id"`
	PicPath   string `json:"pic_Path"`
}
