package notifier

import (
	"context"
	"news_ai_feed/internal/model"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkPosted(ctx context.Context, id int64)
}

type Summarizer interface {
	Summarize(ctx context.Context, text string) (string, error)
}

type Notifier struct {
	articles         ArticleProvider
	summarizer       Summarizer
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	lookupTimeWindow time.Duration
	channelID        int64
}

func New(articles ArticleProvider, summarizer Summarizer, bot *tgbotapi.BotAPI, sendInterval, lookupTimeWindow time.Duration, channelID int64) *Notifier {
	return &Notifier{
		articles:         articles,
		summarizer:       summarizer,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelID:        channelID,
	}
}

func (n *Notifier) SelectAndSend(ctx context.Context) error {

}
