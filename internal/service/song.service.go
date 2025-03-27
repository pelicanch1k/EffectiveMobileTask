package service

import (
	"errors"

	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/dto"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/model"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/repository"
)

type SongsService struct {
	repo        *repository.Repository
	externalAPI ExternalAPI
}

func NewSongsService(repo *repository.Repository, externalAPI ExternalAPI) *SongsService {
	return &SongsService{
		repo:        repo,
		externalAPI: externalAPI,
	}
}

func (s SongsService) GetSongs(req dto.GetSongsRequest) ([]model.Song, error) {
	return s.repo.GetSongs(req)
}

func (s SongsService) GetSongById(id int) (model.Song, error) {
	return s.repo.GetSongById(id)
}

func (s SongsService) AddSong(req dto.AddSongRequest) (int, error) {
	details, err := s.externalAPI.GetSongDetails(req.Group, req.Song)
	if err != nil {
		return 0, err
	}

	req.ReleaseDate = details.ReleaseDate
	req.Text = details.Text
	req.Link = details.Link

	return s.repo.AddSong(req)
}

func (s SongsService) UpdateSong(req dto.UpdateSongRequest) error {
	_, err := s.GetSongById(req.Id)
	if err != nil {
		return err
	}

	return s.repo.UpdateSong(req)
}

func (s SongsService) DeleteSong(id int) error {
	return s.repo.DeleteSong(id)
}

func (s SongsService) GetSongLyrics(req dto.GetSongLyricsRequest) ([]string, error) {
	return s.repo.GetSongLyrics(req)
}

func (s SongsService) SearchSongs(query string) ([]model.Song, error) {
	if query == "" {
		return nil, errors.New("поисковый запрос не может быть пустым")
	}

	return s.repo.SearchSongs(query)
}
