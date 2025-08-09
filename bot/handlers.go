package bot

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/ksyusha123/procrastinator-library/storage/articles"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleUpdate(ctx context.Context, update *tgbotapi.Update) {
	b.sendReply(update.Message.Chat.ID, "🐈 Meeeeow!")

	//err := b.userStorage.Save(ctx, update.Message.Chat.ID)
	//if err != nil {
	//	return
	//}
	//if update.Message == nil {
	//	return
	//}
	//
	//if update.Message.IsCommand() {
	//	b.handleCommand(ctx, update.Message)
	//	return
	//}
	//
	//b.handleMessage(ctx, update.Message)
}

func (b *Bot) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	cmdParts := strings.SplitN(msg.Command(), "_", 2)
	baseCmd := cmdParts[0]
	var arg string
	if len(cmdParts) > 1 {
		arg = cmdParts[1]
	}

	switch baseCmd {
	case "start":
		b.handleStart(msg)
	case "save":
		b.handleSave(ctx, msg)
	case "list":
		b.handleList(ctx, msg)
	case "read":
		if arg == "" {
			b.sendReply(msg.Chat.ID, "Usage: /read_<article_id>")
			return
		}
		b.handleMarkRead(ctx, msg, arg)
	case "delete":
		if arg == "" {
			b.sendReply(msg.Chat.ID, "Usage: /delete_<article_id>")
			return
		}
		b.handleDelete(ctx, msg, arg)
	case "help":
		b.handleHelp(msg)
	default:
		b.sendReply(msg.Chat.ID, "Unknown command. Type /help for available commands.")
	}
}

func (b *Bot) handleMarkRead(ctx context.Context, msg *tgbotapi.Message, articleID string) {
	id, err := uuid.Parse(articleID)
	if err != nil {
		b.sendReply(msg.Chat.ID, "Invalid article ID. Must be a UUID.")
		return
	}

	err = b.articleStorage.MarkAsRead(ctx, id, msg.Chat.ID)
	if err != nil {
		b.sendReply(msg.Chat.ID, "Failed to mark article as read.")
		log.Printf("Error marking as read: %v", err)
		return
	}

	b.sendReply(msg.Chat.ID, fmt.Sprintf("✅ Article #%d marked as read", id))
}

func (b *Bot) handleDelete(ctx context.Context, msg *tgbotapi.Message, articleID string) {
	id, err := uuid.Parse(articleID)
	if err != nil {
		b.sendReply(msg.Chat.ID, "Invalid article ID. Must be a number.")
		return
	}

	err = b.articleStorage.Delete(ctx, id, msg.Chat.ID)
	if err != nil {
		b.sendReply(msg.Chat.ID, "Failed to delete article.")
		log.Printf("Error deleting article: %v", err)
		return
	}

	b.sendReply(msg.Chat.ID, fmt.Sprintf("🗑 Article #%d deleted", id))
}

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	urls := extractLinks(msg.Text)
	if urls == nil {
		b.sendReply(msg.Chat.ID, "Please use commands to interact with me. Type /help for available commands.")
		return
	}
	for _, u := range urls {
		b.innerHandleSave(ctx, u, getTitle(u), msg.Chat.ID)
	}
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

func extractLinks(text string) []string {
	re := regexp.MustCompile(`(https?://[^\s]+)`)
	links := re.FindAllString(text, -1)
	return links
}

func (b *Bot) handleSave(ctx context.Context, msg *tgbotapi.Message) {
	sharedArticle := strings.TrimSpace(msg.CommandArguments())
	if sharedArticle == "" {
		if msg.ReplyToMessage != nil {
			sharedArticle = msg.ReplyToMessage.Text
		} else {
			sharedArticle = msg.Text
		}
	}

	urls := extractLinks(sharedArticle)

	if urls == nil {
		b.sendReply(msg.Chat.ID, "Please provide a valid URL starting with http:// or https://")
		return
	}

	for _, u := range urls {
		b.innerHandleSave(ctx, u, getTitle(u), msg.Chat.ID)
	}
}

func getTitle(u string) string {
	resp, err := http.Get(u)
	if err != nil {
		log.Fatalf("failed to download %s: %v\n", u, err)
	}
	defer resp.Body.Close()

	parsedURL, err := url.Parse(u)
	if err != nil {
		log.Fatalf("error parsing url")
	}

	article, err := readability.FromReader(resp.Body, parsedURL)
	if err != nil {
		log.Fatalf("failed to parse %s: %v\n", u, err)
	}

	return article.Title
}

func (b *Bot) innerHandleSave(ctx context.Context, url string, title string, chatID int64) {
	article := &articles.Article{
		ID:     uuid.New(),
		URL:    url,
		Title:  title,
		UserID: chatID,
	}

	if err := b.articleStorage.Save(ctx, article); err != nil {
		log.Printf("Error saving article: %v", err)
		b.sendReply(chatID, "Failed to save article. Please try again.")
		return
	}

	reply := fmt.Sprintf("✅ *Article saved!*\n\n*Title:* %s\n*URL:* %s",
		article.Title, article.URL)
	b.sendReply(chatID, reply)
}

func (b *Bot) handleList(ctx context.Context, msg *tgbotapi.Message) {
	articles, err := b.articleStorage.Get(ctx, msg.Chat.ID)
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
		readCommand := fmt.Sprintf("/read\\_%d", article.ID)
		deleteCommand := fmt.Sprintf("/delete\\_%d", article.ID)

		sb.WriteString(fmt.Sprintf("%d. %s [%s]\n%s\n%s | %s\n\n",
			i+1, article.Title, status, article.URL, readCommand, deleteCommand))

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

func (b *Bot) sendReply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown
	if _, err := b.botAPI.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}
