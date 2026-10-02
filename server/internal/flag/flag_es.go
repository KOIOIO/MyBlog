package flag

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"server/internal/model/elasticsearch"
)

// Elasticsearch 创建 ES 索引。
func Elasticsearch(deps Deps) error {
	indexExists, err := deps.ES.Indices.Exists(elasticsearch.ArticleIndex()).Do(context.TODO())
	if err != nil {
		return err
	}

	if indexExists {
		fmt.Println("The index already exists. Do you want to delete the data and recreate the index? (y/n)")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "y":
			fmt.Println("Proceeding to delete the data and recreate the index...")
			if _, err := deps.ES.Indices.Delete(elasticsearch.ArticleIndex()).Do(context.TODO()); err != nil {
				return err
			}
		case "n":
			fmt.Println("Exiting the program.")
			os.Exit(0)
		default:
			fmt.Println("Invalid input. Please enter 'y' to delete and recreate the index, or 'n' to exit.")
			return Elasticsearch(deps)
		}
	}

	_, err = deps.ES.Indices.Create(elasticsearch.ArticleIndex()).Mappings(elasticsearch.ArticleMapping()).Do(context.TODO())
	return err
}
