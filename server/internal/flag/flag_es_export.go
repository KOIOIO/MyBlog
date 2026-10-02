package flag

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"server/internal/model/elasticsearch"
	"server/internal/model/other"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// ElasticsearchExport 导出 ES 中的数据到 JSON 文件。
func ElasticsearchExport(deps Deps) error {
	var response other.ESIndexResponse

	res, err := deps.ES.Search().
		Index(elasticsearch.ArticleIndex()).
		Scroll("1m").
		Size(1000).
		Query(&types.Query{MatchAll: &types.MatchAllQuery{}}).
		Do(context.TODO())
	if err != nil {
		return err
	}

	for _, hit := range res.Hits.Hits {
		response.Data = append(response.Data, other.Data{ID: hit.Id_, Doc: hit.Source_})
	}

	for {
		res, err := deps.ES.Scroll().ScrollId(*res.ScrollId_).Scroll("1m").Do(context.TODO())
		if err != nil {
			return err
		}
		if len(res.Hits.Hits) == 0 {
			break
		}
		for _, hit := range res.Hits.Hits {
			response.Data = append(response.Data, other.Data{ID: hit.Id_, Doc: hit.Source_})
		}
	}

	_, err = deps.ES.ClearScroll().ScrollId(*res.ScrollId_).Do(context.TODO())
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("es_%s.json", time.Now().Format("20060102"))
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	byteData, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = file.Write(byteData)
	return err
}
