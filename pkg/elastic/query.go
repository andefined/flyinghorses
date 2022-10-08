package elastic

import "strings"

var (
	availableFields = []string{
		"actions",
		"actions_per_day",
		"count",
		"created_at",
		"dates_since",
		"default_profile",
		"description",
		"favorite_count",
		"favorites",
		"ffr",
		"fl",
		"followers",
		"following",
		"full_text",
		"fvfr",
		"hashtags",
		"id",
		"lang",
		"langs",
		"link",
		"listed",
		"location",
		"media",
		"mentions",
		"name",
		"observed_at",
		"party",
		"profile_image_url",
		"protected",
		"retweet_count",
		"source",
		"stf",
		"stfv",
		"target",
		"toxicity",
		"toxicity.*",
		"toxicity.anti_refugee",
		"toxicity.hate",
		"toxicity.insult",
		"toxicity.law_and_order",
		"toxicity.nationalism",
		"toxicity.racism",
		"toxicity.sexism",
		"toxicity.threat",
		"toxicity.toxic",
		"toxicity.severe_toxic",
		"toxicity.identity_hate",
		"tweets",
		"type",
		"urls",
		"user_classification.label.keyword",
		"user_classification.score",
		"user_name",
		"user_partisanship.hyper_partisanship_ratio",
		"user_partisanship.parties.count",
		"user_partisanship.parties.label",
		"user_partisanship.parties.label.keyword",
		"user_partisanship.parties.normalized_ratio",
		"user_partisanship.parties.ratio",
		"user_partisanship.parties.valence",
		"user_partisanship.party_majoriy.label.keyword",
		"user_partisanship.party_majoriy.normalized_ratio",
		"user_partisanship.valence",
		"verified",
		"_source",
	}
)

// Query Struct
type Query struct {
	Query        BoolQuery                `json:"query,omitempty"`
	Aggregations interface{}              `json:"aggs,omitempty"`
	Fields       []string                 `json:"fields,omitempty"`
	Source       []string                 `json:"_source,omitempty"`
	Sort         []map[string]interface{} `json:"sort,omitempty"`
	Size         int                      `json:"size" default:"0"`
	From         int                      `json:"from,omitempty"`
	SearchAfter  []int                    `json:"search_after,omitempty"`
}

// BoolQuery Struct
type BoolQuery struct {
	Bool          *BoolQueryParams `json:"bool,omitempty"`
	FunctionScore interface{}      `json:"function_score,omitempty"`
}

// BoolQueryParams Struct
type BoolQueryParams struct {
	MustNot            []map[string]interface{} `json:"must_not,omitempty"`
	Must               []map[string]interface{} `json:"must,omitempty"`
	Should             []map[string]interface{} `json:"should,omitempty"`
	Filter             []map[string]interface{} `json:"filter,omitempty"`
	MinimunShouldMatch interface{}              `json:"minimum_should_match,omitempty"`
}

// Validate query fields
func (q *Query) Validate() {
	q.Source = make([]string, 0)
	for _, s := range q.Fields {
		if IsValidKey(s, availableFields) {
			q.Source = append(q.Source, s)
		}
	}
}

// IsValidKey checks if a string is in an array
func IsValidKey(str string, list []string) bool {
	for _, s := range list {
		if strings.ToLower(s) == strings.ToLower(str) {
			return true
		}
	}
	return false
}

// RemoveKey removes an invalid key from keys
func RemoveKey(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// NewQuery return a new query
func NewQuery() *Query {
	q := new(Query)
	q.Query.Bool = new(BoolQueryParams)
	return q
}
