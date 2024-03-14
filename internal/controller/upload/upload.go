package upload

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	ctr "github.com/GintGld/fizteh-radio-bot/internal/controller"
	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/storage"
	localModels "github.com/GintGld/fizteh-radio-bot/internal/models"
)

const (
	cmdBase   ctr.Command = ""
	cmdManual ctr.Command = "manual"
	cmdLink   ctr.Command = "link"
	cmdBack   ctr.Command = "back"

	// manual upload
	cmdFile ctr.Command = "file"

	// link upload
	cmdGetLink ctr.Command = "get-link"

	// settings
	cmdSettings      ctr.Command = "settings"
	cmdGetData       ctr.Command = "get-data"
	cmdSubmit        ctr.Command = "submit"
	cmdCancel        ctr.Command = "cancel"
	cmdCancelSetting ctr.Command = "cancel-settings"

	// filler
	cmdNoOp ctr.Command = "no-op"
)

type upload struct {
	router      *ctr.Router
	auth        Auth
	mediaUpload MediaUpload
	session     ctr.Session
	onError     bot.ErrorsHandler
	tmpDir      string

	fileStorage            storage.Storage[string]
	mediaConfigStorage     storage.Storage[localModels.MediaConfig]
	settingTargetStorage   storage.Storage[string]
	linkDownloadResStorage storage.Storage[localModels.LinkDownloadResult]
	msgIdStorage           storage.Storage[int]
}

type Auth interface {
	IsKnown(ctx context.Context, id int64) bool
}

type MediaUpload interface {
	NewMedia(ctx context.Context, id int64, media localModels.MediaConfig, source string) error
	LinkDownload(ctx context.Context, id int64, link string) (localModels.LinkDownloadResult, error)
	LinkUpload(ctx context.Context, id int64, res localModels.LinkDownloadResult) error
}

func Register(
	router *ctr.Router,
	auth Auth,
	mediaUpload MediaUpload,
	session ctr.Session,
	onError bot.ErrorsHandler,
	tmpDir string,
) {
	u := &upload{
		router:      router,
		auth:        auth,
		mediaUpload: mediaUpload,
		session:     session,
		onError:     onError,
		tmpDir:      tmpDir,

		fileStorage:            storage.New[string](),
		mediaConfigStorage:     storage.New[localModels.MediaConfig](),
		settingTargetStorage:   storage.New[string](),
		linkDownloadResStorage: storage.New[localModels.LinkDownloadResult](),
		msgIdStorage:           storage.New[int](),
	}

	router.RegisterCommand(u.init)

	// manual upload
	router.RegisterCallback(cmdManual, u.manualUpload)
	router.RegisterHandler(cmdFile, u.manualUploadFile)

	// link upload
	router.RegisterCallback(cmdLink, u.linkUpload)
	router.RegisterHandler(cmdGetLink, u.getLink)

	// settings
	router.RegisterCallbackPrefix(cmdSettings, u.updateSettings)
	router.RegisterHandler(cmdGetData, u.getSettingNewData)
	router.RegisterCallback(cmdCancelSetting, u.cancelSubTask)
	router.RegisterCallbackPrefix(cmdSubmit, u.submit)
	router.RegisterCallback(cmdCancel, u.returnToMainMenu)

	// filler
	router.RegisterCallback(cmdNoOp, u.nullHandler)
}

func (u *upload) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	chatId := update.Message.Chat.ID

	if !u.auth.IsKnown(ctx, chatId) {
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrUnknown,
		}); err != nil {
			u.onError(err)
		}
		return
	}
	msg, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.LibUpload,
		ReplyMarkup: u.mainMenuMarkup(),
	})
	if err != nil {
		u.onError(err)
	}

	u.msgIdStorage.Set(chatId, msg.ID)
}

func (u *upload) submit(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	var err error

	switch u.router.GetState(update.CallbackQuery.Data) {
	case "manual":
		err = u.mediaUpload.NewMedia(ctx, chatId, u.mediaConfigStorage.Get(chatId), u.fileStorage.Get(chatId))
	case "link":
		err = u.mediaUpload.LinkUpload(ctx, chatId, u.linkDownloadResStorage.Get(chatId))
	}

	if err != nil {
		// handle errors
		if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatId,
			Text:   ctr.ErrorMessage,
		}); err != nil {
			u.onError(err)
		}
		return
	}

	u.fileStorage.Del(chatId)
	u.linkDownloadResStorage.Del(chatId)
	u.mediaConfigStorage.Del(chatId)
	u.settingTargetStorage.Del(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:    chatId,
		MessageID: u.msgIdStorage.Get(chatId),
		Text:      ctr.LibUploadSuccess,
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) returnToMainMenu(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	u.fileStorage.Del(chatId)
	u.linkDownloadResStorage.Del(chatId)
	u.mediaConfigStorage.Del(chatId)
	u.settingTargetStorage.Del(chatId)

	if _, err := b.EditMessageText(ctx, &bot.EditMessageTextParams{
		ChatID:      chatId,
		MessageID:   u.msgIdStorage.Get(chatId),
		Text:        ctr.LibUpload,
		ReplyMarkup: u.mainMenuMarkup(),
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) nullHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)
}

func (u *upload) callbackAnswer(ctx context.Context, b *bot.Bot, callbackQuery *models.CallbackQuery) {
	ok, err := b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
		CallbackQueryID: callbackQuery.ID,
	})
	if err != nil {
		u.onError(err)
		return
	}
	if !ok {
		u.onError(fmt.Errorf("callback answer failed"))
	}
}
