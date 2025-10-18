package circuitBreaker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sony/gobreaker"
)

type Client struct {
	client         *resty.Client
	circuitBreaker *gobreaker.CircuitBreaker
	name           string
}

type Config struct {
	Name          string
	MaxRequests   uint32
	Interval      time.Duration
	Timeout       time.Duration
	ReadyToTrip   func(counts gobreaker.Counts) bool
	OnStateChange func(name string, from gobreaker.State, to gobreaker.State)
}

// DefaultCircuitBreakerConfig Default configuration
func DefaultCircuitBreakerConfig(name string) *Config {
	return &Config{
		Name:        name,
		MaxRequests: 5,
		Interval:    60 * time.Second,
		Timeout:     60 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures > 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			fmt.Printf("CircuitBreaker '%s' changed from %s to %s\n", name, from, to)
		},
	}
}

func NewCircuitBreakerClient(config *Config) *Client {
	if config == nil {
		config = DefaultCircuitBreakerConfig("default")
	}

	// Create circuit breaker
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:          config.Name,
		MaxRequests:   config.MaxRequests,
		Interval:      config.Interval,
		Timeout:       config.Timeout,
		ReadyToTrip:   config.ReadyToTrip,
		OnStateChange: config.OnStateChange,
	})

	// Create Resty client
	client := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second)

	return &Client{
		client:         client,
		circuitBreaker: cb,
		name:           config.Name,
	}
}

// ExecuteWithCB executes a request with circuit breaker protection
func (c *Client) ExecuteWithCB(ctx context.Context, req *resty.Request) (*resty.Response, error) {
	result, err := c.circuitBreaker.Execute(func() (interface{}, error) {
		resp, err := req.Execute(req.Method, req.URL)
		if err != nil {
			return nil, err
		}

		// Consider 5xx errors as failures for circuit breaker
		if resp.StatusCode() >= 500 {
			return resp, fmt.Errorf("server error: %d", resp.StatusCode())
		}

		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	return result.(*resty.Response), nil
}

// Convenience methods - FIXED VERSION
func (c *Client) Get(ctx context.Context, url string) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx)
	return c.ExecuteWithCB(ctx, req.SetDoNotParseResponse(false).SetResult(nil))
}

func (c *Client) GetWithResult(ctx context.Context, url string, result interface{}) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx).SetResult(result)
	return c.ExecuteWithCB(ctx, req)
}

func (c *Client) Post(ctx context.Context, url string, body interface{}) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx).SetBody(body)
	return c.ExecuteWithCB(ctx, req)
}

func (c *Client) PostWithResult(ctx context.Context, url string, body interface{}, result interface{}) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx).SetBody(body).SetResult(result)
	return c.ExecuteWithCB(ctx, req)
}

func (c *Client) Put(ctx context.Context, url string, body interface{}) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx).SetBody(body)
	return c.ExecuteWithCB(ctx, req)
}

func (c *Client) PutWithResult(ctx context.Context, url string, body interface{}, result interface{}) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx).SetBody(body).SetResult(result)
	return c.ExecuteWithCB(ctx, req)
}

func (c *Client) Delete(ctx context.Context, url string) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx)
	return c.ExecuteWithCB(ctx, req)
}

func (c *Client) DeleteWithResult(ctx context.Context, url string, result interface{}) (*resty.Response, error) {
	req := c.client.R().SetContext(ctx).SetResult(result)
	return c.ExecuteWithCB(ctx, req)
}

// GetCircuitBreakerState returns the current state of the circuit breaker
func (c *Client) GetCircuitBreakerState() gobreaker.State {
	return c.circuitBreaker.State()
}

// GetCircuitBreakerCounts returns the current counts of the circuit breaker
func (c *Client) GetCircuitBreakerCounts() gobreaker.Counts {
	return c.circuitBreaker.Counts()
}

// SetBaseURL sets the base URL for all requests
func (c *Client) SetBaseURL(url string) *Client {
	c.client.SetBaseURL(url)
	return c
}

// SetHeaders sets common headers for all requests
func (c *Client) SetHeaders(headers map[string]string) *Client {
	c.client.SetHeaders(headers)
	return c
}
