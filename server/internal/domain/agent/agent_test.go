package agent

import "testing"

func TestMakeTitle(t *testing.T) {
	got := MakeTitle("这是一段超过三十个字的消息内容用来验证标题截取逻辑是否正确执行")
	if len([]rune(got)) != 30 {
		t.Fatalf("want 30 runes, got %d: %q", len([]rune(got)), got)
	}
	if got := MakeTitle("短消息"); got != "短消息" {
		t.Fatalf("want original, got %q", got)
	}
}

func TestValidateAccess(t *testing.T) {
	if err := ValidateAccess(1, 1); err != nil {
		t.Fatalf("same user should pass: %v", err)
	}
	if err := ValidateAccess(1, 2); err != ErrForbidden {
		t.Fatalf("want ErrForbidden, got %v", err)
	}
}
