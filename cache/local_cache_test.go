package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	c := NewLocalCache()

	if c.data == nil {
		t.Fatal("NewCache() returned nil")
	}

	if len(c.data) != 0 {
		t.Errorf("Expected empty cache, got length %d", len(c.data))
	}
}

func TestCache_SetAndGet(t *testing.T) {
	c := NewLocalCache()
	outletId := "outlet-123"
	expectedPower := "25W"
	expectedUsedMinutes := int64(10)

	// Test setting a value
	info := OutletInfo{Power: expectedPower, UsedMinutes: expectedUsedMinutes}
	c.Set(outletId, info)

	// Test getting the value
	retrievedInfo, exists := c.Get(outletId)

	if !exists {
		t.Fatal("Expected outlet to exist in cache")
	}

	if retrievedInfo.Power != expectedPower {
		t.Errorf("Expected power %s, got %s", expectedPower, retrievedInfo.Power)
	}

	if retrievedInfo.UsedMinutes != expectedUsedMinutes {
		t.Errorf("Expected used minutes %d, got %d", expectedUsedMinutes, retrievedInfo.UsedMinutes)
	}

	// Verify that UpdatedAt was set
	if retrievedInfo.UpdatedAt == 0 {
		t.Error("Expected UpdatedAt to be set to current timestamp")
	}

	// Verify UpdatedAt is recent (within last 5 seconds)
	now := time.Now().Unix()
	if retrievedInfo.UpdatedAt < now-5 || retrievedInfo.UpdatedAt > now {
		t.Errorf("UpdatedAt timestamp seems incorrect: %d, expected around %d", retrievedInfo.UpdatedAt, now)
	}
}

func TestCache_GetNonExistent(t *testing.T) {
	c := NewLocalCache()

	info, exists := c.Get("non-existent-outlet")

	if exists {
		t.Error("Expected outlet to not exist in cache")
	}

	// Check that zero value is returned
	expectedInfo := OutletInfo{}
	if info != expectedInfo {
		t.Errorf("Expected zero value OutletInfo, got %+v", info)
	}
}

func TestCache_UpdateExisting(t *testing.T) {
	c := NewLocalCache()
	outletId := "outlet-123"

	// Set initial value
	initialInfo := OutletInfo{Power: "10W", UsedMinutes: 5}
	c.Set(outletId, initialInfo)

	// Get the initial timestamp
	retrievedInfo, _ := c.Get(outletId)
	initialTimestamp := retrievedInfo.UpdatedAt

	// Wait a full second to ensure timestamp difference (Unix timestamps are in seconds)
	time.Sleep(time.Second * 1)

	// Update with new value
	updatedInfo := OutletInfo{Power: "20W", UsedMinutes: 15}
	c.Set(outletId, updatedInfo)

	// Verify the update
	finalInfo, exists := c.Get(outletId)

	if !exists {
		t.Fatal("Expected outlet to exist in cache after update")
	}

	if finalInfo.Power != "20W" {
		t.Errorf("Expected updated power 20W, got %s", finalInfo.Power)
	}

	if finalInfo.UsedMinutes != 15 {
		t.Errorf("Expected updated used minutes 15, got %d", finalInfo.UsedMinutes)
	}

	if finalInfo.UpdatedAt <= initialTimestamp {
		t.Errorf("Expected UpdatedAt to be updated to a newer timestamp. Initial: %d, Final: %d", initialTimestamp, finalInfo.UpdatedAt)
	}
}

func TestCache_JSON_EmptyCache(t *testing.T) {
	c := NewLocalCache()

	jsonBytes := c.JSON()
	jsonStr := string(jsonBytes)

	// Empty cache should return "{}"
	expected := "{}"
	if jsonStr != expected {
		t.Errorf("Expected empty cache JSON to be %s, got %s", expected, jsonStr)
	}
}

