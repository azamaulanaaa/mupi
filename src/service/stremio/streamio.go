package stremio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gojek/heimdall/v7/httpclient"
)

type stremio_impl struct {
	config     Config
	httpClient *httpclient.Client
}

func New(config Config) Stremio {
	stremio := stremio_impl{
		config:     config,
		httpClient: httpclient.NewClient(httpclient.WithHTTPTimeout(30 * time.Second)),
	}

	return &stremio
}

func (stremio *stremio_impl) Stream(ctx context.Context, typee string, id string) (res StreamsResponse, err error) {
	var url string
	{
		var path string
		path = fmt.Sprintf("/stream/%s/%s.json", typee, id)
		url = stremio.url(path)
	}

	{
		var req *http.Request
		{
			req, err = http.NewRequestWithContext(
				ctx,

				http.MethodGet,
				url,
				nil,
			)
			if err != nil {
				return res, err
			}
		}

		httpRes, err := stremio.httpClient.Do(req)
		if err != nil {
			return res, err
		}

		if httpRes.StatusCode != http.StatusOK {
			return res, ErrSomethingWrong
		}

		{
			var buff bytes.Buffer
			buff.ReadFrom(httpRes.Body)
			err = json.Unmarshal(buff.Bytes(), &res)
			if err != nil {
				return res, err
			}
		}

		httpRes.Body.Close()
	}

	return
}

func (stremio *stremio_impl) url(path string) string {
	return fmt.Sprintf("%s%s", stremio.config.BaseUrl, path)
}
