package cache

import (
	"context"
	"github.com/amundfpl/Assignment-2/utils"
	"log"
	"time"
)

// purgeTask defines a single cache purge job.
// It includes the name of the cache collection, the purge function to run, and an error message format.
type purgeTask struct {
	Name string
	Func func(context.Context) error
	Err  string
}

// StartCachePurgeLoop launches a background loop that regularly purges expired cache entries.
// It runs every interval defined by utils.CachePurgeInterval.
func StartCachePurgeLoop() {
	ticker := time.NewTicker(utils.CachePurgeInterval) // Schedule regular purging using a ticker set to every hour
	defer ticker.Stop()                                // Stop the ticker if the loop ever exits (good practice)

	// Define the purge tasks for each cache type, using constants from the utils package
	tasks := []purgeTask{
		{Name: utils.CountryCacheCollection, Func: PurgeOldCountryCache, Err: utils.ErrPurgeCountryCache},
		{Name: utils.WeatherCacheCollection, Func: PurgeOldWeatherCache, Err: utils.ErrPurgeWeatherCache},
		{Name: utils.CurrencyCacheCollection, Func: PurgeOldCurrencyCache, Err: utils.ErrPurgeCurrencyCache},
	}

	// Infinite loop that performs cache purging at the specified interval
	for {
		ctx := context.Background()           // Create a new background context for this purge cycle
		log.Println(utils.MsgCachePurgeStart) // Log the start of the cache purge

		// Execute each purge task in order
		for _, task := range tasks {
			if err := task.Func(ctx); err != nil {
				// Log the task-specific error message if the purge function fails
				log.Printf(task.Err, err)
			}
		}

		log.Println(utils.MsgCachePurgeDone) // Log the end of the purge cycle
		<-ticker.C                           // Wait until the next tick to repeat
	}
}
