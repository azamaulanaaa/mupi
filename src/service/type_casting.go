package service

import (
	"context"
	"errors"
	"io"
	"mupi/src/service/imdb"
	"time"

	"github.com/spf13/afero"
	"go.uber.org/zap"
)

var (
	ErrNotFound = errors.New("not found")
)

const (
	TypeMovie uint = iota
	TypeSeries
)

type Config struct {
	Logger               *zap.Logger
	CacheFs              afero.Fs
	CacheLifetime        time.Duration
	CacheCleanUpInterval time.Duration
	Port                 uint16
	Timeout              time.Duration
}

type Service interface {
	Search(ctx context.Context, query string) (movies []Movie, err error)
	Stream(ctx context.Context, typee uint, id string) (streams []Stream, err error)
	Reader(ctx context.Context, infohash string, index int) (r io.ReadSeekCloser, err error)
}

type Movie struct {
	imdb.Movie
}

type Stream struct {
	Name     string
	InfoHash string
	Index    uint8
}
