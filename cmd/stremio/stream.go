package main

import (
	"fmt"
	"mupi/src/service"
	"mupi/src/service/stremio"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StreamHandler(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var typee uint
		{
			switch c.Param("type") {
			case "movie":
				typee = service.TypeMovie
			case "series":
				typee = service.TypeSeries
			}
		}

		var id string
		{
			id = c.Param("id")
			id = id[:len(id)-5]
		}

		var res stremio.StreamsResponse
		{
			streams, err := svc.Stream(c.Request.Context(), typee, id)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}

			for _, v := range streams {
				res.Streams = append(res.Streams, stremio.StreamItem{
					Title: v.Name,
					Url:   fmt.Sprintf("http://%s/torrent/%s/%d", c.Request.Host, v.InfoHash, v.Index),
				})
			}

			res.CacheMaxAge = 14400
			res.StaleRevalidate = 14400
			res.StaleError = 604800
		}

		c.JSON(http.StatusOK, res)
	}
}
