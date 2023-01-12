package assembler

import "elastic-query-service/shared/structs"

func ElasticSearchToUser(hit interface{}) structs.User {
	i := hit.(map[string]interface{})["_source"]
	return structs.User{
		Id:    i.(map[string]interface{})["id"].(string),
		Name:  i.(map[string]interface{})["name"].(string),
		Email: i.(map[string]interface{})["email"].(string),
		Phone: i.(map[string]interface{})["phone"].(string),
	}
}
