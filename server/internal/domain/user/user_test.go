package user

import (
	"encoding/json"
	"testing"
	"time"

	"server/internal/domain/shared"
	"server/internal/model/appTypes"
	"server/internal/model/database"

	"github.com/gofrs/uuid"
)

func mustUUID(t *testing.T) uuid.UUID {
	t.Helper()
	u, err := uuid.NewV4()
	if err != nil {
		t.Fatal(err)
	}
	return u
}

func TestUserFreezeRules(t *testing.T) {
	u := &User{}
	if u.IsFrozen() {
		t.Fatal("new user must not be frozen")
	}
	u.SetFrozen(true)
	if !u.IsFrozen() {
		t.Fatal("user should be frozen after SetFrozen(true)")
	}
	u.SetFrozen(false)
	if u.IsFrozen() {
		t.Fatal("user should be unfrozen after SetFrozen(false)")
	}
}

func TestUserRoleRules(t *testing.T) {
	normal := &User{RoleID: shared.User}
	if !normal.IsNormalUser() || normal.IsAdmin() {
		t.Fatal("RoleID=User should be normal user only")
	}
	admin := &User{RoleID: shared.Admin}
	if !admin.IsAdmin() || admin.IsNormalUser() {
		t.Fatal("RoleID=Admin should be admin only")
	}
}

func TestUserResetRegisterInfo(t *testing.T) {
	u := &User{Username: "wwy", Email: "a@b.c"}
	uid := mustUUID(t)
	u.ResetRegisterInfo(uid, "hash", "/image/avatar.jpg")

	if u.UUID != uid {
		t.Fatalf("uuid not set: %v", u.UUID)
	}
	if u.Password != "hash" {
		t.Fatalf("password hash not set: %q", u.Password)
	}
	if u.Avatar != "/image/avatar.jpg" {
		t.Fatalf("avatar not set: %q", u.Avatar)
	}
	if u.RoleID != shared.User {
		t.Fatalf("role must be User, got %v", u.RoleID)
	}
	if u.Register != shared.Email {
		t.Fatalf("register must be Email, got %v", u.Register)
	}
}

func TestUserSetPasswordHash(t *testing.T) {
	u := &User{}
	u.SetPasswordHash("new-hash")
	if u.Password != "new-hash" {
		t.Fatalf("password not set: %q", u.Password)
	}
}

// TestUserJSONParity 校验领域实体序列化与 database.User 完全一致（HTTP 兼容性关键）。
func TestUserJSONParity(t *testing.T) {
	uid, _ := uuid.FromString("fbd5364d-bb15-11f1-b230-16b7a2303b52")
	now := time.Date(2026, 10, 2, 12, 0, 0, 0, time.UTC)

	dom := &User{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
		UUID:      uid,
		Username:  "Xiaoyu_Wang",
		Password:  "secret-hash",
		Email:     "2652777599@qq.com",
		Avatar:    "/image/avatar.jpg",
		Address:   "郑州",
		Signature: "hello",
		RoleID:    shared.Admin,
		Register:  shared.Email,
		Freeze:    false,
	}

	db := database.User{
		MODEL:     database.MODEL{ID: 1, CreatedAt: now, UpdatedAt: now},
		UUID:      uid,
		Username:  "Xiaoyu_Wang",
		Password:  "secret-hash",
		Email:     "2652777599@qq.com",
		Avatar:    "/image/avatar.jpg",
		Address:   "郑州",
		Signature: "hello",
		RoleID:    appTypes.Admin,
		Register:  appTypes.Email,
		Freeze:    false,
	}

	domJSON, _ := json.Marshal(dom)
	dbJSON, _ := json.Marshal(db)
	if string(domJSON) != string(dbJSON) {
		t.Fatalf("JSON mismatch:\ndomain=%s\ndb    =%s", domJSON, dbJSON)
	}
}
