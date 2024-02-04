package storage

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"news_ai_feed/internal/model"
	"time"
)

type ArticlePostgresStorage struct {
	db *sqlx.DB
}

func NewArticleStorage(db *sqlx.DB) *ArticlePostgresStorage {
	return &ArticlePostgresStorage{db: db}
}

func (s *ArticlePostgresStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	query := `INSERT INTO articles (source_id, title, link, summary, published_at) VALUES ($1, $2, $3, $4, $5)
ON CONFLICT DO NOTHING;`
	if _, err := conn.ExecContext(ctx, query, article.SourceID, article.Title, article.Link, article.Summary, article.PublishedAt); err != nil {
		return err
	}

	return nil
}

func (s *ArticlePostgresStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	query := `SELECT * FROM articles WHERE posted_at IS NULL AND published_at >= $1::timestamp ORDER BY published_at DESC LIMIT $2;`
	if err := conn.SelectContext(ctx, &articles, query, since.UTC().Format(time.RFC3339), limit); err != nil {
		return nil, err
	}

	var modelArticles []model.Article
	for _, article := range articles {
		modelArticles = append(modelArticles, model.Article{
			ID:          article.ID,
			SourceID:    article.SourceID,
			Title:       article.Title,
			Link:        article.Link,
			Summary:     article.Summary,
			PublishedAt: article.PublishedAt,
		})
	}

	return modelArticles, nil
}

func (s *ArticlePostgresStorage) MarkPosted(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	query := `UPDATE articles SET posted_at = $1::timestamp WHERE id = $2;`
	if _, err := conn.ExecContext(ctx, query, time.Now().UTC().Format(time.RFC3339), id); err != nil {
		return err
	}

	return nil
}

type dbArticle struct {
	ID          int64        `db:"id"`
	SourceID    int64        `db:"source_id"`
	Title       string       `db:"title"`
	Link        string       `db:"link"`
	Summary     string       `db:"summary"`
	PublishedAt time.Time    `db:"published_at"`
	PostedAt    sql.NullBool `db:"posted_at"`
}
