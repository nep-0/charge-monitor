package cache

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type LocalCache struct {
	data map[string]OutletInfo
	mu   sync.RWMutex
}

func NewLocalCache() *LocalCache {
	return &LocalCache{
		data: make(map[string]OutletInfo),
		mu:   sync.RWMutex{},
	}
}

func (c *LocalCache) Get(outletId string) (OutletInfo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	info, exists := c.data[outletId]
	return info, exists
}

func (c *LocalCache) Set(outletId string, info OutletInfo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	info.UpdatedAt = time.Now().Unix()
	c.data[outletId] = info
}

func (c *LocalCache) JSON() []byte {
	// Simple JSON serialization without external libraries, using strings.Builder
	c.mu.RLock()
	defer c.mu.RUnlock()
	if len(c.data) == 0 {
		return []byte("{}")
	}

	var builder strings.Builder
	builder.WriteString("{")
	first := true
	for id, info := range c.data {
		if !first {
			builder.WriteString(",")
		}
		builder.WriteString("\"")
		builder.WriteString(id)
		builder.WriteString("\": {")
		builder.WriteString("\"power\": \"")
		builder.WriteString(info.Power)
		builder.WriteString("\",\"updated_at\": ")
		builder.WriteString(fmt.Sprintf("%d", info.UpdatedAt))
		builder.WriteString(",\"used_minutes\": ")
		builder.WriteString(fmt.Sprintf("%d", info.UsedMinutes))
		builder.WriteString("}")
		first = false
	}
	builder.WriteString("}")
	return []byte(builder.String())
}

func (c *LocalCache) LoadFromJSON(data []byte) error {
	// Simple JSON deserialization without external libraries
	// This is a naive implementation and assumes well-formed input
	c.mu.Lock()
	defer c.mu.Unlock()
	var jsonData map[string]map[string]any
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return err
	}
	for id, info := range jsonData {
		outletInfo := OutletInfo{
			Power:       info["power"].(string),
			UsedMinutes: int64(info["used_minutes"].(float64)),
			UpdatedAt:   int64(info["updated_at"].(float64)),
		}
		// Directly assign to preserve the original UpdatedAt timestamp
		c.data[id] = outletInfo
	}
	return nil
}
