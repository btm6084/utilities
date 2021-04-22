package elasticsearch

import (
	elasticsearch "github.com/elastic/go-elasticsearch/v6"
)

// New configures and returns an Elasticsearch client from parameters.
func New(host, index string) (*elasticsearch.Client, error) {
	// Connect to elastic search
	cfg := elasticsearch.Config{
		Addresses: []string{
			host,
		},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return es, err
	}

	// Do an info request to check the status
	info, err := es.Info()
	if err != nil {
		return es, err
	}

	defer info.Body.Close()

	// Make sure the index exists
	indices, err := es.Indices.Get([]string{index})
	if err != nil {
		return es, err
	}

	defer indices.Body.Close()

	return es, nil
}
