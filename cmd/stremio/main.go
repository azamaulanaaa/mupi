package main

import (
	"fmt"
	"mupi/src/service"
	"os"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
	"go.uber.org/zap"
)

func main() {
	var port uint16
	{
		port = 8080
		portEnvStr := os.Getenv("PORT")
		if portEnv, err := strconv.Atoi(portEnvStr); err != nil && portEnv != 0 {
			port = uint16(portEnv)
		}
	}

	var svc service.Service
	{
		var dir string
		{
			var err error
			dir, err = os.MkdirTemp("", "mupi-*")
			must(err)
			defer os.RemoveAll(dir)
		}

		var err error
		svcConfig := service.Config{
			Logger:               zap.NewExample(),
			CacheFs:              afero.NewBasePathFs(afero.NewOsFs(), dir),
			CacheLifetime:        30 * time.Minute,
			CacheCleanUpInterval: 5 * time.Second,
			Timeout:              30 * time.Second,
			Port:                 port,
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

	{
		router.GET("manifest.json", ManifestHandler(svc))
		router.GET("stream/:type/:id", StreamHandler(svc))
		router.GET("torrent/:infohash/:index", TorrentHandler(svc))
	}

	{
		host := fmt.Sprintf(":%d", port)
		router.Run(host)
	}
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
