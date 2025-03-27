package service

import (
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/dto"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/model"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/repository"
)

type Songs interface {
	GetSongs(req dto.GetSongsRequest) ([]model.Song, error)
	GetSongById(id int) (model.Song, error)
	AddSong(song dto.AddSongRequest) (int, error)
	UpdateSong(song dto.UpdateSongRequest) error
	DeleteSong(id int) error
	GetSongLyrics(req dto.GetSongLyricsRequest) ([]string, error)
	SearchSongs(query string) ([]model.Song, error)
}

type ExternalAPI interface {
	GetSongDetails(group, song string) (*dto.SongDetails, error)
}

type Service struct {
	Songs
	externalAPI ExternalAPI
}

func NewService(repo *repository.Repository, externalAPI ExternalAPI) *Service {
	return &Service{
		Songs:       NewSongsService(repo, externalAPI),
		externalAPI: externalAPI,
	}
}
