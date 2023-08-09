package initialize

import (
	"context"
	"github.com/allegro/bigcache/v3"
	"time"
)

func TokenCache() *bigcache.BigCache {
	config := bigcache.Config{
		Shards:             1024,
		LifeWindow:         time.Hour * 24,
		CleanWindow:        5 * time.Minute,
		MaxEntriesInWindow: 3000,
		MaxEntrySize:       300,
		StatsEnabled:       false,
		Verbose:            true,
		HardMaxCacheSize:   1024, // memory limit, value in MB
		Logger:             bigcache.DefaultLogger(),
	}
	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		panic(err)
	}
	return cache
}
