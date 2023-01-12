package config

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

var elasticSearchInstance *elasticsearch.Client
var elasticSearchInstanceSingleton bool

func init() {
	ElasticSearchClient()
}

func SetUp() {
	// ...
}

func ElasticSearchClient() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	log.Println(elasticsearch.Version)

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()
	log.Println(res)
	elasticSearchInstance = es
	elasticSearchInstanceSingleton = true
}

func GetElasticClient() *elasticsearch.Client {
	if elasticSearchInstanceSingleton {
		return elasticSearchInstance
	}
	panic("there is not a elastic client")
}
