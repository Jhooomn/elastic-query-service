package repository

import (
	"bytes"
	"context"
	"elastic-query-service/infra/config"
	"elastic-query-service/shared/assembler"
	"elastic-query-service/shared/structs"
	"encoding/json"
	"github.com/brianvoe/gofakeit"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"log"
	"sync"
)

var elasticSearchInstance *elasticsearch.Client

const (
	maxGenerator     = 200
	elasticIndexName = "user"
)

func init() {
	elasticSearchInstance = config.GetElasticClient()
}

func Save() []string {
	ids := make([]string, 0, maxGenerator)
	var wg sync.WaitGroup
	var mu sync.RWMutex
	wg.Add(maxGenerator)
	for i := 0; i < maxGenerator; i++ {
		go func() {
			defer wg.Done()
			user := structs.User{
				Id:    uuid.NewString(),
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
				Phone: gofakeit.Phone(),
			}
			data, err := json.Marshal(user)
			if err != nil {
				log.Fatalf("Error marshaling document: %s", err)
			}
			req := esapi.IndexRequest{
				Index:      elasticIndexName,
				DocumentID: user.Id,
				Body:       bytes.NewReader(data),
				Refresh:    "true",
			}
			res, err := req.Do(context.Background(), elasticSearchInstance)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()
			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%d", res.Status(), i+1)
			} else {
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					mu.Lock()
					log.Printf("[%s] %s; version=%d; DocumentID=%s", res.Status(), r["result"], int(r["_version"].(float64)), user.Id)
					ids = append(ids, user.Id)
					mu.Unlock()
				}
			}

		}()
	}
	wg.Wait()
	return ids
}

func FindById(ids []string) []structs.User {
	usersResponse := make([]structs.User, 0, maxGenerator)
	var wg sync.WaitGroup
	var mu sync.RWMutex
	wg.Add(maxGenerator)
	for _, id := range ids {
		go func() {
			defer wg.Done()
			var buf bytes.Buffer
			query := map[string]interface{}{
				"query": map[string]interface{}{
					"match": map[string]interface{}{
						"id": id,
					},
				},
			}
			if err := json.NewEncoder(&buf).Encode(query); err != nil {
				log.Fatalf("Error encoding query: %s", err)
			}

			res, err := elasticSearchInstance.Search(
				elasticSearchInstance.Search.WithContext(context.Background()),
				elasticSearchInstance.Search.WithIndex(elasticIndexName),
				elasticSearchInstance.Search.WithBody(&buf),
				elasticSearchInstance.Search.WithTrackTotalHits(true),
				elasticSearchInstance.Search.WithPretty(),
			)
			if err != nil {
				log.Fatalf("Error getting response: %s", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				var e map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
					log.Fatalf("Error parsing the response body: %s", err)
				} else {
					log.Fatalf("[%s] %s: %s",
						res.Status(),
						e["error"].(map[string]interface{})["type"],
						e["error"].(map[string]interface{})["reason"],
					)
				}
			}
			var r map[string]interface{}
			if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
				log.Fatalf("Error parsing the response body: %s", err)
			}
			log.Printf(
				"[%s] %d hits; took: %dms",
				res.Status(),
				int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
				int(r["took"].(float64)),
			)
			for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
				mu.Lock()
				usersResponse = append(usersResponse, assembler.ElasticSearchToUser(hit))
				mu.Unlock()
			}

		}()
	}
	wg.Wait()

	return usersResponse
}

func FindByName() {

}
