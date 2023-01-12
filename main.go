package main

import (
	"elastic-query-service/infra/config"
	"elastic-query-service/infra/repository"
	"fmt"
)

func main() {
	config.SetUp()
	config.GetElasticClient()
	save := repository.Save()
	ids := repository.FindById(save)
	fmt.Println(ids)
	fmt.Println(len(ids))
}
