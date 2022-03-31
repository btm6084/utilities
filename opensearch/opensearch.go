package opensearch

import (
	opensearch "github.com/opensearch-project/opensearch-go"
)

// New configures and returns an Opensearch client from parameters.
func New(host, index string) (*opensearch.Client, error) {
	// Connect to Opensearch
	cfg := opensearch.Config{
		Addresses: []string{
			host,
		},
	}

	search, err := opensearch.NewClient(cfg)
	if err != nil {
		return search, err
	}

	// Do an info request to check the status
	info, err := search.Info()
	if err != nil {
		return search, err
	}

	defer info.Body.Close()

	// Make sure the index exists
	indices, err := search.Indices.Get([]string{index})
	if err != nil {
		return search, err
	}

	defer indices.Body.Close()

	return search, nil
}
