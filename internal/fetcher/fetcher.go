package fetcher

import (
	"context"
	"news_ai_feed/internal/model"
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

//func (f *Fetcher) Fetch(ctx context.Context) error {
//	sources, err := f.sources.Sources(ctx)
//	if err != nil {
//		return err
//	}
//
//	var wg sync.WaitGroup
//
//	for _
//}
