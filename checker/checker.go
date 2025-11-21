package checker

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Checker struct {
	client *http.Client
}

func NewChecker() *Checker {
	return &Checker{
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Checker) CheckLinks(ctx context.Context, links []string) map[string]string {

	var wg sync.WaitGroup
	results := make(map[string]string)
	var mu sync.Mutex

	for _, link := range links {
		wg.Add(1)
		go func(l string) {
			defer wg.Done()

			req, err := http.NewRequestWithContext(ctx, "GET", "http://"+l, nil)
			if err != nil {
				mu.Lock()
				results[l] = "not available"
				mu.Unlock()
				return
			}

			resp, err := c.client.Do(req)
			if err != nil {
				mu.Lock()
				results[l] = "not available"
				mu.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode < 400 {
				mu.Lock()
				results[l] = "available"
				mu.Unlock()
			} else {
				mu.Lock()
				results[l] = "not available"
				mu.Unlock()
			}
		}(link)
	}

	wg.Wait()
	return results
}
