package elastic

import (
	"crypto/tls"
	"log"
	"net/http"

	"github.com/elastic/go-elasticsearch/v7"
)

type Elastic struct {
	Client *elasticsearch.Client
}

func NewElastic(host, user, pass string) (*Elastic, error) {
	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{host},
			Username:  user,
			Password:  pass,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	)

	if err != nil {
		return nil, err
	}

	res, _ := es.Info()
	log.Println(res)

	return &Elastic{es}, nil
}

func (elastic *Elastic) CreateIndex(name, template string) {
	res, err := elastic.Client.Indices.Create(name)
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}
	if res.IsError() {
		log.Fatalf("Cannot create index: %s", res)
	}
	res.Body.Close()
}

func (elastic *Elastic) CheckIfIndexExists(name string) bool {
	res, err := elastic.Client.Indices.Exists([]string{name})
	if err != nil {
		log.Fatalf("Cannot create index: %s", err)
	}

	res.Body.Close()
	if res.StatusCode != 200 {
		return false
	}
	return true
}
