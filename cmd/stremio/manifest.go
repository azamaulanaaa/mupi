package main

import (
	"mupi/src/service"
	"mupi/src/service/stremio"

	"github.com/gin-gonic/gin"
)

func ManifestHandler(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, stremio.Manifest{
			Id:          "com.mupi.addon",
			Version:     "0.0.1",
			Name:        "MUPI",
			Description: "-",
			Types:       []string{"movie", "series"},
			Resources: []stremio.ResourceItem{
				{
					Name:       "stream",
					Types:      []string{"movie", "series"},
					IdPrefixes: []string{"tt"},
				},
			},
			Catalogs: []stremio.CatalogItem{},
		})
	}
}
