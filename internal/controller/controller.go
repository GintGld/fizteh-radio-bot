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

type OnSelectHandler func(ctx context.Context, b *bot.Bot, mes models.MaybeInaccessibleMessage)

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
	WelcomeMessage       = "Добро пожаловать, %s!\nТеперь можешь тыкать куда угодно."
	ErrAuthorizedMessage = "Логин или пароль неверны. Попробуем еще раз сначала."
	ErrEmptyLogin        = "Логин не может быть пустым."
	ErrEmptyPass         = "Пароль не может быть пустым."

	// "/lib/search"
	LibSearchInit               = "Настрой поиск, а потом нажми 'искать'."
	LibSearchAskNameAuthor      = "Отлично, введи название/автора."
	LibSearchAskFormat          = "Отлично, выбери формат медиа."
	LibSearchAskPlaylist        = "Отлично, введи плейлисты через зяпятую."
	LibSearchAskPodcast         = "Отлично, введи подкасты через зяпятую."
	LibSearchAskGenre           = "Отлично, введи жанры через запятую."
	LibSearchAskLang            = "Отлично, введи языки через запятую."
	LibSearchAskMood            = "Отлично, введи настроения через запятую."
	LibSearchErrNameAuthorEmpty = "А почему название пустое?"
	LibSearchErrNilOption       = "Ты так получишь фиг знает что, настрой поиск получше."
	LibSearchErrEmptyRes        = "По твоему запросу ничего не нашлось."

	// "/lib/search" update
	LibSearchUpdatedSuccess    = "Успешно обновлено."
	LibSearchUpdateAskName     = "Название."
	LibSearchUpdateAskAuthor   = "Имя автора."
	LibSearchUpdateAskGenre    = "Выбирай жанры."
	LibSearchUpdateAskPlaylist = "Введи через запятую плейлисты, куда добавить песню."
	LibSearchUpdateAskPodcast  = "Введи через запятую сезоны, куда добавить подкаст."
	LibSearchUpdateAskLang     = "Выбирай языки."
	LibSearchUpdateAskMood     = "Выбирай настроения."
	LibSearchUpdateErrEmptyMsg = "Не надо делать пустое поле..."

	// "/lib/search" delete
	LibSearchDeleteSubmit  = "Точно ли хочешь удалить?"
	LibSearchDeleteSuccess = "Успешно удалено."

	// "/lib/search/pick"
	LibSearchPickSelecting = "Выбор даты и времени."

	// "/lib/upload"
	LibUpload                      = "Выбери вариант загрузки."
	LibUploadAskFile               = "Отправь мне файл для скачивания."
	LibUploadFileNotFound          = "Ты не отправил(а) мне файл."
	LibUploadInvalidMimeType       = "Я пока могу кушать только .mp3 файлы :("
	LibUploadAskName               = "Название."
	LibUploadAskAuthor             = "Имя автора."
	LibUploadAskAlbum              = "Введи через запятую альбомы."
	LibUploadAskGenre              = "Выбирай жанры."
	LibUploadAskPlaylist           = "Введи через запятую плейлисты, куда добавить песню."
	LibUploadAskPodcast            = "Введи через запятую сезоны, куда добавить подкаст."
	LibUploadAskLang               = "Выбирай языки."
	LibUploadAskMood               = "Выбирай настроения."
	LibUploadAskLink               = "Отправь мне ссылку на скачивание. Поддерживаемые сервисы на данный момент: Яндекс."
	LibUploadSuccess               = "Загружено."
	LibUploadErrEmptyMsg           = "Не надо делать пустое поле..."
	LibUploadErrInvalidLink        = "Не могу распознать твою ссылку"
	LibUploadErrMediaAlreadyExists = "Композиция с таким названием и автором уже существует. Если хочешь ее отредактировать, используй поиск в библиотеке."

	// "/sch" command
	SchEmptySchedule = "Расписание пока пусто"

	// "/sch/autodj"
	SchAutoDJAskGenre    = "Выбирай жанры."
	SchAutoDJAskPlaylist = "Введи через запятую плейлисты."
	SchAutoDJAskLanguage = "Выбирай языки."
	SchAutoDJAskMood     = "Выбирай настроения."
	SchAutoDJSuccess     = "Успешно обновлено."

	LiveAskName    = "Введи название эфира."
	LiveNotPlaying = "Сейчас эфир не идет."
	LiveStarted    = "Эфир запущен."
	LiveSubmitStop = "Точно ли хочешь остановить эфир?"
	LiveStopped    = "Эфир остановлен."
	LiveNameEmpty  = "Название эфира не может быть пустым."

	// TODO: write help message

	// "/help" command
	HelpMessage = "Здесь будет большое описание функционала (потом)."

	// in progress
	InProgress = "Работаю..."

	// Unknown user
	ErrUnknown = "Ты кто, сталкер?"

	// undefined behavior
	UndefMsg = "Все плохо, пиши админу."

	// unexpected message
	UnexpectedMsg = "Ну и зачем ты мне это прислал(а)?"

	// Default error message
	ErrorMessage = "Произошла какая-то ошибка.\nТыкай админа."
)
