package service

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"golang.org/x/net/context"
	"server/global"
)

type EsService struct {
}

func (e *EsService) IndexCreate(indexName string, mapping *types.TypeMapping) error {
	_, err := global.ESClient.Indices.Create(indexName).Mappings(mapping).Do(context.TODO())
	return err
}

func (e *EsService) IndexDelete(indexName string) error {
	_, err := global.ESClient.Indices.Delete(indexName).Do(context.TODO())
	return err
}

func (e *EsService) IndexExists(indexName string) (bool, error) {
	return global.ESClient.Indices.Exists(indexName).Do(context.TODO())
}
