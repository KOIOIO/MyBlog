package forum

import (
	"encoding/json"
	"testing"
	"time"

	"server/internal/domain/shared"
	"server/internal/model/appTypes"
	"server/internal/model/database"

	"github.com/gofrs/uuid"
)

func TestForumPostJSONEqualsDatabaseModel(t *testing.T) {
	now := time.Date(2026, 10, 1, 10, 0, 0, 0, time.FixedZone("CST", 8*3600))
	uu := uuid.FromStringOrNil("fbd5364d-bb15-11f1-b230-16b7a2303b52")

	old := database.ForumPost{
		MODEL:        database.MODEL{ID: 3, CreatedAt: now, UpdatedAt: now},
		UserID:       1,
		Title:        "测试帖",
		Content:      "内容",
		Category:     "技术",
		Tags:         database.JSONStringArray{"Go", "Gin"},
		Images:       database.JSONStringArray{},
		LikeCount:    2,
		CommentCount: 1,
		ViewCount:    9,
		User: database.User{
			MODEL:     database.MODEL{ID: 1},
			UUID:      uu,
			Username:  "Xiaoyu_Wang",
			Avatar:    "/uploads/avatar/x.jpg",
			Address:   "河南省郑州市",
			Signature: "没有理想的人不伤心",
			Register:  appTypes.Email,
		},
	}

	dp := &ForumPost{
		ID:           3,
		CreatedAt:    now,
		UpdatedAt:    now,
		UserID:       1,
		Title:        "测试帖",
		Content:      "内容",
		Category:     "技术",
		Tags:         []string{"Go", "Gin"},
		Images:       []string{},
		LikeCount:    2,
		CommentCount: 1,
		ViewCount:    9,
		User: UserBrief{
			ID:        1,
			UUID:      uu,
			Username:  "Xiaoyu_Wang",
			Avatar:    "/uploads/avatar/x.jpg",
			Address:   "河南省郑州市",
			Signature: "没有理想的人不伤心",
			Register:  shared.Email,
		},
	}

	a, _ := json.Marshal(old)
	b, _ := json.Marshal(dp)
	if string(a) != string(b) {
		t.Fatalf("JSON mismatch:\nold: %s\nnew: %s", a, b)
	}
}

func TestAllowedCategory(t *testing.T) {
	for _, c := range []string{"技术", "生活"} {
		if _, ok := AllowedCategory[c]; !ok {
			t.Fatalf("category %q must be allowed", c)
		}
	}
	if _, ok := AllowedCategory["娱乐"]; ok {
		t.Fatal("category 娱乐 must not be allowed")
	}
}
