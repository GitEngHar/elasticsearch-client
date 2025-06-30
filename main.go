package main

// https://www.elastic.co/docs/reference/elasticsearch/clients/go/getting-started#_connecting
import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"log"
	"net/http"
	"os"
	"strings"
)

// Docs structure
type Tweet struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

func main() {
	// init client
	cfg := elasticsearch.Config{
		Addresses: []string{"https://localhost:9200"},
		Username:  "elastic",
		Password:  os.Getenv("ELASTIC_PASSWD"),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	es, err := elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating the client : %s", err)
	}

	// Docs apply index
	tweet := Tweet{User: "bob", Message: "Hello Elastic"}
	data, _ := json.Marshal(tweet)
	req := esapi.IndexRequest{
		Index:      "tweets",
		DocumentID: "1",
		Body:       strings.NewReader(string(data)),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Index error %s", err)
	}
	defer res.Body.Close()
	fmt.Printf("Indexed document, status: %s\n", res.Status())

	// Exec query of search
	query := `{
		"query": {
			"match": {
				"message": "Elastic"
			}
		}
	}`

	searchReq := esapi.SearchRequest{
		Index: []string{"tweets"},
		Body:  strings.NewReader(query),
	}
	searchRes, err := searchReq.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Search error: %s", err)
	}

	// View Output after parse
	var sr struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source Tweet `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(searchRes.Body).Decode(&sr); err != nil {
		log.Fatalf("Parsing search response: %s", err)
	}

	fmt.Printf("Found %d hits:\n", sr.Hits.Total.Value)
	for _, hit := range sr.Hits.Hits {
		fmt.Printf("- %s: %s \n", hit.Source.User, hit.Source.Message)
	}
}
