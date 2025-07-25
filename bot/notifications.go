package bot

import (
	"log"
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) StartNotificationScheduler() {
	c := cron.New(cron.WithLocation(time.UTC))

	// Schedule for every Monday at 9 AM (UTC) "0 9 * * 1"
	_, err := c.AddFunc("0 9 * * 1", func() {
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
	users, err := b.userStorage.GetUsersReceivingNotifications()
	if err != nil {
		log.Printf("Error getting users for notifications: %v", err)
		return
	}

	for _, user := range users {
		articles, err := b.articleStorage.GetUnreadArticles(user.ID)
		if err != nil {
			log.Printf("Error getting unread articles for notifications: %v", err)
		}
		b.sendUserNotification(user.ID, articles) // congratulate if no unread
	}
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

		sb.WriteString(fmt.Sprintf(
			"Use `%s` to mark as read",
			b.generateReadCommand(articles[0].ID),
		))
	}

	msg := tgbotapi.NewMessage(chatID, sb.String())
	msg.ParseMode = tgbotapi.ModeMarkdown
	msg.DisableWebPagePreview = true

	if _, err := b.botAPI.Send(msg); err != nil {
		log.Printf("Error sending notification to %d: %v", chatID, err)
	}
}

// Helper function to escape MarkdownV2 special characters
// func escapeMarkdown(text string) string {
// 	replacer := strings.NewReplacer(
// 		"_", "\\_",
// 		"*", "\\*",
// 		"[", "\\[",
// 		"]", "\\]",
// 		"(", "\\(",
// 		")", "\\)",
// 		"~", "\\~",
// 		"`", "\\`",
// 		">", "\\>",
// 		"#", "\\#",
// 		"+", "\\+",
// 		"-", "\\-",
// 		"=", "\\=",
// 		"|", "\\|",
// 		"{", "\\{",
// 		"}", "\\}",
// 		".", "\\.",
// 		"!", "\\!",
// 	)
// 	return replacer.Replace(text)
// }

// Updated command generation with escaping
func (b *Bot) generateReadCommand(articleID int) string {
	return fmt.Sprintf("/read\\_%d", articleID) // Escaped underscore
}

// func (b *Bot) generateDeleteCommand(articleID int) string {
//     return fmt.Sprintf("/delete\\_%d", articleID) // Escaped underscore
// }
