package models

type Song struct {
	ID          int64  `json:"id"  bson:"id"`
	Group       string `json:"group"  bson:"group"`
	Song        string `json:"song"  bson:"song"`
	Text        string `json:"text"  bson:"text"`
	Link        string `json:"link"  bson:"link"`
	ReleaseDate string `json:"release_date"  bson:"release_date"`
}

type SongInfo struct {
	Group       string `json:"group"  bson:"group"`
	Song        string `json:"song"  bson:"song"`
	Text        string `json:"text"  bson:"text"`
	Link        string `json:"link"  bson:"link"`
	ReleaseDate string `json:"release_date"  bson:"release_date"`
}

type CreateSongReq struct {
	Group string `json:"group"  bson:"group"`
	Song  string `json:"song"  bson:"song"`
}
