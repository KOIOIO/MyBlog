package article

import (
	"encoding/json"
	"testing"
)

func TestValidateCategory(t *testing.T) {
	if err := ValidateCategory("技术"); err != nil {
		t.Fatalf("技术 must be valid: %v", err)
	}
	if err := ValidateCategory("生活"); err != nil {
		t.Fatalf("生活 must be valid: %v", err)
	}
	if err := ValidateCategory("科技"); err == nil || err.Error() != "分类只能是 技术 或 生活" {
		t.Fatalf("invalid category must be rejected with old message, got %v", err)
	}
}

func TestIllustrationsOf(t *testing.T) {
	content := "![alt](https://a.com/1.png) 正文 ![x](/uploads/2.jpg)"
	got := IllustrationsOf(content)
	if len(got) != 2 || got[0] != "https://a.com/1.png" || got[1] != "/uploads/2.jpg" {
		t.Fatalf("illustrations mismatch: %v", got)
	}
	if len(IllustrationsOf("no images")) != 0 {
		t.Fatal("no images expected")
	}
}

func TestDiffTags(t *testing.T) {
	added, removed := DiffTags([]string{"a", "b", "c"}, []string{"b", "c", "d"})
	if len(added) != 1 || added[0] != "d" {
		t.Fatalf("added mismatch: %v", added)
	}
	if len(removed) != 1 || removed[0] != "a" {
		t.Fatalf("removed mismatch: %v", removed)
	}
}

// TestArticleJSONParity 校验领域实体与 model/elasticsearch.Article 序列化等价（HTTP 兼容关键）。
func TestArticleJSONParity(t *testing.T) {
	dom := &Article{
		CreatedAt: "2026-10-02 10:00:00",
		UpdatedAt: "2026-10-02 10:00:00",
		Cover:     "/uploads/image/cover.png",
		Title:     "Gin 框架中间件机制深度解析",
		Keyword:   "Gin 框架中间件机制深度解析",
		Category:  "技术",
		Tags:      []string{"Gin", "Go"},
		Abstract:  "摘要",
		Content:   "内容",
		Views:     100,
		Comments:  3,
		Likes:     5,
	}
	legacy := struct {
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
		Cover     string   `json:"cover"`
		Title     string   `json:"title"`
		Keyword   string   `json:"keyword"`
		Category  string   `json:"category"`
		Tags      []string `json:"tags"`
		Abstract  string   `json:"abstract"`
		Content   string   `json:"content"`
		Views     int      `json:"views"`
		Comments  int      `json:"comments"`
		Likes     int      `json:"likes"`
	}{
		CreatedAt: dom.CreatedAt, UpdatedAt: dom.UpdatedAt, Cover: dom.Cover, Title: dom.Title,
		Keyword: dom.Keyword, Category: dom.Category, Tags: dom.Tags, Abstract: dom.Abstract,
		Content: dom.Content, Views: dom.Views, Comments: dom.Comments, Likes: dom.Likes,
	}
	a, _ := json.Marshal(dom)
	b, _ := json.Marshal(legacy)
	if string(a) != string(b) {
		t.Fatalf("JSON mismatch:\ndomain=%s\nlegacy=%s", a, b)
	}
}
