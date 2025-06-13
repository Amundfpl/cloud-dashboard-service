// cache/cache_purge_logic_test.go
package cache

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/amundfpl/Assignment-2/db"
	"github.com/amundfpl/Assignment-2/testsetup"
	"github.com/amundfpl/Assignment-2/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// Koble til testdatabasen én gang før alle tester
	testsetup.InitTestFirebase()
	os.Exit(m.Run())
}

func TestPurgeOldCountryCache_NoError(t *testing.T) {
	ctx := context.Background()
	err := PurgeOldCountryCache(ctx)
	assert.NoError(t, err)
}

func TestPurgeOldWeatherCache_NoError(t *testing.T) {
	ctx := context.Background()
	err := PurgeOldWeatherCache(ctx)
	assert.NoError(t, err)
}

func TestPurgeOldCurrencyCache_NoError(t *testing.T) {
	ctx := context.Background()
	err := PurgeOldCurrencyCache(ctx)
	assert.NoError(t, err)
}

func TestPurgeEmptyCollection(t *testing.T) {
	ctx := context.Background()
	collection := "test_cache"

	// Legg til dokument med gammel timestamp
	_, _, err := db.FirestoreClient().Collection(collection).Add(ctx, map[string]interface{}{
		utils.TimestampField: time.Now().Add(-2 * time.Hour),
	})
	assert.NoError(t, err)

	// Slett dokumenter eldre enn 1 time
	err = purgeCacheCollection(ctx, collection, 1*time.Hour)
	assert.NoError(t, err)

	// Verifiser at alle dokumenter er slettet
	docs, err := db.FirestoreClient().Collection(collection).Documents(ctx).GetAll()
	assert.NoError(t, err)
	assert.Len(t, docs, 0)

	// Rydd opp i tilfelle noe gjenstår
	for _, doc := range docs {
		_, _ = doc.Ref.Delete(ctx)
	}
}
