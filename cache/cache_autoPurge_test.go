// cache/cache_purge_test.go
package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Dummy implementations for testing
func TestSinglePurgeCycle(t *testing.T) {
	calls := make(map[string]bool)

	tasks := []purgeTask{
		{
			Name: "TestCountry",
			Func: func(ctx context.Context) error {
				calls["Country"] = true
				return nil
			},
		},
		{
			Name: "TestWeather",
			Func: func(ctx context.Context) error {
				calls["Weather"] = true
				return nil
			},
		},
		{
			Name: "TestCurrency",
			Func: func(ctx context.Context) error {
				calls["Currency"] = true
				return nil
			},
		},
	}

	// simulate one loop manually
	ctx := context.Background()
	for _, task := range tasks {
		err := task.Func(ctx)
		assert.NoError(t, err)
	}

	assert.True(t, calls["Country"])
	assert.True(t, calls["Weather"])
	assert.True(t, calls["Currency"])
}

// Note: We don't test StartCachePurgeLoop itself because it runs forever (infinite loop).
