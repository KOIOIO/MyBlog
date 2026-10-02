// Package page 提供 MySQL 与 Elasticsearch 分页查询工具（数据库/ES 客户端经参数注入）。
package page

import (
	"context"

	"server/internal/model/other"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gorm.io/gorm"
)

// MySQLPagination 实现 MySQL 数据分页查询。
func MySQLPagination[T any](db *gorm.DB, model *T, option other.MySQLOption) (list []T, total int64, err error) {
	if option.Page < 1 {
		option.Page = 1
	}
	if option.PageSize < 1 {
		option.PageSize = 10
	}
	if option.Order == "" {
		option.Order = "id desc"
	}

	query := db.Model(model)
	if option.Where != nil {
		query = query.Where(option.Where)
	}

	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	for _, preload := range option.Preload {
		query = query.Preload(preload)
	}

	err = query.Order(option.Order).
		Limit(option.PageSize).
		Offset((option.Page - 1) * option.PageSize).
		Find(&list).Error

	return list, total, err
}

// EsPagination 实现 Elasticsearch 数据分页查询。
func EsPagination(ctx context.Context, es *elasticsearch.TypedClient, option other.EsOption) (list []types.Hit, total int64, err error) {
	if option.Page < 1 {
		option.Page = 1
	}
	if option.PageSize < 1 {
		option.PageSize = 10
	}

	from := (option.Page - 1) * option.PageSize
	option.Request.Size = &option.PageSize
	option.Request.From = &from

	res, err := es.Search().
		Index(option.Index).
		Request(option.Request).
		SourceIncludes_(option.SourceIncludes...).
		Do(ctx)
	if err != nil {
		return nil, 0, err
	}

	list = res.Hits.Hits
	total = res.Hits.Total.Value
	return list, total, nil
}
