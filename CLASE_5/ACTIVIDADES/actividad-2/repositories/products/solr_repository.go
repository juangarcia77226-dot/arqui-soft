package products

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	models "main/models/products"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, httpClient *http.Client) *Client {
	return &Client{baseURL: baseURL, http: httpClient}
}

func (c *Client) Search(ctx context.Context, q models.Query) ([]models.Result, error) {
	params := url.Values{}
	params.Set("q", q.Text)
	params.Set("defType", "edismax")
	params.Set("qf", "title description")
	params.Set("rows", fmt.Sprintf("%d", q.Limit))
	params.Set("fl", "id,title,brand,category,price,score")
	params.Set("wt", "json")
	if q.Category != "" {
		params.Add("fq", "{!term f=category}"+q.Category)
	}
	if q.Brand != "" {
		params.Add("fq", "{!term f=brand}"+q.Brand)
	}

	reqURL := c.baseURL + "/select?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("solr respondió %d", resp.StatusCode)
	}

	var body response
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}

	results := make([]models.Result, 0, len(body.Response.Docs))
	for _, doc := range body.Response.Docs {
		results = append(results, models.Result{
			ID:       doc.ID,
			Title:    doc.Title,
			Brand:    doc.Brand,
			Category: doc.Category,
			Price:    doc.Price,
			Score:    doc.Score,
		})
	}
	return results, nil
}
