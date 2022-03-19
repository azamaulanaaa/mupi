package torrentio

import "go.uber.org/zap"

func DefaultConfig() Config {
	return Config{
		Logger: zap.NewExample(),
		Host:   "https://torrentio.strem.fun/lite",
	}
}
