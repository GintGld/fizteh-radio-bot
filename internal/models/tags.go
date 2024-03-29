package models

import "slices"

// In agreement with server migrations.
// For Genre, Moog, Language type id is not equal
// tag they coreespond to.
var (
	// TODO reqrite tag types avail as variables (like others)
	TagTypesAvail = map[string]TagType{
		"format":   {ID: 1, Name: "format"},
		"genre":    {ID: 2, Name: "genre"},
		"playlist": {ID: 3, Name: "playlist"},
		"mood":     {ID: 4, Name: "mood"},
		"language": {ID: 5, Name: "language"},
		"podcast":  {ID: 6, Name: "podcast"},
		"album":    {ID: 7, Name: "album"},
	}

	Pop          = Genre{Id: 1, Name: "Поп"}
	HipHop       = Genre{Id: 2, Name: "Хип-хоп"}
	Rock         = Genre{Id: 3, Name: "Рок"}
	Jazz         = Genre{Id: 4, Name: "Джаз"}
	Electro      = Genre{Id: 5, Name: "Электро"}
	Instrumental = Genre{Id: 6, Name: "Инструментальный"}
	Rap          = Genre{Id: 7, Name: "Рэп"}
	LoFi         = Genre{Id: 8, Name: "Lo-fi"}

	Aggresive  = Mood{Id: 1, Name: "агрессивное"}
	Optimistic = Mood{Id: 2, Name: "оптимистичное"}
	Calm       = Mood{Id: 3, Name: "спокойное"}
	Anxious    = Mood{Id: 4, Name: "тревожное"}
	Rhythmic   = Mood{Id: 5, Name: "ритмичное"}
	Romantic   = Mood{Id: 6, Name: "романтичное"}
	Sad        = Mood{Id: 7, Name: "печальное"}

	Russian   = Language{Id: 1, Name: "русский"}
	English   = Language{Id: 2, Name: "английский"}
	French    = Language{Id: 3, Name: "французский"}
	Italian   = Language{Id: 4, Name: "итальянский"}
	German    = Language{Id: 5, Name: "немецкий"}
	Spanish   = Language{Id: 6, Name: "испанский"}
	Mongolian = Language{Id: 7, Name: "монгольский"}
	Korean    = Language{Id: 8, Name: "корейский"}
	Japanese  = Language{Id: 9, Name: "японский"}
	Chinese   = Language{Id: 10, Name: "китайский"}
	NoWords   = Language{Id: 11, Name: "без слов"}

	GenresAvail = [GenreNumber]Genre{Pop, HipHop, Rock, Jazz, Electro, Instrumental, Rap, LoFi}
	MoodsAvail  = [MoodNumber]Mood{Aggresive, Optimistic, Calm, Anxious, Rhythmic, Romantic, Sad}
	LangsAvail  = [LangNumber]Language{Russian, English, French, Italian, German, Spanish, Mongolian, Korean, Japanese, Chinese, NoWords}
)

const (
	GenreNumber = 8
	MoodNumber  = 7
	LangNumber  = 11
)

func (a Album) String() string {
	return a.Name
}

func (a Album) Tag() Tag {
	return Tag{
		Name: a.Name,
		Type: TagTypesAvail["album"],
		Meta: map[string]string{
			"author": a.Author,
		},
	}
}

func (g Genre) String() string {
	return g.Name
}

func (g Genre) Tag() Tag {
	return Tag{
		Type: TagTypesAvail["genre"],
		Name: g.Name,
	}
}

func (m Mood) String() string {
	return m.Name
}

func (m Mood) Tag() Tag {
	return Tag{
		Type: TagTypesAvail["mood"],
		Name: m.Name,
	}
}

func (l Language) String() string {
	return l.Name
}

func (l Language) Tag() Tag {
	return Tag{
		Type: TagTypesAvail["language"],
		Name: l.Name,
	}
}

func (t Tag) AsAlbum() Album {
	return Album{
		Name:   t.Name,
		Author: t.Meta["author"],
	}
}

func (t Tag) AsGenre() Genre {
	return Genre{
		Id: GenresAvail[slices.IndexFunc(GenresAvail[:], func(g Genre) bool {
			return g.Name == t.Name
		})].Id,
		Name: t.Name,
	}
}

func (t Tag) AsLang() Language {
	return Language{
		Id: LangsAvail[slices.IndexFunc(LangsAvail[:], func(l Language) bool {
			return l.Name == t.Name
		})].Id,
		Name: t.Name,
	}
}

func (t Tag) AsMood() Mood {
	return Mood{
		Id: MoodsAvail[slices.IndexFunc(MoodsAvail[:], func(m Mood) bool {
			return m.Name == t.Name
		})].Id,
		Name: t.Name,
	}
}
