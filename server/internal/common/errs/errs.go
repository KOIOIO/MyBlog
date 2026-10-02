// Package errs 提供跨限界上下文共享的错误定义与包装工具（无业务语义）。
package errs

import (
	"errors"
	"fmt"
)

// 通用错误。
var (
	// ErrNotFound 表示目标资源不存在。
	ErrNotFound = errors.New("record not found")
	// ErrConflict 表示资源冲突（如重复创建）。
	ErrConflict = errors.New("record conflict")
)

// Wrap 包装错误并附带上下文信息。
func Wrap(err error, format string, args ...interface{}) error {
	return fmt.Errorf(format+": %w", append(args, err)...)
}