func TestCache_JSON_SingleEntry(t *testing.T) {
	c := NewLocalCache()
	outletId := "outlet-123"
	power := "25W"
	usedMinutes := int64(30)

	info := OutletInfo{Power: power, UsedMinutes: usedMinutes}
	c.Set(outletId, info)

	jsonBytes := c.JSON()

	// Parse the JSON to verify structure
	var parsed map[string]map[string]interface{}
	err := json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Generated JSON is not valid: %v", err)
	}

	// Verify the outlet exists in JSON
	outletData, exists := parsed[outletId]
	if !exists {
		t.Fatalf("Expected outlet %s to exist in JSON", outletId)
	}

	// Verify power value
	if outletData["power"] != power {
		t.Errorf("Expected power %s, got %v", power, outletData["power"])
	}

	// Verify used_minutes value
	if int64(outletData["used_minutes"].(float64)) != usedMinutes {
		t.Errorf("Expected used_minutes %d, got %v", usedMinutes, outletData["used_minutes"])
	}

	// Verify updated_at exists and is a number
	updatedAt, exists := outletData["updated_at"]
	if !exists {
		t.Error("Expected updated_at field to exist")
	}

	if _, ok := updatedAt.(float64); !ok {
		t.Errorf("Expected updated_at to be a number, got %T", updatedAt)
	}
}

func TestCache_JSON_MultipleEntries(t *testing.T) {
	c := NewLocalCache()

	// Add multiple outlets
	outlets := map[string]OutletInfo{
		"outlet-1": {Power: "10W", UsedMinutes: 1},
		"outlet-2": {Power: "20W", UsedMinutes: 2},
		"outlet-3": {Power: "30W", UsedMinutes: 3},
	}

	for id, info := range outlets {
		c.Set(id, info)
	}

	jsonBytes := c.JSON()

	// Parse and verify
	var parsed map[string]map[string]interface{}
	err := json.Unmarshal(jsonBytes, &parsed)
	if err != nil {
		t.Fatalf("Generated JSON is not valid: %v", err)
	}

	// Verify all outlets are present
	if len(parsed) != len(outlets) {
		t.Errorf("Expected %d outlets in JSON, got %d", len(outlets), len(parsed))
	}

	for id, expectedInfo := range outlets {
		outletData, exists := parsed[id]
		if !exists {
			t.Errorf("Expected outlet %s to exist in JSON", id)
			continue
		}

		if outletData["power"] != expectedInfo.Power {
			t.Errorf("Expected power %s for outlet %s, got %v", expectedInfo.Power, id, outletData["power"])
		}

		if int64(outletData["used_minutes"].(float64)) != expectedInfo.UsedMinutes {
			t.Errorf("Expected used_minutes %d for outlet %s, got %v", expectedInfo.UsedMinutes, id, outletData["used_minutes"])
		}
	}
}

func TestCache_LoadFromJSON_ValidData(t *testing.T) {
	c := NewLocalCache()

	// Create test JSON data
	testData := `{
		"outlet-1": {
			"power": "15W",
			"used_minutes": 45,
			"updated_at": 1695456789
		},
		"outlet-2": {
			"power": "25W",
			"used_minutes": 60,
			"updated_at": 1695456790
		}
	}`

	err := c.LoadFromJSON([]byte(testData))
	if err != nil {
		t.Fatalf("LoadFromJSON failed: %v", err)
	}

	// Verify outlet-1
	info1, exists := c.Get("outlet-1")
	if !exists {
		t.Error("Expected outlet-1 to exist after loading from JSON")
	} else {
		if info1.Power != "15W" {
			t.Errorf("Expected power 15W for outlet-1, got %s", info1.Power)
		}
		if info1.UsedMinutes != 45 {
			t.Errorf("Expected used minutes 45 for outlet-1, got %d", info1.UsedMinutes)
		}
		if info1.UpdatedAt != 1695456789 {
			t.Errorf("Expected UpdatedAt 1695456789 for outlet-1, got %d", info1.UpdatedAt)
		}
	}

	// Verify outlet-2
	info2, exists := c.Get("outlet-2")
	if !exists {
		t.Error("Expected outlet-2 to exist after loading from JSON")
	} else {
		if info2.Power != "25W" {
			t.Errorf("Expected power 25W for outlet-2, got %s", info2.Power)
		}
		if info2.UsedMinutes != 60 {
			t.Errorf("Expected used minutes 60 for outlet-2, got %d", info2.UsedMinutes)
		}
		if info2.UpdatedAt != 1695456790 {
			t.Errorf("Expected UpdatedAt 1695456790 for outlet-2, got %d", info2.UpdatedAt)
		}
	}
}

