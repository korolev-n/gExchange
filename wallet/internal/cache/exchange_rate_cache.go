package cache

import (
	"sync"
	"time"
	"golang.org/x/sync/singleflight"
)

type ExchangeRateCache struct {
	mu      sync.RWMutex
	rates   map[string]float64
	expires time.Time
	ttl     time.Duration
	sf singleflight.Group
}

func NewExchangeRateCache(ttl time.Duration) *ExchangeRateCache {
	return &ExchangeRateCache{
		rates: make(map[string]float64),
		ttl:   ttl,
	}
}

func (c *ExchangeRateCache) Get() (map[string]float64, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if time.Now().Before(c.expires) && len(c.rates) > 0 {
		// Вернуть копию карты, чтобы не утекла гонка
		result := make(map[string]float64)
		for k, v := range c.rates {
			result[k] = v
		}
		return result, true
	}
	return nil, false
}

func (c *ExchangeRateCache) Set(rates map[string]float64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.rates = make(map[string]float64)
	for k, v := range rates {
		c.rates[k] = v
	}
	c.expires = time.Now().Add(c.ttl)
}

func (c *ExchangeRateCache) GetOrFetch(fetchFunc func() (map[string]float64, error)) (map[string]float64, error) {
	if rates, ok := c.Get(); ok {
		return rates, nil
	}

	// Используем singleflight для защиты от одновременных fetch
	v, err, _ := c.sf.Do("rates", func() (interface{}, error) {
		rates, err := fetchFunc()
		if err != nil {
			return nil, err
		}
		c.Set(rates)
		return rates, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(map[string]float64), nil
}