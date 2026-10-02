package comment

import (
	"encoding/json"
	"testing"
	"time"

	"server/internal/domain/shared"
	"server/internal/model/appTypes"
	"server/internal/model/database"

	"github.com/gofrs/uuid"
)

var testUUID = uuid.FromStringOrNil("fbd5364d-bb15-11f1-b230-16b7a2303b52")

func TestDedupByRootUser(t *testing.T) {
	u1 := testUUID
	u2 := uuid.FromStringOrNil("11111111-1111-1111-1111-111111111111")

	// 用户 u1 的评论构成根列表：rootA(10, u1) 的子评论 c1(u2)、c2(u1)；c2 的子评论 c3(u1)、c4(u2)。
	// 规则：作为某根评论（同用户）子孙出现的评论从根列表剔除 → c2、c3 被剔除，rootA 保留。
	c4 := &Comment{ID: 4, UserUUID: u2}
	c3 := &Comment{ID: 3, UserUUID: u1, Children: []*Comment{c4}}
	c2 := &Comment{ID: 2, UserUUID: u1, Children: []*Comment{c3}}
	c1 := &Comment{ID: 1, UserUUID: u2}
	rootA := &Comment{ID: 10, UserUUID: u1, Children: []*Comment{c1, c2}}
	// 根评论 B(u2) 的子评论 d1(u1) 不受 A 影响（只看各自根）。
	d1 := &Comment{ID: 21, UserUUID: u1}
	rootB := &Comment{ID: 20, UserUUID: u2, Children: []*Comment{d1}}

	out := DedupByRootUser([]*Comment{rootA, c2, c3, rootB})
	ids := make([]uint, 0, len(out))
	for _, c := range out {
		ids = append(ids, c.ID)
	}
	// rootA 保留，c2/c3 剔除，rootB 保留（其树内无同用户子孙）
	if len(ids) != 2 || ids[0] != 10 || ids[1] != 20 {
		t.Fatalf("want roots [10 20], got %v", ids)
	}
	// 树的 children 结构不被修改：rootA.Children 仍为 [c1 c2]
	if len(rootA.Children) != 2 || rootA.Children[1].ID != 2 {
		t.Fatalf("rootA children must stay [c1 c2], got %+v", rootA.Children)
	}
}

func TestCommentJSONEqualsDatabaseModel(t *testing.T) {
	now := time.Date(2026, 9, 29, 8, 36, 26, 78000000, time.FixedZone("CST", 8*3600))
	pid := uint(55)

	// database.Comment（子评论未加载 → children nil）
	dbc := database.Comment{
		MODEL:     database.MODEL{ID: 11, CreatedAt: now, UpdatedAt: now},
		ArticleID: "15",
		PID:       &pid,
		Children:  nil,
		UserUUID:  testUUID,
		Content:   "哈哈",
		User: database.User{
			MODEL:     database.MODEL{ID: 0},
			UUID:      testUUID,
			Username:  "Xiaoyu_Wang",
			Avatar:    "/uploads/avatar/x.jpg",
			Address:   "河南省郑州市",
			Signature: "没有理想的人不伤心",
			Register:  appTypes.Email,
		},
	}

	// 对应 domain.Comment
	dc := &Comment{
		ID:        11,
		CreatedAt: now,
		UpdatedAt: now,
		ArticleID: "15",
		PID:       &pid,
		Children:  nil,
		UserUUID:  testUUID,
		Content:   "哈哈",
		User: UserBrief{
			ID:        0,
			UUID:      testUUID,
			Username:  "Xiaoyu_Wang",
			Avatar:    "/uploads/avatar/x.jpg",
			Address:   "河南省郑州市",
			Signature: "没有理想的人不伤心",
			Register:  shared.Email,
		},
	}

	a, _ := json.Marshal(dbc)
	b, _ := json.Marshal(dc)
	if string(a) != string(b) {
		t.Fatalf("JSON mismatch:\nold: %s\nnew: %s", a, b)
	}
}

func TestCommentChildrenNullVsEmpty(t *testing.T) {
	// 未加载子评论 → children null；加载后（空）→ children []
	var m map[string]interface{}
	c := &Comment{ID: 1, Children: nil}
	json.Unmarshal(mustJSON(t, c), &m)
	if _, ok := m["children"]; !ok || m["children"] != nil {
		t.Fatalf("children want null, got %v", m["children"])
	}
	c2 := &Comment{ID: 2, Children: []*Comment{}}
	json.Unmarshal(mustJSON(t, c2), &m)
	ch, ok := m["children"].([]interface{})
	if !ok || len(ch) != 0 {
		t.Fatalf("children want [], got %v", m["children"])
	}
}

func mustJSON(t *testing.T, v interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
