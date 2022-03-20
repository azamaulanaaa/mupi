package service

import (
	"context"
	"io"
	"mupi/src/service/imdb"
	"mupi/src/service/stremio"
	"regexp"

	"github.com/azamaulanaaa/gotor/src/torrentlib"
)

type service_impl struct {
	config        Config
	torrentClient *torrentlib.Client
	imdb          imdb.IMDB
	stremio       stremio.Stremio
	resRe         *regexp.Regexp
}

func New(config Config) (Service, error) {
	var torrentClient *torrentlib.Client
	{
		torrentConfig := torrentlib.ClientConfig{
			Timeout:         config.Timeout,
			FileSystem:      config.CacheFs,
			ClenaUpInterval: config.CacheCleanUpInterval,
			Lifetime:        config.CacheLifetime,
		}

		var err error
		torrentClient, err = torrentlib.NewClient(torrentConfig)
		if err != nil {
			return nil, err
		}
	}

	var IMDb imdb.IMDB
	{
		imdbConfig := imdb.Config{
			Logger:  config.Logger.Named("imdb"),
			BaseUrl: "https://v2.sg.media-imdb.com",
		}

		IMDb = imdb.New(imdbConfig)
	}

	var Torrentio stremio.Stremio
	{
		torrentioConfig := stremio.Config{
			Logger:  config.Logger.Named("torrentio"),
			BaseUrl: "https://torrentio.strem.fun/providers=yts,eztv,rarbg,1337x,thepiratebay,kickasstorrents|qualityfilter=other,scr,cam,unknown|limit=1",
		}

		Torrentio = stremio.New(torrentioConfig)
	}

	var resRe *regexp.Regexp
	{
		var err error
		resRe, err = regexp.Compile(`[\s\.\-\_\(\[](2160p|1080p|720p|480p)[\s\.\-\_\)\]]`)
		if err != nil {
			return nil, err
		}
	}

	service := service_impl{
		config:        config,
		torrentClient: torrentClient,
		imdb:          IMDb,
		stremio:       Torrentio,
		resRe:         resRe,
	}

	return &service, nil
}

func (service *service_impl) Search(ctx context.Context, query string) (movies []Movie, err error) {
	moviesRaw, err := service.imdb.Search(ctx, query)
	if err != nil {
		return
	}

	for _, v := range moviesRaw {
		movies = append(movies, Movie{v})
	}

	return
}

func (service *service_impl) Stream(ctx context.Context, typee uint, id string) (streams []Stream, err error) {
	var typeStr string
	{
		switch typee {
		case TypeMovie:
			typeStr = "movie"
		case TypeSeries:
			typeStr = "series"
		}
	}

	res, err := service.stremio.Stream(ctx, typeStr, id)
	if err != nil {
		return
	}

	for _, v := range res.Streams {
		var name string
		{
			name = v.Title
			match := service.resRe.FindStringSubmatch(name)
			if len(match) == 2 {
				name = match[1]
			}
		}

		streams = append(streams, Stream{
			Name:     name,
			InfoHash: v.InfoHash,
			Index:    v.FileIdx,
		})

	}

	return
}

func (service *service_impl) Reader(ctx context.Context, infohash string, index int) (r io.ReadSeekCloser, err error) {
	var torrent torrentlib.Torrent
	{
		torrent, err = service.torrentClient.AddHash(infohash)
		if err != nil {
			return nil, err
		}
	}

	files := torrent.Files()
	if len(files) < index {
		return nil, ErrNotFound
	}

	return files[index].Reader(), nil
}
