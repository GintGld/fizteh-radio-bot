package random

import (
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/GintGld/fizteh-radio-bot/internal/models"
)

var (
	types = []string{"format", "genre", "playlist", "mood", "language"}
)

func Media() models.Media {
	return models.Media{
		Name:   gofakeit.MovieName(),
		Author: gofakeit.Name(),
		Tags:   TagList(),
	}
}

func TagList() models.TagList {
	const maxLen = 10

	list := make(models.TagList, rand.Intn(maxLen))

	for i := 0; i < len(list); i++ {
		list[i] = Tag()
	}

	return list
}

func Tag() models.Tag {
	return models.Tag{
		Name: gofakeit.Adjective(),
		Type: models.TagTypesAvail[types[rand.Intn(len(models.TagTypesAvail))]],
	}

}

func Segment() models.Segment {
	begin := time.Duration(gofakeit.Uint32())
	stop := begin + time.Duration(gofakeit.Uint32())

	return models.Segment{
		Start:    gofakeit.Date(),
		Media:    Media(),
		BeginCut: begin.Truncate(time.Microsecond),
		StopCut:  stop.Truncate(time.Microsecond),
	}
}
