package flag

import (
	"context"
	"encoding/json"
	"os"

	"server/internal/model/elasticsearch"
	"server/internal/model/other"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/bulk"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/refresh"
)

// ElasticsearchImport 从指定的 JSON 文件导入数据到 ES。
func ElasticsearchImport(deps Deps, jsonPath string) (int, error) {
	byteData, err := os.ReadFile(jsonPath)
	if err != nil {
		return 0, err
	}

	var response other.ESIndexResponse
	if err := json.Unmarshal(byteData, &response); err != nil {
		return 0, err
	}

	indexExists, err := deps.ES.Indices.Exists(elasticsearch.ArticleIndex()).Do(context.TODO())
	if err != nil {
		return 0, err
	}
	if indexExists {
		if _, err := deps.ES.Indices.Delete(elasticsearch.ArticleIndex()).Do(context.TODO()); err != nil {
			return 0, err
		}
	}
	if _, err := deps.ES.Indices.Create(elasticsearch.ArticleIndex()).Mappings(elasticsearch.ArticleMapping()).Do(context.TODO()); err != nil {
		return 0, err
	}

	var request bulk.Request
	for _, data := range response.Data {
		request = append(request, types.OperationContainer{Index: &types.IndexOperation{Id_: data.ID}})
		request = append(request, data.Doc)
	}

	_, err = deps.ES.Bulk().
		Request(&request).
		Index(elasticsearch.ArticleIndex()).
		Refresh(refresh.True).
		Do(context.TODO())
	if err != nil {
		return 0, err
	}

	return len(response.Data), nil
}
