package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"github.com/go-telegram/bot"
	botmodel "github.com/go-telegram/bot/models"
)

const (
	botToken = ""
)

func main() {
	b, err := bot.New(botToken)
	if err != nil {
		log.Fatalf("create bot: %v", err)
	}

	// Set up the command handler
	b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, handleHelp)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/base64_enc", bot.MatchTypePrefix, handleBase64Encode)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/base64_dec", bot.MatchTypePrefix, handleBase64Decode)

	// Start listening for updates
	ctx := context.Background()
	log.Println("bot is running...")
	b.Start(ctx)
}

// handleHelp handles the /help command
func handleHelp(ctx context.Context, b *bot.Bot, update *botmodel.Update) {
	helpText := `
I can help you with the following commands:

/help - Show this help message
/base64_enc <text> - Encode the given text to Base64
/base64_dec <base64_string> - Decode the given Base64 string

Just send any of these commands to get started!
`
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   helpText,
	})
	if err != nil {
		log.Printf("send message to chat %d: %v", update.Message.Chat.ID, err)
	}
}

// handleBase64Encode handles the /base64_enc command
func handleBase64Encode(ctx context.Context, b *bot.Bot, update *botmodel.Update) {
	inputText := getUserInput(update.Message.Text)

	if inputText == "" {
		// If there's no text after the command, send a usage message
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please provide a string to encode. Usage: /base64_enc your_text",
		})
		if err != nil {
			log.Printf("send message to chat %d: %v", update.Message.Chat.ID, err)
		}
		return
	}

	// Encode the string to base64
	encodedText := base64.StdEncoding.EncodeToString([]byte(inputText))

	// Send the encoded string back to the user
	responseText := fmt.Sprintf("```\n%s\n```", encodedText)
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      responseText,
		ParseMode: botmodel.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("send message to chat %d: %v", update.Message.Chat.ID, err)
	}
}

// handleBase64Decode handles the /base64_dec command
func handleBase64Decode(ctx context.Context, b *bot.Bot, update *botmodel.Update) {
	inputText := getUserInput(update.Message.Text)

	if inputText == "" {
		// If there's no text after the command, send a usage message
		_, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Please provide a base64 encoded string to decode. Usage: /base64_dec your_base64_string",
		})
		if err != nil {
			log.Printf("send message to chat %d: %v", update.Message.Chat.ID, err)
		}
		return
	}

	// Decode the base64 string
	decodedBytes, err := base64.StdEncoding.DecodeString(inputText)
	if err != nil {
		log.Printf("base64 decode of [%s] failed: %v", inputText, err)
		// If there's an error decoding, send an error message
		_, err = b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: update.Message.Chat.ID,
			Text:   "Failed to decode base64 string. Please make sure it's a valid base64 encoded string.",
		})
		if err != nil {
			log.Printf("send message to chat %d: %v", update.Message.Chat.ID, err)
		}
		return
	}

	// Send the decoded string back to the user
	responseText := fmt.Sprintf("```\n%s\n```", string(decodedBytes))
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      responseText,
		ParseMode: botmodel.ParseModeMarkdown,
	})
	if err != nil {
		log.Printf("send message to chat %d: %v", update.Message.Chat.ID, err)
	}
}

func getUserInput(messageText string) string {
	parts := strings.SplitN(messageText, " ", 2)

	if len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
