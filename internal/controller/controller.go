package controller

import (
	_ "embed"
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

// TODO: create struct, parse replicas to it, use in all controllers.
// And move this code to nother file in the same package.
// // // go:embed replicas.yaml
// var messages string

const (
	// "/start" command
	HelloMessage         = "Привет! Для начала тебе надо авторизироваться, введи логин от админа."
	GotLoginAskPass      = "А тепепь пароль."
	AuthorizedMessage    = "Кастуй или пиздуй, %s."
	WelcomeMessage       = "Добро пожаловать, %s!\nТеперь можешь тыкать куда угодно."
	ErrAuthorizedMessage = "Логин или пароль неверны."
	ErrEmptyLogin        = "Логин не может быть пустым."
	ErrEmptyPass         = "Пароль не может быть пустым."

	// "/lib" command
	LibMainMenuMessage = "Можем что-нибудь поискать или загрузить новенькое."

	// "/lib/search"
	LibSearchInit               = "Настрой поиск, а потом нажми 'искать'."
	LibSearchAskNameAuthror     = "Отлично, введи название/автора."
	LibSearchAskFormat          = "Отлично, выбери формат медиа."
	LibSearchAskPlaylist        = "Отлично, введи плейлисты через зяпятую."
	LibSearchAskGenre           = "Отлично, введи жанры через запятую."
	LibSearchAskLanguage        = "Отлично, введи языки через запятую."
	LibSearchAskMood            = "Отлично, введи нвстроения через запятую."
	LibSearchErrNameAuthorEmpty = "А почему название пустое?"
	LibSearchErrNilOption       = "Ты так получишь фиг знает что, настрой поиск получше."

	// "/sch" command
	SchMainMenuMessage = "Можем посмотреть расписание или настроить авто диджея."

	// TODO: write help message

	// "/help" command
	HelpMessage = "Здесь будет большое описание функционала (потом)."

	// Unknown user
	ErrUnknown = "Ты кто, сталкер?"

	// Default error message
	ErrorMessage = "Произошла какая-то ошибка.\nТыкай админа."
)
