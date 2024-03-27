package upload

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	mp3MimeType = "audio/mpeg"
)

func (u *upload) manualUpload(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.manualUpload"

	u.CallbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	u.session.Redirect(chatId, u.router.Path(cmdFile))

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: u.msgIdStorage.Get(chatId),
		Text:      ctr.LibUploadAskFile,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

func (u *upload) manualUploadFile(ctx context.Context, b *bot.Bot, update *models.Update) {
	const op = "upload.manualUploadFile"

	chatId := update.Message.Chat.ID

	if update.Message.Audio == nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.LibUploadFileNotFound,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.LibUploadAskFile,
		})
		if err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		u.msgIdStorage.Set(chatId, msg.ID)
		return
	}

	if update.Message.Audio.MimeType != mp3MimeType {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.LibUploadInvalidMimeType,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.LibUploadAskFile,
		})
		if err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		u.msgIdStorage.Set(chatId, msg.ID)
		return
	}

	file, err := b.GetFile(ctx, &bot.GetFileParams{
		FileID: update.Message.Audio.FileID,
	})
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}

	filepath, err := u.downloadFile(b.FileDownloadLink(file))
	if err != nil {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
		}
		return
	}
	// Automatic removal after 1 hour
	time.AfterFunc(time.Hour, func() {
		u.deleteFile(filepath)
	})

	author, name, found := strings.Cut(update.Message.Audio.FileName, " - ")
	if !found {
		author = ""
		name = ""
	}

	conf := localModels.MediaConfig{
		Name:       name,
		Author:     author,
		Duration:   getMediaDuration(filepath),
		SourcePath: filepath,
	}

	u.mediaConfigStorage.Set(chatId, conf)

	if _, err := b.DeleteMessage(ctx, &bot.DeleteMessageParams{
		ChatID:    chatId,
		MessageID: update.Message.ID,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        conf.String(),
		ReplyMarkup: u.mediaConfMarkup(conf),
		ParseMode:   models.ParseModeHTML,
	}); err != nil {
		u.onError(fmt.Errorf("%s [%d]: %w", op, chatId, err))
	}
}

// download downloads file from telegram.
func (u *upload) downloadFile(link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	out, err := os.CreateTemp(u.tmpDir, "media-upload-*.mp3")
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return out.Name(), err
}

func (u *upload) deleteFile(path string) error {
	return os.Remove(path)
}

// FIXME
func getMediaDuration(_ string) time.Duration {
	return time.Duration(0)
}
