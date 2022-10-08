package elastic

// ParseDocument
func ParseDocument(source map[string]interface{}) (map[string]interface{}, error) {
	doc := source["_source"].(map[string]interface{})
	return doc, nil
}

// ParseDocuments
func ParseDocuments(res map[string]interface{}) (map[string]interface{}, error) {
	docs := make([]map[string]interface{}, 0)
	for _, hit := range res["hits"].(map[string]interface{})["hits"].([]interface{}) {
		docs = append(docs, hit.(map[string]interface{})["_source"].(map[string]interface{}))
	}
	data := make(map[string]interface{})
	data["data"] = docs
	data["total"] = res["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	return data, nil
}

// ParserAggregation
func ParserAggregation(res map[string]interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	data["data"] = res["aggregations"]
	data["total"] = res["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)
	return data, nil
}
