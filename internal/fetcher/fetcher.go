package fetcher

import (
	"context"
	"log"
	"news_ai_feed/internal/model"
	"news_ai_feed/internal/source"
	"strings"
	"sync"
	"time"
)

type ArticleStorage interface {
	Store(ctx context.Context, article model.Article) error
}

type SourceProvider interface {
	Sources(ctx context.Context) ([]model.Source, error)
}

type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

type Fetcher struct {
	articles ArticleStorage
	sources  SourceProvider

	fetchInterval  time.Duration
	filterKeywords []string
}

func NewFetcher(articles ArticleStorage, sources SourceProvider, fetchInterval time.Duration, filterKeywords []string) *Fetcher {
	return &Fetcher{
		articles:       articles,
		sources:        sources,
		fetchInterval:  fetchInterval,
		filterKeywords: filterKeywords,
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	for _, src := range sources {
		wg.Add(1)

		rssSource := source.NewRSSSource(src)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				log.Printf("[ERROR] Fetching items from source %s: %v", source, err)
				return
			}

			if err := f.processItems(ctx, source, items); err != nil {
				log.Printf("[ERROR] Processing items from source %s: %v", source, err)
				return
			}
		}(rssSource)

		wg.Wait()
	}

	return nil
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInterval)
	defer ticker.Stop()

	if err := f.Fetch(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}

func (f *Fetcher) processItems(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.itemFilter(item) {
			continue
		}

		if err := f.articles.Store(ctx, model.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			FeedURL:     item.Link,
			Summary:     item.Summary,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fetcher) itemFilter(item model.Item) bool {
	categories := make(map[string]interface{})

	for _, keyword := range f.filterKeywords {
		keywordInTitle := strings.Contains(strings.ToLower(item.Title), keyword)

		if _, ok := categories[keyword]; ok {
			return true
		}

		if keywordInTitle {
			return true
		}
	}

	return false
}
