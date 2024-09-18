package client

import "time"

type User struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}

type Video struct {
	Id          string     `json:"id"`
	UserId      string     `json:"user_id"`
	Status      string     `json:"status"`
	Path        string     `json:"path"`
	Size        int64      `json:"size"`
	MimeType    string     `json:"mime_type"`
	Metadata    string     `json:"metadata"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Metadata struct {
	Resolution string `json:"resolution"`
	Duration   int    `json:"duration"`
	Format     string `json:"format"`
	Codec      string `json:"codec"`
	Bitrate    int    `json:"bitrate"`
}

type VideoSession struct {
	Id           string    `json:"id"`
	UserId       string    `json:"user_id"`
	VideoId      string    `json:"video_id"`
	FragmentHash string    `json:"fragment_hash"`
	FragmentPath string    `json:"fragment_path"`
	Token        string    `json:"token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
