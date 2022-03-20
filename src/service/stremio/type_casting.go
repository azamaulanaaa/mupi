package stremio

import (
	"context"
	"errors"

	"go.uber.org/zap"
)

var (
	ErrSomethingWrong = errors.New("something wrong")
)

type Config struct {
	Logger  *zap.Logger
	BaseUrl string
}

type Stremio interface {
	Stream(ctx context.Context, typee string, id string) (res StreamsResponse, err error)
}

type StreamsResponse struct {
	Streams         []StreamItem `json:"streams,omitempty"`
	CacheMaxAge     uint32       `json:"cacheMaxAge"`
	StaleRevalidate uint32       `json:"staleRevalidate"`
	StaleError      uint32       `json:"staleError"`
}

type StreamItem struct {
	Title       string `json:"title"`
	InfoHash    string `json:"infoHash,omitempty"`
	FileIdx     uint8  `json:"fileIdx,omitempty"`
	Url         string `json:"url,omitempty"`
	YtId        string `json:"ytId,omitempty"`
	ExternalUrl string `json:"externalUrl,omitempty"`
}

type MetaItem struct {
	Name   string   `json:"name"`
	Genres []string `json:"genres,omitempty"`
	Poster string   `json:"-"`
}

type MetaItemJson struct {
	Id          string   `json:"id"`
	Type        string   `json:"type"`
	Name        string   `json:"name"`
	Genres      []string `json:"genres"`
	Poster      string   `json:"poster"`
	PosterShape string   `json:"posterShape,omitempty"`
}

type CatalogItem struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

type Manifest struct {
	Id          string         `json:"id"`
	Version     string         `json:"version"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Types       []string       `json:"types"`
	Catalogs    []CatalogItem  `json:"catalogs"`
	Resources   []ResourceItem `json:"resources"`
}

type ResourceItem struct {
	Name       string   `json:"name"`
	Types      []string `json:"types"`
	IdPrefixes []string `json:"idPrefixes"`
}
