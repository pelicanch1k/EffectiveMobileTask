package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/dto"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/model"
	"github.com/pelicanch1k/EffectiveMobileTestTask/pkg/logging"
)

type SongsPostgres struct {
	db     *sqlx.DB
	logger *logging.Logger
}

func NewSongsPostgres(db *sqlx.DB) *SongsPostgres {
	return &SongsPostgres{db: db, logger: logging.GetLogger()}
}

func (s SongsPostgres) GetSongs(req dto.GetSongsRequest) ([]model.Song, error) {
	var songs []model.Song

	query := `
		SELECT s.id, s.song, s.genre, TO_CHAR(s.releaseDate, 'DD.MM.YYYY') as releaseDate, s.text, s.link, 
			   s.group_id, g.name as group_name
		FROM songs s
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE 1=1
	`

	params := []interface{}{}
	paramCount := 1

	if req.Id != 0 {
		query += fmt.Sprintf(" AND s.id = $%d", paramCount)
		params = append(params, req.Id)
		paramCount++
	}

	if req.Genre != "" {
		query += fmt.Sprintf(" AND s.genre ILIKE $%d", paramCount)
		params = append(params, "%"+req.Genre+"%")
		paramCount++
	}

	if req.Song != "" {
		query += fmt.Sprintf(" AND s.song ILIKE $%d", paramCount)
		params = append(params, "%"+req.Song+"%")
		paramCount++
	}

	if req.ReleaseDate != "" {
		query += fmt.Sprintf(" AND TO_CHAR(s.releaseDate, 'DD.MM.YYYY') = $%d", paramCount)
		params = append(params, req.ReleaseDate)
		paramCount++
	}

	if req.Text != "" {
		query += fmt.Sprintf(" AND s.text ILIKE $%d", paramCount)
		params = append(params, "%"+req.Text+"%")
		paramCount++
	}

	if req.Link != "" {
		query += fmt.Sprintf(" AND s.link ILIKE $%d", paramCount)
		params = append(params, "%"+req.Link+"%")
		paramCount++
	}

	if req.GroupId != 0 {
		query += fmt.Sprintf(" AND s.group_id = $%d", paramCount)
		params = append(params, req.GroupId)
		paramCount++
	}

	if req.Group != "" {
		query += fmt.Sprintf(" AND g.name ILIKE $%d", paramCount)
		params = append(params, "%"+req.Group+"%")
		paramCount++
	}

	query += " ORDER BY s.id"

	if req.Limit != 0 {
		query += fmt.Sprintf(" LIMIT $%d", paramCount)
		params = append(params, req.Limit)
		paramCount++

		if req.Offset != 0 {
			query += fmt.Sprintf(" OFFSET $%d", paramCount)
			params = append(params, req.Offset)
		}
	}

	if err := s.db.Select(&songs, query, params...); err != nil {
		s.logger.Errorf("Ошибка при получении песен: %v", err)
		return nil, err
	}

	return songs, nil
}

