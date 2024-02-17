-- +goose Up
-- +goose StatementBegin
INSERT INTO sources (name, feed_url, created_at, updated_at) VALUES
('The Go Blog', 'https://go.dev/blog/', now(), now()),
('Habr', 'https://habr.com/ru/hubs/go/articles/', now(), now()),
('Reddit', 'https://www.reddit.com/r/golang/', now(), now());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM sources WHERE feed_url IN
('https://go.dev/blog/', 'https://habr.com/ru/hubs/go/articles/', 'https://www.reddit.com/r/golang/');
-- +goose StatementEnd
