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
	cmdFormat        ctr.Command = "format"
	cmdSubmit        ctr.Command = "submit"
	cmdCancel        ctr.Command = "cancel"
	cmdCancelSetting ctr.Command = "cancel-settings"

	// filler
	cmdNoOp ctr.Command = "no-op"
)

type upload struct {
	router      *ctr.Router
	mediaUpload MediaUpload
	session     ctr.Session
	onCancel    ctr.OnCancelHandler
	onError     bot.ErrorsHandler
	tmpDir      string

	fileStorage            storage.Storage[string]
	mediaConfigStorage     storage.Storage[localModels.MediaConfig]
	settingTargetStorage   storage.Storage[string]
	linkDownloadResStorage storage.Storage[localModels.LinkDownloadResult]
}

type MediaUpload interface {
	NewMedia(media localModels.MediaConfig, source string) error
	LinkDownload(link string) (localModels.LinkDownloadResult, error)
	LinkUpload(localModels.LinkDownloadResult) error
}

func Register(
	router *ctr.Router,
	mediaUpload MediaUpload,
	session ctr.Session,
	onCancel ctr.OnCancelHandler,
	onError bot.ErrorsHandler,
	tmpDir string,
) {
	u := &upload{
		router:      router,
		mediaUpload: mediaUpload,
		session:     session,
		onCancel:    onCancel,
		onError:     onError,
		tmpDir:      tmpDir,

		fileStorage:            storage.Storage[string]{},
		mediaConfigStorage:     storage.Storage[localModels.MediaConfig]{},
		settingTargetStorage:   storage.Storage[string]{},
		linkDownloadResStorage: storage.Storage[localModels.LinkDownloadResult]{},
	}

	router.RegisterCallback(cmdBase, u.init)
	router.RegisterCallback(cmdBack, u.back)

	// manual upload
	router.RegisterCallback(cmdManual, u.manualUpload)
	router.RegisterHandler(cmdFile, u.manualUploadFile)

	// link upload
	router.RegisterCallback(cmdLink, u.linkUpload)
	router.RegisterHandler(cmdGetLink, u.getLink)

	// settings
	router.RegisterCallback(cmdSettings, u.updateSettings)
	router.RegisterCallback(cmdFormat, u.updateFormat)
	router.RegisterHandler(cmdGetData, u.getSettingNewData)
	router.RegisterCallback(cmdCancelSetting, u.cancelSubTask)
	router.RegisterCallback(cmdSubmit, u.submit)
	router.RegisterCallback(cmdCancel, u.returnToMainMenu)

	// filler
	router.RegisterCallback(cmdNoOp, u.nullHandler)
}

func (u *upload) init(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	chatId := update.CallbackQuery.Message.Message.Chat.ID

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        ctr.LibUpload,
		ReplyMarkup: u.mainMenuMarkup(),
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) back(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	u.onCancel(ctx, b, update.CallbackQuery.Message)
}

func (u *upload) submit(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	var err error

	switch u.router.GetState(update.CallbackQuery.Data) {
	case "manual":
		err = u.mediaUpload.NewMedia(u.mediaConfigStorage.Get(userId), u.fileStorage.Get(userId))
	case "link":
		err = u.mediaUpload.LinkUpload(u.linkDownloadResStorage.Get(userId))
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

	u.fileStorage.Del(userId)
	u.linkDownloadResStorage.Del(userId)
	u.mediaConfigStorage.Del(userId)
	u.settingTargetStorage.Del(userId)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: chatId,
		Text:   ctr.LibUploadSuccess,
	}); err != nil {
		u.onError(err)
	}
}

func (u *upload) returnToMainMenu(ctx context.Context, b *bot.Bot, update *models.Update) {
	u.callbackAnswer(ctx, b, update.CallbackQuery)

	userId := update.CallbackQuery.From.ID
	chatId := update.CallbackQuery.Message.Message.Chat.ID

	u.fileStorage.Del(userId)
	u.linkDownloadResStorage.Del(userId)
	u.mediaConfigStorage.Del(userId)
	u.settingTargetStorage.Del(userId)

	if _, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
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
