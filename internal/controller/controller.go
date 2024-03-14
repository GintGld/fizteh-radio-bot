package controller

import (
	"context"
	_ "embed"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Analogue for web cookies.
// Used to store sessions.
type Session interface {
	// Extract current status.
	Status(id int64) string
	// Redirect to another route path.
	Redirect(id int64, cmd string)
}

type Command string

// Status used to define
// when session is closed.
const NullStatus string = "/null"

type OnCancelHandler func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage)

// TODO make one function for answer callback for all controllers.
// Or make an inheritance after struct with this method.

// TODO: create struct, parse replicas to it, use in all controllers.
// And move this code to nother file in the same package.
// // // go:embed replicas.yaml
// var messages string

const (
	// "/start" command
	HelloMessage         = "Привет, Для начала тебе надо авторизироваться, введи логин от админа."
	GotLoginAskPass      = "А тепепь пароль."
	AuthorizedMessage    = "Кастуй или пиздуй, %s."
	WelcomeMessage       = "Добро пожаловать, %s\\!\nТеперь можешь тыкать куда угодно."
	ErrAuthorizedMessage = "Логин или пароль неверны."
	ErrEmptyLogin        = "Логин не может быть пустым."
	ErrEmptyPass         = "Пароль не может быть пустым."

	// "/lib" command
	LibMainMenuMessage = "Можем что-нибудь поискать или загрузить новенькое."

	// "/lib/search"
	LibSearchInit               = "Настрой поиск, а потом нажми 'искать'."
	LibSearchAskNameAuthor      = "Отлично, введи название/автора."
	LibSearchAskFormat          = "Отлично, выбери формат медиа."
	LibSearchAskPlaylist        = "Отлично, введи плейлисты через зяпятую."
	LibSearchAskPodcast         = "Отлично, введи подкасты через зяпятую."
	LibSearchAskGenre           = "Отлично, введи жанры через запятую."
	LibSearchAskLang            = "Отлично, введи языки через запятую."
	LibSearchAskMood            = "Отлично, введи нвстроения через запятую."
	LibSearchErrNameAuthorEmpty = "А почему название пустое?"
	LibSearchErrNilOption       = "Ты так получишь фиг знает что, настрой поиск получше."
	LibSearchErrEmptyRes        = "По твоему запросу ничего не нашлось."

	// "/lib/search/pick"
	LibSearchPickSelecting = "Выбор даты и времени."

	// "/lib/upload"
	LibUpload                = "Выбери вариант загрузки."
	LibUploadAskFile         = "Отправь мне файл для скачивания."
	LibUploadFileNotFound    = "Ты не отправил(а) мне файл."
	LibUploadInvalidMimeType = "Я пока могу кушать только .mp3 файлы :("
	LibUploadAskName         = "Название для файла."
	LibUploadAskAuthor       = "Имя автора."
	LibUploadAskGenre        = "Введи через запятую жанры."
	LibUploadAskFormat       = "Песня или подкаст?"
	LibUploadAskPlaylist     = "Введи через запятую плейлисты, куда добавить песню."
	LibUploadAskPodcast      = "Введи через запятую сезоны, куда добавить подкаст."
	LibUploadAskLang         = "Введи через запятую языки."
	LibUploadAskMood         = "Введи через запятую настроения."
	LibUploadAskLink         = "Отправь мне ссылку на скачивание. Поддерживаемые сервисы на данный момент: Яндекс."
	LibUploadSuccess         = "Загружено."
	LibUploadErrEmptyMsg     = "Не надо делать пустое поле..."

	// "/sch" command
	SchMainMenuMessage = "Можем посмотреть расписание или настроить авто диджея."

	// "/sch/autodj"
	SchAutoDJAskGenre    = "Введи через запятую жанры."
	SchAutoDJAskPlaylist = "Введи через запятую плейлисты, куда добавить песню."
	SchAutoDJAskLanguage = "Введи через запятую языки."
	SchAutoDJAskMood     = "Введи через запятую настроения."
	SchAutoDJSuccess     = "Успешно обновлено."

	// TODO: write help message

	// "/help" command
	HelpMessage = "Здесь будет большое описание функционала (потом)."

	// Unknown user
	ErrUnknown = "Ты кто, сталкер?"

	// undefined behavior
	UndefMsg = "Все плохо, пиши админу."

	// Default error message
	ErrorMessage = "Произошла какая-то ошибка.\nТыкай админа."
)
