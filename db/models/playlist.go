package models

import (
	"fmt"

	"github.com/anhthii/go-echo/db"
	"github.com/jinzhu/gorm"
)

// Song that belongs to many Playlists and have many artists
type Song struct {
	ID        string      `json:"id" gorm:"PRIMARY_KEY"`
	Name      string      `json:"name"`
	Thumbnail string      `json:"thumbnail" sql:"type:text;"`
	Playlists []*Playlist `gorm:"many2many:playlist_songs;"`
	Artists   []*Artist   `json:"artists" gorm:"foreignkey:SongID"`
}

// Artist belong to a song
type Artist struct {
	Name   string `json:"name" gorm:"PRIMARY_KEY"`
	Link   string `json:"link"`
	SongID string
}

// Playlist that has many songs
type Playlist struct {
	gorm.Model
	Songs     []*Song `json:"songs" gorm:"many2many:playlist_songs;"`
	Title     string  `json:"title"`
	UserRefer string
}

// AfterFind replace nil value wiht empty slice for client to handle
func (playlist *Playlist) AfterFind() (err error) {
	if playlist.Songs == nil {
		playlist.Songs = []*Song{}
	}
	return
}

// AfterFind replace nil value wiht empty slice for client to handle
func (song *Song) AfterFind() (err error) {
	if song.Artists == nil {
		song.Artists = []*Artist{}
	}
	return
}

func checkRecordExistence(model interface{}, query string, field ...interface{}) bool {
	if err := db.GetDB().Where(query, field...).First(model).Error; gorm.IsRecordNotFoundError(err) {
		return false
	}
	return true
}

// AddToPlaylist add a song to playlist
func (song *Song) AddToPlaylist(username, playlistTitle string) (ok bool, errStr string) {
	// check if song exists in playlist
	playlist := &Playlist{}
	if exist := checkRecordExistence(playlist, "user_refer = ?", username); !exist {
		return false, fmt.Sprintf("playlist %s does not exist", playlistTitle)
	}

	var songs []Song
	if db.GetDB().Model(&playlist).Related(&songs, "Songs").
		Where("id = ? ", song.ID).First(&song).RecordNotFound() {
		playlist.Songs = append(playlist.Songs, song)
		db.GetDB().Save(playlist)
		return true, ""
	}

	return false, fmt.Sprintf("%s song already exists %s in playlist", song.Name, playlistTitle)
}

// GetPlaylists of a specific user
func GetPlaylists(username string) (playlists []Playlist) {
	playlists = []Playlist{}
	user := &User{}
	if exist := checkRecordExistence(user, "username = ?", username); !exist {
		return
	}
	db.GetDB().Preload("Songs.Artists").Preload("Songs").Model(user).Related(&playlists, "UserRefer")
	return
}

// CreatePlaylist create a playlist for a user with a title, if a playlist with
// the provided title exists then do not create one
func CreatePlaylist(username, playlistTitle string) (data map[string]interface{}, errStr string) {
	user := &User{}
	if exist := checkRecordExistence(user, "username = ?", username); !exist {
		return nil, fmt.Sprintf("user %s does not exist", username)
	}

	// if user doesn't have a playlist with this title, create a new one
	playlist := &Playlist{}
	if exist := checkRecordExistence(playlist, "user_refer = ? AND title = ", username, playlistTitle); !exist {
		newPlaylist := &Playlist{Title: playlistTitle, UserRefer: username}
		db.GetDB().Create(newPlaylist)
		return map[string]interface{}{
			"_username": username,
			"playlists": GetPlaylists(username),
		}, ""
	}

	return nil, fmt.Sprintf("playlist with title %s already exists", playlistTitle)
}

// DeletePlaylist delete a playlist and it's associations to other tables
func DeletePlaylist(username, playlistTitle string) (ok bool, errStr string) {

	user := &User{}
	if exist := checkRecordExistence(user, "username = ?", username); !exist {
		return false, fmt.Sprintf("user %s does not exist", username)
	}
	playlist := &Playlist{}
	if err := db.GetDB().Where("user_refer = ? and title = ?", username, playlistTitle).Delete(playlist).Error; err != nil {
		return false, fmt.Sprintf("Errors occured while deleting playlist")
	}
	return true, ""
}
