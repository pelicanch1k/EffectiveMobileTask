package external_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/dto"
)

type SongAPI struct {
	baseURL string
}

func NewSongAPI() *SongAPI {
	baseURL := os.Getenv("URL_ADD_SONG")
	if baseURL == "" {
		baseURL = "http://localhost:8080/api"
	}

	return &SongAPI{
		baseURL: baseURL,
	}
}

// GetSongDetails получает информацию о песне из внешнего API
func (api *SongAPI) GetSongDetails(group, song string) (*dto.SongDetails, error) {
	if group == "" || song == "" {
		return nil, errors.New("отсутствуют обязательные поля group и song")
	}

	requestURL := fmt.Sprintf("%s/info?group=%s&song=%s",
		api.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(song))

	resp, err := http.Get(requestURL)
	if err != nil {
		return nil, errors.New("ошибка при выполнении запроса к внешнему API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("внешний API вернул код ошибки: %d", resp.StatusCode)
	}

	var songDetails dto.SongDetails
	if err := json.NewDecoder(resp.Body).Decode(&songDetails); err != nil {
		return nil, errors.New("ошибка декодирования ответа от внешнего API")
	}

	return &songDetails, nil
}
