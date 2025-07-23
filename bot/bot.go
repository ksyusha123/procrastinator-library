package bot

import (
	"context"
	"log"

	// "time"
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ksyusha123/procrastinator-library/storage"
)

type Article = storage.Article

type Bot struct {
	botAPI   *tgbotapi.BotAPI
	articleStorage       storage.ArticleStorage
	commands map[string]string
}

func New(botAPI *tgbotapi.BotAPI, db storage.ArticleStorage) *Bot {
	return &Bot{
		botAPI: botAPI,
		articleStorage:     db,
		commands: map[string]string{
			"save":   "Save an article (reply to message or provide URL)",
			"list":   "List your saved articles",
			"read":   "Mark article as read (provide article ID)",
			"delete": "Delete article (provide article ID)",
			"help":   "Show available commands",
		},
	}
}

func (b *Bot) Start(ctx context.Context) {
	log.Println("Starting article bot...")

	_, err := readLastUpdateId()
	if err != nil {
		return
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.botAPI.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			b.handleUpdate(&update)
		case <-ctx.Done():
			log.Println("Stopping bot updates")
			return
		}
	}
}

func (b *Bot) StartNotificationScheduler() {
	c := cron.New(cron.WithLocation(time.UTC))

	// Schedule for every Monday at 9 AM (UTC) "0 9 * * 1"
	_, err := c.AddFunc("* * * * *", func() {
		b.sendWeeklyNotifications()
	})
	if err != nil {
		log.Printf("Failed to schedule notifications: %v", err)
		return
	}

	log.Println("Notification scheduler started")
	c.Start()
}

func (b *Bot) sendWeeklyNotifications() {
	articles, err := b.articleStorage.GetAllUnreadArticles()
	if err != nil {
		log.Printf("Error getting unread articles for notifications: %v", err)
		return
	}

	userArticles := groupByUserID(articles)

	for userID := range userArticles {
		b.sendUserNotification(userID, userArticles[userID]) // congratulate if no unread
	}
}

func groupByUserID(articles []Article) map[int64][]Article {
    groups := make(map[int64][]Article)
    
    for _, article := range articles {
        groups[article.UserID] = append(groups[article.UserID], article)
    }
    
    return groups
}

func (b *Bot) sendUserNotification(chatID int64, articles []Article) {
	var sb strings.Builder
	sb.WriteString("📚 *You have unread articles:*\n\n")

	for i, article := range articles {
		sb.WriteString(fmt.Sprintf(
			"%d. [%s](%s)\nSaved: %s\n\n",
			i+1,
			article.Title,
			article.URL,
			article.CreatedAt.Format("2006-01-02"),
		))
	}

	sb.WriteString(fmt.Sprintf(
		"Use `%s` to mark as read",
		b.generateReadCommand(articles[0].ID),
	))

	msg := tgbotapi.NewMessage(chatID, sb.String())
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableWebPagePreview = true

	if _, err := b.botAPI.Send(msg); err != nil {
		log.Printf("Error sending notification to %d: %v", chatID, err)
	}
}

// Helper function to escape MarkdownV2 special characters
func escapeMarkdown(text string) string {
    replacer := strings.NewReplacer(
        "_", "\\_",
        "*", "\\*",
        "[", "\\[",
        "]", "\\]",
        "(", "\\(",
        ")", "\\)",
        "~", "\\~",
        "`", "\\`",
        ">", "\\>",
        "#", "\\#",
        "+", "\\+",
        "-", "\\-",
        "=", "\\=",
        "|", "\\|",
        "{", "\\{",
        "}", "\\}",
        ".", "\\.",
        "!", "\\!",
    )
    return replacer.Replace(text)
}

// Updated command generation with escaping
func (b *Bot) generateReadCommand(articleID int) string {
    return fmt.Sprintf("/read\\_%d", articleID) // Escaped underscore
}

// func (b *Bot) generateDeleteCommand(articleID int) string {
//     return fmt.Sprintf("/delete\\_%d", articleID) // Escaped underscore
// }
