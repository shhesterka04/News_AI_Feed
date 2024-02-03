package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
	"news_ai_feed/internal/model"
	"time"
)

type SourcePostgresStorage struct {
	db *sqlx.DB
}

func (s *SourcePostgresStorage) Sources(ctx context.Context) ([]model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var sources []dbSource
	query := `SELECT * FROM sources;`
	if err := conn.SelectContext(ctx, &sources, query); err != nil {
		return nil, err
	}

	var modelSources []model.Source
	for _, source := range sources {
		modelSources = append(modelSources, model.Source{
			ID:        source.ID,
			Name:      source.Name,
			FeedURL:   source.FeedURL,
			CreatedAt: source.CreatedAt,
		})
	}
	return modelSources, nil
}

func (s *SourcePostgresStorage) SourceByID(ctx context.Context, id int64) (*model.Source, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var source dbSource
	query := `SELECT * FROM sources WHERE id = $1;`
	if err := conn.GetContext(ctx, &source, query, id); err != nil {
		return nil, err
	}

	return &model.Source{
		ID:        source.ID,
		Name:      source.Name,
		FeedURL:   source.FeedURL,
		CreatedAt: source.CreatedAt,
	}, nil
}

func (s *SourcePostgresStorage) Add(ctx context.Context, source model.Source) (int64, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	var id int64
	query := `INNSERT INTO sources (name, feed_url, created_at) VALUES ($1, $2, $3) RETURNING id;`
	row := conn.QueryRowxContext(ctx, query, source.Name, source.FeedURL, source.CreatedAt)
	if err := row.Err(); err != nil {
		return -1, err
	}

	if err := row.Scan(&id); err != nil {
		return -1, err
	}

	return id, nil

}

func (s *SourcePostgresStorage) Delete(ctx context.Context, id int64) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	query := `DELETE FROM sources WHERE id = $1;`

	if _, err := conn.ExecContext(ctx, query, id); err != nil {
		return err
	}

	return nil
}

type dbSource struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	FeedURL   string    `db:"feed_url"`
	CreatedAt time.Time `db:"created_at"`
}
