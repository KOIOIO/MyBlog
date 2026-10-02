package flag

import (
	"errors"
	"fmt"
	"os"
	"syscall"

	"server/internal/common/crypto"
	"server/internal/model/appTypes"
	"server/internal/model/database"

	"github.com/gofrs/uuid"
	"golang.org/x/term"
)

// Admin 用于创建一个管理员用户。
func Admin(deps Deps) error {
	var user database.User

	fmt.Print("Enter email: ")
	var email string
	if _, err := fmt.Scanln(&email); err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}
	user.Email = email

	fd := int(syscall.Stdin)
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer term.Restore(fd, oldState)

	fmt.Print("Enter password: ")
	password, err := readPassword()
	fmt.Println()
	if err != nil {
		return err
	}

	fmt.Print("Confirm password: ")
	rePassword, err := readPassword()
	fmt.Println()
	if err != nil {
		return err
	}

	if password != rePassword {
		return errors.New("passwords do not match")
	}
	if len(password) < 8 || len(password) > 20 {
		return errors.New("password length should be between 8 and 20 characters")
	}

	user.UUID = uuid.Must(uuid.NewV4())
	user.Username = deps.Cfg.Website.Name
	user.Password = crypto.BcryptHash(password)
	user.RoleID = appTypes.Admin
	user.Avatar = "/image/avatar.jpg"
	user.Address = deps.Cfg.Website.Address

	if err := deps.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

// readPassword 用于读取密码并且避免回显。
func readPassword() (string, error) {
	var password string
	var buf [1]byte
	for {
		_, err := os.Stdin.Read(buf[:])
		if err != nil {
			return "", err
		}
		char := buf[0]
		if char == '\n' || char == '\r' {
			break
		}
		password += string(char)
	}
	return password, nil
}