func (s SongsPostgres) AddSong(req dto.AddSongRequest) (int, error) {
	var songId int
	var groupId int

	// Начинаем транзакцию
	tx, err := s.db.Beginx()
	if err != nil {
		s.logger.Errorf("Ошибка при начале транзакции: %v", err)
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	err = tx.QueryRow("SELECT id FROM groups WHERE name = $1", req.Group).Scan(&groupId)
	if err == sql.ErrNoRows {
		// Группа не существует, создаем новую
		err = tx.QueryRow("INSERT INTO groups (name) VALUES ($1) RETURNING id", req.Group).Scan(&groupId)
		if err != nil {
			s.logger.Errorf("Ошибка при создании группы: %v", err)
			return 0, err
		}
	} else if err != nil {
		s.logger.Errorf("Ошибка при поиске группы: %v", err)
		return 0, err
	}

	query := `INSERT INTO songs (song, genre, releaseDate, text, link, group_id) 
              VALUES ($1, $2, TO_DATE($3, 'DD.MM.YYYY'), $4, $5, $6) RETURNING id`

	err = tx.QueryRow(query, req.Song, "Unknown", req.ReleaseDate, req.Text, req.Link, groupId).Scan(&songId)
	if err != nil {
		s.logger.Errorf("Ошибка при добавлении песни: %v", err)
		return 0, err
	}

	// Фиксируем транзакцию
	if err = tx.Commit(); err != nil {
		s.logger.Errorf("Ошибка при фиксации транзакции: %v", err)
		return 0, err
	}

	return songId, nil
}

func (s SongsPostgres) UpdateSong(req dto.UpdateSongRequest) error {
	tx, err := s.db.Beginx()
	if err != nil {
		s.logger.Errorf("Ошибка при начале транзакции: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	updateSongQuery := "UPDATE songs SET "
	updateParams := []interface{}{}
	paramCount := 1

	if req.Song != nil {
		updateSongQuery += fmt.Sprintf("song = $%d, ", paramCount)
		updateParams = append(updateParams, *req.Song)
		paramCount++
	}

	if req.Genre != nil {
		updateSongQuery += fmt.Sprintf("genre = $%d, ", paramCount)
		updateParams = append(updateParams, *req.Genre)
		paramCount++
	}

	if req.ReleaseDate != nil {
		updateSongQuery += fmt.Sprintf("releaseDate = TO_DATE($%d, 'DD.MM.YYYY'), ", paramCount)
		updateParams = append(updateParams, *req.ReleaseDate)
		paramCount++
	}

	if req.Text != nil {
		updateSongQuery += fmt.Sprintf("text = $%d, ", paramCount)
		updateParams = append(updateParams, *req.Text)
		paramCount++
	}

	if req.Link != nil {
		updateSongQuery += fmt.Sprintf("link = $%d, ", paramCount)
		updateParams = append(updateParams, *req.Link)
		paramCount++
	}

	if len(updateParams) > 0 {
		updateSongQuery = updateSongQuery[:len(updateSongQuery)-2]

		updateSongQuery += fmt.Sprintf(" WHERE id = $%d", paramCount)
		updateParams = append(updateParams, req.Id)

		_, err = tx.Exec(updateSongQuery, updateParams...)
		if err != nil {
			s.logger.Errorf("Ошибка при обновлении песни: %v", err)
			return err
		}
	}

	if req.Group != nil {
		var groupId int

		err = tx.QueryRow("SELECT id FROM groups WHERE name = $1", *req.Group).Scan(&groupId)
		if err == sql.ErrNoRows {
			err = tx.QueryRow("INSERT INTO groups (name) VALUES ($1) RETURNING id", *req.Group).Scan(&groupId)
			if err != nil {
				s.logger.Errorf("Ошибка при создании новой группы: %v", err)
				return err
			}
		} else if err != nil {
			s.logger.Errorf("Ошибка при поиске группы: %v", err)
			return err
		}

		_, err = tx.Exec("UPDATE songs SET group_id = $1 WHERE id = $2", groupId, req.Id)
		if err != nil {
			s.logger.Errorf("Ошибка при обновлении привязки к группе: %v", err)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		s.logger.Errorf("Ошибка при фиксации транзакции: %v", err)
		return err
	}

	return nil
}

func (s SongsPostgres) DeleteSong(id int) error {
	query := "DELETE FROM songs WHERE id = $1"

	_, err := s.db.Exec(query, id)
	return err
}

func (s SongsPostgres) GetSongLyrics(req dto.GetSongLyricsRequest) ([]string, error) {
	var lyrics string

	query := "SELECT text FROM songs WHERE id = $1"

	err := s.db.QueryRow(query, req.Id).Scan(&lyrics)
	if err == sql.ErrNoRows {
		return nil, errors.New("song not found")
	} else if err != nil {
		return nil, errors.New("failed to fetch song lyrics")
	}

	verses := strings.Split(lyrics, "\n\n")
	start := (req.Offset - 1) * req.Limit
	end := start + req.Limit
	if start >= len(verses) {
		return verses, nil
	}
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], nil
}

func (s SongsPostgres) GetSongById(id int) (model.Song, error) {
	var song model.Song

	query := `
		SELECT s.id, s.song, s.genre, TO_CHAR(s.releaseDate, 'DD.MM.YYYY') as releaseDate, s.text, s.link, 
			   s.group_id, g.name as group_name
		FROM songs s
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE s.id = $1
	`

	err := s.db.Get(&song, query, id)
	if err != nil {
		s.logger.Errorf("Ошибка при получении песни по ID %d: %v", id, err)
		return model.Song{}, err
	}

	return song, nil
}

func (s SongsPostgres) SearchSongs(query string) ([]model.Song, error) {
	var songs []model.Song

	
	searchParam := "%" + query + "%"

	sqlQuery := `
		SELECT s.id, s.song, s.genre, TO_CHAR(s.releaseDate, 'DD.MM.YYYY') as releaseDate, s.text, s.link, 
			   s.group_id, g.name as group_name
		FROM songs s
		LEFT JOIN groups g ON s.group_id = g.id
		WHERE 
			s.song ILIKE $1 OR
			g.name ILIKE $1 OR
			s.text ILIKE $1 OR
			s.genre ILIKE $1
		ORDER BY s.id
	`

	if err := s.db.Select(&songs, sqlQuery, searchParam); err != nil {
		s.logger.Errorf("Ошибка при поиске песен: %v", err)
		return nil, err
	}

	return songs, nil
}
