package imdb

import (
	"go.uber.org/zap"
)

func DefaultConfig() Config {
	return Config{
		Logger:  zap.NewExample(),
		BaseUrl: "https://v2.sg.media-imdb.com",
	}
}
