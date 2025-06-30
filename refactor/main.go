package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

func main() {
	// Elasticsearch クライアントはあらかじめ生成済みとします
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{"https://localhost:9200"},
		Username:  "elastic",
		Password:  os.Getenv("ELASTIC_PASSWD"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		// TLS や Transport は省略
	})
	if err != nil {
		log.Fatalf("Error creating client: %s", err)
	}

	// シンプルな match_all クエリ
	query := `{ "query": { "match_all": {} } }`

	// Search 実行
	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("my_index"),
		es.Search.WithBody(strings.NewReader(query)),
		es.Search.WithTrackTotalHits(true), // 合計件数も取得したい場合
		es.Search.WithPretty(),             // レスポンスを見やすく整形
	)
	if err != nil {
		log.Fatalf("Search error: %s", err)
	}
	defer res.Body.Close()

	// レスポンスのパース
	var r struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Decode error: %s", err)
	}

	fmt.Printf("Found %d documents\n", r.Hits.Total.Value)
	for i, hit := range r.Hits.Hits {
		fmt.Printf("#%d: %s\n", i+1, string(hit.Source))
	}
}
