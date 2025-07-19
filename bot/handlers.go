package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	// "time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/ksyusha123/procrastinator-library/storage"
)

type Bot struct {
	botAPI   *tgbotapi.BotAPI
	db       storage.Storage
	commands map[string]string
}

func New(botAPI *tgbotapi.BotAPI, db storage.Storage) *Bot {
	return &Bot{
		botAPI: botAPI,
		db:     db,
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

func (b *Bot) handleUpdate(update *tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	if update.Message.IsCommand() {
		b.handleCommand(update.Message)
		return
	}

	b.handleMessage(update.Message)

	writeLastUpdateId(update.UpdateID)
}

func readLastUpdateId() (int, error) {
	content, err := os.ReadFile("offset")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return 0, err
	}
	num, err := strconv.Atoi(string(content))
	if err != nil {
		fmt.Printf("Error converting line: %v\n", err)
		return 0, err
	}
	return num, nil
}

func writeLastUpdateId(lastUpdateId int) error {
	return os.WriteFile("offset", []byte(strconv.Itoa(lastUpdateId+1)), 0644)
}

func (b *Bot) handleCommand(msg *tgbotapi.Message) {
	switch msg.Command() {
	case "start":
		b.handleStart(msg)
	case "save":
		b.handleSave(msg)
	case "list":
		b.handleList(msg)
	case "read":
		b.handleMarkRead(msg)
	case "delete":
		b.handleDelete(msg)
	case "help":
		b.handleHelp(msg)
	default:
		b.sendReply(msg.Chat.ID, "Unknown command. Type /help for available commands.")
	}
}

func (b *Bot) handleMessage(msg *tgbotapi.Message) {
	link, title := extractLinks(msg.Text)
	if link == "" {
		b.sendReply(msg.Chat.ID, "Please use commands to interact with me. Type /help for available commands.")
		return
	}
	b.innerHandleSave(link, title, msg.Chat.ID)
}

func (b *Bot) handleStart(msg *tgbotapi.Message) {
	text := "📚 *Article Bot*\n\n" +
		"I help you save and organize articles you want to read later.\n\n" +
		"*Available commands:*\n"

	for cmd, desc := range b.commands {
		text += fmt.Sprintf("/%s - %s\n", cmd, desc)
	}

	b.sendReply(msg.Chat.ID, text)
}

func (b *Bot) handleHelp(msg *tgbotapi.Message) {
	text := "🛠 *Available commands:*\n"
	for cmd, desc := range b.commands {
		text += fmt.Sprintf("/%s - %s\n", cmd, desc)
	}
	b.sendReply(msg.Chat.ID, text)
}

func extractLinks(text string) (string, string) {
	re := regexp.MustCompile(`(https?://[^\s]+)`)

	link := re.FindString(text)

	remainingText := re.ReplaceAllString(text, "")

	return link, strings.TrimSpace(remainingText)
}

func (b *Bot) handleSave(msg *tgbotapi.Message) {
	sharedArticle := strings.TrimSpace(msg.CommandArguments())
	if sharedArticle == "" {
		if msg.ReplyToMessage != nil {
			sharedArticle = msg.ReplyToMessage.Text
		} else {
			sharedArticle = msg.Text
		}
	}

	url, title := extractLinks(sharedArticle)

	if url == "" {
		b.sendReply(msg.Chat.ID, "Please provide a valid URL starting with http:// or https://")
		return
	}

	b.innerHandleSave(url, title, msg.Chat.ID)
}

func (b *Bot) innerHandleSave(url string, title string, chatId int64) {
	article := &storage.Article{
		URL:     url,
		Title:   title,
		Summary: "Summary would be generated here",
		UserID:  chatId,
	}

	if err := b.db.SaveArticle(article); err != nil {
		log.Printf("Error saving article: %v", err)
		b.sendReply(chatId, "Failed to save article. Please try again.")
		return
	}

	reply := fmt.Sprintf("✅ *Article saved!*\n\n*Title:* %s\n*URL:* %s",
		article.Title, article.URL)
	b.sendReply(chatId, reply)
}

func (b *Bot) handleList(msg *tgbotapi.Message) {
	articles, err := b.db.GetArticles(msg.Chat.ID)
	if err != nil {
		log.Printf("Error getting articles: %v", err)
		b.sendReply(msg.Chat.ID, "Failed to retrieve articles. Please try again.")
		return
	}

	if len(articles) == 0 {
		b.sendReply(msg.Chat.ID, "You have no saved articles yet.")
		return
	}

	var sb strings.Builder
	sb.WriteString("📚 *Your saved articles:*\n\n")

	for i, article := range articles {
		status := "🔴"
		if article.IsRead {
			status = "✅"
		}
		sb.WriteString(fmt.Sprintf("%d. %s [%s]\n%s\n\n",
			i+1, article.Title, status, article.URL))

		// Telegram has message length limits, so we send in chunks
		if i > 0 && i%5 == 0 {
			b.sendReply(msg.Chat.ID, sb.String())
			sb.Reset()
		}
	}

	if sb.Len() > 0 {
		b.sendReply(msg.Chat.ID, sb.String())
	}
}

func (b *Bot) handleMarkRead(msg *tgbotapi.Message) {
	articleID, err := parseArticleID(msg.CommandArguments())
	if err != nil {
		b.sendReply(msg.Chat.ID, "Please provide a valid article ID (number from /list)")
		return
	}

	if err := b.db.MarkAsRead(articleID, msg.Chat.ID); err != nil {
		log.Printf("Error marking article as read: %v", err)
		b.sendReply(msg.Chat.ID, "Failed to mark article as read. Please check the ID and try again.")
		return
	}

	b.sendReply(msg.Chat.ID, fmt.Sprintf("✅ Marked article #%d as read", articleID))
}

func (b *Bot) handleDelete(msg *tgbotapi.Message) {
	articleID, err := parseArticleID(msg.CommandArguments())
	if err != nil {
		b.sendReply(msg.Chat.ID, "Please provide a valid article ID (number from /list)")
		return
	}

	if err := b.db.DeleteArticle(articleID, msg.Chat.ID); err != nil {
		log.Printf("Error deleting article: %v", err)
		b.sendReply(msg.Chat.ID, "Failed to delete article. Please check the ID and try again.")
		return
	}

	b.sendReply(msg.Chat.ID, fmt.Sprintf("🗑 Deleted article #%d", articleID))
}

func (b *Bot) sendReply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := b.botAPI.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func parseArticleID(input string) (int, error) {
	var id int
	_, err := fmt.Sscanf(strings.TrimSpace(input), "%d", &id)
	return id, err
}
