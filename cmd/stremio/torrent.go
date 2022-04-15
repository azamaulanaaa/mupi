package main

import (
	"mupi/src/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func TorrentHandler(svc service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var infohash string
		{
			infohash = c.Param("infohash")
		}

		var index int
		{
			var err error
			indexStr := c.Param("index")
			index, err = strconv.Atoi(indexStr)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}
		}

		{
			r, err := svc.Reader(c.Request.Context(), infohash, index)
			if err != nil {
				c.Status(http.StatusNotFound)
				return
			}

			http.ServeContent(c.Writer, c.Request, "-", time.Unix(0, 0), r)
		}
	}
}
