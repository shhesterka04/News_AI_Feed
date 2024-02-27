-- +goose Up
-- +goose StatementBegin
INSERT INTO sources (name, feed_url, created_at, updated_at) VALUES
('Habr', 'https://habr.com/ru/rss/hubs/go/articles/?fl=ru', now(), now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sources WHERE feed_url IN
('https://habr.com/ru/rss/hubs/go/articles/?fl=ru');
-- +goose StatementEnd
