package imdb

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrSomethingWrong = errors.New("something wrong")
)

type Config struct {
	Logger  *zap.Logger
	BaseUrl string
}

type IMDB interface {
	Search(ctx context.Context, query string) (movies []Movie, err error)
}

type Movie struct {
	ID    string     `json:"id"`
	Title string     `json:"l"`
	Image MovieImage `json:"i"`
}

type MovieImage struct {
	Height   uint   `json:"height"`
	Width    uint   `json:"weight"`
	ImageUrl string `json:"imageUrl"`
}
