package imdb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gojek/heimdall/v7/httpclient"
)

type imdb_impl struct {
	config     Config
	httpClient *httpclient.Client
}

func New(config Config) IMDB {
	imdb := imdb_impl{
		config:     config,
		httpClient: httpclient.NewClient(httpclient.WithHTTPTimeout(config.Timeout)),
	}

	return &imdb
}

type searchRes struct {
	Movies []Movie `json:"d"`
}

func (imdb *imdb_impl) Search(ctx context.Context, query string) (movies []Movie, err error) {
	var url string
	{
		var path string
		path = fmt.Sprintf("/suggestion/%s/%s.json", string(query[0]), query)
		url = imdb.url(path)
	}

	{
		httpRes, err := imdb.httpClient.Get(url, nil)
		if err != nil {
			return movies, err
		}

		if httpRes.StatusCode != http.StatusOK {
			return movies, ErrSomethingWrong
		}

		{
			var res searchRes
			var buff bytes.Buffer
			buff.ReadFrom(httpRes.Body)
			httpRes.Body.Close()
			err := json.Unmarshal(buff.Bytes(), &res)
			if err != nil {
				return movies, err
			}

			movies = res.Movies
		}
	}

	return
}

func (imdb *imdb_impl) url(path string) string {
	return fmt.Sprintf("%s%s", imdb.config.BaseUrl, path)
}
