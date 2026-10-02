package flag

import (
	"server/internal/model/database"
)

// SQL 表结构迁移，如果表不存在，它会创建新表；如果表已经存在，它会根据结构更新表。
func SQL(deps Deps) error {
	return deps.DB.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(
		&database.Advertisement{},
		&database.ArticleCategory{},
		&database.ArticleLike{},
		&database.ArticleTag{},
		&database.Comment{},
		&database.Feedback{},
		&database.FooterLink{},
		&database.FriendLink{},
		&database.Image{},
		&database.JwtBlacklist{},
		&database.Login{},
		&database.User{},
		&database.AgentConversation{},
		&database.AgentMessage{},
		&database.AgentMemory{},
	)
}
