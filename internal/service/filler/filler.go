package filler

import (
	"context"
	"math/rand"
	"sync"

	"github.com/GintGld/fizteh-radio-bot/internal/lib/utils/random"
	"github.com/GintGld/fizteh-radio-bot/internal/models"
	"github.com/brianvoe/gofakeit/v6"
)

type Filler struct {
	conf      models.AutoDJInfo
	confMutex sync.Mutex
}

func New() *Filler {
	return &Filler{}
}

func (f *Filler) IsKnown(_ context.Context, _ int64) bool {
	return true
}

func (f *Filler) Login(_ context.Context, _ int64, _, _ string) error {
	return nil
}

func (f *Filler) Search(_ context.Context, _ int64, filter models.MediaFilter) ([]models.MediaConfig, error) {
	respSize := rand.Intn(filter.MaxRespLen)

	res := make([]models.MediaConfig, respSize)

	for i := 0; i < respSize; i++ {
		res[i] = random.Media().ToConfig()
	}

	return res, nil
}

func (f *Filler) NewMedia(_ context.Context, _ int64, _ models.MediaConfig) error {
	return nil
}

func (f *Filler) LinkDownload(_ context.Context, _ int64, _ string) (models.LinkDownloadResult, error) {
	const maxRespLen = 10

	respSize := rand.Intn(maxRespLen)

	res := make([]models.MediaConfig, respSize)

	for i := range res {
		res[i] = random.Media().ToConfig()
	}

	return models.LinkDownloadResult{
		Type: models.ResPlaylist,
		Playlist: models.Playlist{
			Name:   gofakeit.Name(),
			Values: res,
		},
	}, nil
}

func (f *Filler) LinkUpload(_ context.Context, _ int64, _ models.LinkDownloadResult) error {
	return nil
}

func (f *Filler) Schedule(_ context.Context, _ int64) ([]models.Segment, error) {
	const maxScheduleSize = 50

	respSize := rand.Intn(maxScheduleSize)

	res := make([]models.Segment, respSize)

	for i := 0; i < respSize; i++ {
		res[i] = random.Segment()
	}

	return res, nil
}

func (f *Filler) NewSegment(_ context.Context, _ int64, _ models.Segment) error {
	return nil
}

func (f *Filler) Config(_ context.Context, _ int64) (models.AutoDJInfo, error) {
	f.confMutex.Lock()
	defer f.confMutex.Unlock()

	return f.conf, nil
}

func (f *Filler) SetConfig(_ context.Context, _ int64, conf models.AutoDJInfo) error {
	f.confMutex.Lock()
	defer f.confMutex.Unlock()

	f.conf = conf

	return nil
}

func (f *Filler) StartAutoDJ(_ context.Context, _ int64) error {
	return nil
}

func (f *Filler) StopAutoDJ(_ context.Context, _ int64) error {
	return nil
}
