package dto

// @Description Request to add a new song
type AddSongRequest struct {
	Song        string `json:"song" binding:"required"`
	Genre       string `json:"genre"`
	ReleaseDate string `json:"releaseDate" binding:"required"`
	Text        string `json:"text" binding:"required"`
	Link        string `json:"link"`
	Group       string `json:"group" binding:"required"`
}

// @Description Request to update a song
type UpdateSongRequest struct {
	Id          int     `json:"id" binding:"required"`
	Song        *string `json:"song"`
	Genre       *string `json:"genre"`
	ReleaseDate *string `json:"releaseDate"`
	Text        *string `json:"text"`
	Link        *string `json:"link"`
	Group       *string `json:"group"`
}

// @Description Request to getting songs
type GetSongsRequest struct {
	Id          int
	Genre       string
	Song        string
	ReleaseDate string
	Text        string
	Link        string
	GroupId     int 
	Group       string
	Limit       int
	Offset      int
}

// @Description Request to get lyrics
type GetSongLyricsRequest struct {
	Id     int
	Limit  int
	Offset int
}

type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}
