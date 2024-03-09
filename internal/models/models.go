package localModels

import "time"

// type Editor struct {
// 	Login string `json:"login"`
// 	Pass  string `json:"pass"`
// }

type Media struct {
	ID       int64         `json:"id"`
	Name     string        `json:"name"`
	Author   string        `json:"author"`
	Duration time.Duration `json:"duration"`
	Tags     TagList       `json:"tags"`
}

type MediaFilter struct {
	Name       string
	Author     string
	Tags       []string
	MaxRespLen int
}

type TagTypes []TagType
type TagList []Tag

type Tag struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Type TagType `json:"type"`
}

type TagType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
