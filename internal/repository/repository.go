package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/dto"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/model"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/repository/postgres"
)

type Songs interface {
	GetSongs(resp dto.GetSongsRequest) ([]model.Song, error)
	AddSong(song dto.AddSongRequest) (int, error)
	UpdateSong(song dto.UpdateSongRequest) error
	DeleteSong(id int) error
	GetSongLyrics(req dto.GetSongLyricsRequest) ([]string, error)
	GetSongById(id int) (model.Song, error)
	SearchSongs(query string) ([]model.Song, error)
}

type Repository struct {
	Songs
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Songs: postgres.NewSongsPostgres(db),
	}
}
