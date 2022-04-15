package main

import (
	"mupi/src/service"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

const timeout = 5 * time.Second

func main() {
	var svc service.Service
	{
		var err error
		svcConfig := service.Config{
			Logger:               zap.NewExample(),
			CacheFs:              afero.NewBasePathFs(afero.NewOsFs(), "./cache/"),
			CacheLifetime:        5 * time.Minute,
			CacheCleanUpInterval: 1 * time.Second,
			Timeout:              30 * time.Second,
			Port:                 8080,
		}
		svc, err = service.New(svcConfig)
		must(err)
	}

	var router *gin.Engine
	{
		router = gin.Default()

		corsConfig := cors.DefaultConfig()
		corsConfig.AllowOrigins = []string{"https://app.strem.io"}
		// corsConfig.AllowOrigins = []string{"https://app.strem.io"}
		// corsConfig.AllowAllOrigins = true
		// corsConfig.AllowWebSockets = true
		router.Use(cors.New(corsConfig))
	}

	var host string
	{
		host = "localhost:8080"
	}

	{
		router.GET("manifest.json", ManifestHandler(svc))
		router.GET("stream/:type/:id", StreamHandler(svc, host))
		router.GET("torrent/:infohash/:index", TorrentHandler(svc))
	}

	router.Run(host)
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
