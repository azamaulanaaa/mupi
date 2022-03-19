package torrentio

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gojek/heimdall/v7/httpclient"
)

type torrentio_impl struct {
	config     Config
	httpClient *httpclient.Client
}

func New(config Config) Torrentio {
	torrentio := torrentio_impl{
		config:     config,
		httpClient: httpclient.NewClient(httpclient.WithHTTPTimeout(config.Timeout)),
	}

	return &torrentio
}

func (torrention *torrentio_impl) Stream(ctx context.Context, typee string, id string) (res StremsResponse, err error) {
	var url string
	{
		var path string
		path = fmt.Sprintf("/stream/%s/%s.json", typee, id)
		url = torrention.url(path)
	}

	{
		httpRes, err := torrention.httpClient.Get(url, nil)
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

func (torrentio *torrentio_impl) url(path string) string {
	return fmt.Sprintf("%s%s", torrentio.config.Host, path)
}