func TestCache_LoadFromJSON_InvalidData(t *testing.T) {
	c := NewLocalCache()

	// Test with invalid JSON
	invalidJSON := `{"invalid": json}`

	err := c.LoadFromJSON([]byte(invalidJSON))
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}

func TestCache_LoadFromJSON_EmptyData(t *testing.T) {
	c := NewLocalCache()

	err := c.LoadFromJSON([]byte("{}"))
	if err != nil {
		t.Errorf("Expected no error when loading empty JSON, got %v", err)
	}

	if len(c.data) != 0 {
		t.Errorf("Expected cache to remain empty, got %d entries", len(c.data))
	}
}

func TestCache_RoundTrip_JSONSerialization(t *testing.T) {
	// Test that we can serialize to JSON and deserialize back correctly
	original := NewLocalCache()

	// Add test data
	testData := map[string]OutletInfo{
		"outlet-1": {Power: "10W", UsedMinutes: 10},
		"outlet-2": {Power: "20W", UsedMinutes: 20},
		"outlet-3": {Power: "30W", UsedMinutes: 30},
	}

	for id, info := range testData {
		original.Set(id, info)
	}

	// Serialize to JSON
	jsonBytes := original.JSON()

	// Create new cache and deserialize
	roundTrip := NewLocalCache()
	err := roundTrip.LoadFromJSON(jsonBytes)
	if err != nil {
		t.Fatalf("Round-trip deserialization failed: %v", err)
	}

	// Verify all data is preserved (except UpdatedAt might be different due to Set() call)
	for id, expectedInfo := range testData {
		info, exists := roundTrip.Get(id)
		if !exists {
			t.Errorf("Expected outlet %s to exist after round-trip", id)
			continue
		}

		if info.Power != expectedInfo.Power {
			t.Errorf("Expected power %s for outlet %s after round-trip, got %s", expectedInfo.Power, id, info.Power)
		}

		if info.UsedMinutes != expectedInfo.UsedMinutes {
			t.Errorf("Expected used minutes %d for outlet %s after round-trip, got %d", expectedInfo.UsedMinutes, id, info.UsedMinutes)
		}

		// UpdatedAt should be updated by the Set() call in LoadFromJSON
		if info.UpdatedAt == 0 {
			t.Errorf("Expected UpdatedAt to be set for outlet %s after round-trip", id)
		}
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	// Test concurrent Set and Get operations now that mutex is implemented
	c := NewLocalCache()
	n := 100
	outletId := "concurrent-test"

	var wg sync.WaitGroup

	// Start n goroutines to set the value concurrently
	for i := range n {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(outletId, OutletInfo{Power: fmt.Sprintf("%dW", 100+i), UsedMinutes: int64(i)})
		}(i)
	}

	// Start n goroutines to get the value concurrently
	for range n {
		wg.Go(func() {
			_, _ = c.Get(outletId)
		})
	}

	wg.Wait()

	// After all goroutines, check that the value exists
	retrieved, exists := c.Get(outletId)
	if !exists {
		t.Error("Expected outlet to exist after concurrent access")
	}

	// Power should be in the expected format (last write wins)
	if retrieved.Power == "" {
		t.Error("Expected power to be set after concurrent access")
	}
}
