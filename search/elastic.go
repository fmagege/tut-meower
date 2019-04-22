package search

import (
	"context"
	"encoding/json"
	"log"

	// elastic "github.com/elastic/go-elasticsearch/v7"
	// esapi "github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/fmagege/tut-meower/schema"
	elastic "github.com/olivere/elastic"
)

// ElasticRepository .
type ElasticRepository struct {
	client *elastic.Client
}

// NewElastic .
func NewElastic(url string) (*ElasticRepository, error) {
	client, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}
	return &ElasticRepository{client}, nil
}

// NewElastic .
// func NewElastic() (*ElasticRepository, error) {
// 	client, err := elastic.NewDefaultClient()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &ElasticRepository{client}, nil
// }

// Close .
func (r *ElasticRepository) Close() {}

// InsertMeow .
func (r *ElasticRepository) InsertMeow(ctx context.Context, meow schema.Meow) error {
	_, err := r.client.Index().
		Index("meows").
		Type("meow").
		Id(meow.ID).
		BodyJson(meow).
		Refresh("wait_for").
		Do(ctx)
	return err
}

// InsertMeow .
// func (r *ElasticRepository) InsertMeow(ctx context.Context, meow schema.Meow) error {
// 	meowBody, err := json.Marshal(meow)
// 	if err != nil {
// 		fmt.Printf("Error: %s", err)
// 		return err
// 	}

// 	req := esapi.IndexRequest{
// 		Index:        "meows",
// 		DocumentType: "meow",
// 		DocumentID:   meow.ID,
// 		Body:         strings.NewReader(string(meowBody)),
// 		Refresh:      "true",
// 	}

// 	// Perform request with the client.
// 	res, err := req.Do(context.Background(), r.client)
// 	if err != nil {
// 		log.Fatalf("Error getting response: %s", err)
// 		return err
// 	}

// 	defer res.Body.Close()

// 	if res.IsError() {
// 		log.Printf("[%s] Error indexing document ID=%s", res.Status(), meow.ID)
// 	} else {
// 		// Deserialize the response into a map.
// 		var r map[string]interface{}
// 		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
// 			log.Printf("Error parsing the response body: %s", err)
// 		} else {
// 			// Print the response status and indexed document version.
// 			log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
// 		}
// 		return nil
// 	}
// 	return err
// }

// SearchMeows .
func (r *ElasticRepository) SearchMeows(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Meow, error) {
	result, err := r.client.Search().
		Index("meows").
		Query(
			elastic.NewMultiMatchQuery(query, "body").
				Fuzziness("3").
				PrefixLength(1).
				CutoffFrequency(0.0001),
		).
		From(int(skip)).
		Size(int(take)).
		Do(ctx)
	if err != nil {
		return nil, err
	}
	meows := []schema.Meow{}
	for _, hit := range result.Hits.Hits {
		var meow schema.Meow
		if err = json.Unmarshal(*hit.Source, &meow); err != nil {
			log.Println(err)
		}
		meows = append(meows, meow)
	}
	return meows, nil
}

// SearchMeows .
// func (r *ElasticRepository) SearchMeows(ctx context.Context, query string, skip uint64, take uint64) ([]schema.Meow, error) {

// 	res, err := r.client.Search(
// 		r.client.Search.WithContext(context.Background()),
// 		r.client.Search.WithIndex("meows"),
// 		// r.client.Search.WithBody(strings.NewReader(string(schema.Meow))),
// 		r.client.Search.WithTrackTotalHits(true),
// 		r.client.Search.WithPretty(),
// 		r.client.Search.WithSize(int(take)),
// 		r.client.Search.WithFrom(int(skip)),
// 	)

// 	meows := []schema.Meow{}

// 	if err := json.NewDecoder(res.Body).Decode(&meows); err != nil {
// 		log.Fatalf("Error parsing the response body: %s", err)
// 	}

// 	return meows, err
// }
