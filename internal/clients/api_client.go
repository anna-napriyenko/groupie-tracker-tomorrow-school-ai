package clients

import (
	"fmt"
	"groupie-tracker-visualizations/internal/models"
	"log"
	"time"

	"github.com/go-resty/resty/v2"
)

type APIClient struct {
	client *resty.Client
}

func NewAPIClient() *APIClient {
	client := resty.New().
		SetBaseURL(BaseURL).
		SetTimeout(10 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			log.Printf("[HTTP] → %s %s", req.Method, req.URL)
			return nil
		}).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			log.Printf("[HTTP] ← %s %s [%d] (%v)",
				resp.Request.Method,
				resp.Request.URL,
				resp.StatusCode(),
				resp.Time())
			return nil
		}).
		OnError(func(req *resty.Request, err error) {
			log.Printf("[HTTP ERROR] %s %s: %v", req.Method, req.URL, err)
		})

	return &APIClient{
		client: client,
	}
}

func (c *APIClient) FetchArtists() ([]models.Artist, error) {
	log.Println("[APIClient] Fetching artists...")

	var artists []models.Artist
	resp, err := c.client.R().
		SetResult(&artists).
		Get(ArtistsPath)

	if err != nil {
		log.Printf("[APIClient] Failed to fetch artists: %v", err)
		return nil, fmt.Errorf("failed to fetch artists: %w", err)
	}

	if resp.IsError() {
		log.Printf("[APIClient] API returned error status %d for artists", resp.StatusCode())
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	log.Printf("[APIClient] Successfully fetched %d artists", len(artists))
	return artists, nil
}

func (c *APIClient) FetchLocations() ([]models.Location, error) {
	log.Println("[APIClient] Fetching locations...")

	var wrapper struct {
		Index []models.Location `json:"index"`
	}
	resp, err := c.client.R().
		SetResult(&wrapper).
		Get(LocationsPath)

	if err != nil {
		log.Printf("[APIClient] Failed to fetch locations: %v", err)
		return nil, fmt.Errorf("failed to fetch locations: %w", err)
	}

	if resp.IsError() {
		log.Printf("[APIClient] API returned error status %d for locations", resp.StatusCode())
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	log.Printf("[APIClient] Successfully fetched %d locations", len(wrapper.Index))
	return wrapper.Index, nil
}

func (c *APIClient) FetchDates() ([]models.Date, error) {
	log.Println("[APIClient] Fetching dates...")

	var wrapper struct {
		Index []models.Date `json:"index"`
	}
	resp, err := c.client.R().
		SetResult(&wrapper).
		Get(DatesPath)

	if err != nil {
		log.Printf("[APIClient] Failed to fetch dates: %v", err)
		return nil, fmt.Errorf("failed to fetch dates: %w", err)
	}

	if resp.IsError() {
		log.Printf("[APIClient] API returned error status %d for dates", resp.StatusCode())
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	log.Printf("[APIClient] Successfully fetched %d dates", len(wrapper.Index))
	return wrapper.Index, nil
}

func (c *APIClient) FetchRelation() ([]models.Relation, error) {
	log.Println("[APIClient] Fetching relations...")

	var wrapper struct {
		Index []models.Relation `json:"index"`
	}
	resp, err := c.client.R().
		SetResult(&wrapper).
		Get(RelationPath)

	if err != nil {
		log.Printf("[APIClient] Failed to fetch relations: %v", err)
		return nil, fmt.Errorf("failed to fetch relations: %w", err)
	}

	if resp.IsError() {
		log.Printf("[APIClient] API returned error status %d for relations", resp.StatusCode())
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode())
	}

	log.Printf("[APIClient] Successfully fetched %d relations", len(wrapper.Index))
	return wrapper.Index, nil
}
