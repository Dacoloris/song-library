package repository

import (
	"song-library/internal/models"
	"song-library/pkg/db"
)

type SongRepository interface {
	GetAllSongs(group string, title string, page int, limit int) ([]models.Song, error)
	GetSongByID(id uint) (models.Song, error)
	CreateSong(song models.Song) error
	UpdateSong(song models.Song) error
	DeleteSong(id uint) error
}

type songRepository struct {
	db *db.Db
}

func NewSongRepository(db *db.Db) SongRepository {
	return &songRepository{db: db}
}

func (r *songRepository) GetAllSongs(group string, title string, page int, limit int) ([]models.Song, error) {
	var songs []models.Song
	query := r.db.Model(&models.Song{})

	if group != "" {
		query = query.Where("group = ?", group)
	}
	if title != "" {
		query = query.Where("title = ?", title)
	}

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&songs).Error
	return songs, err
}

func (r *songRepository) GetSongByID(id uint) (models.Song, error) {
	var song models.Song
	err := r.db.First(&song, id).Error
	return song, err
}

func (r *songRepository) CreateSong(song models.Song) error {
	return r.db.Create(&song).Error
}

func (r *songRepository) UpdateSong(song models.Song) error {
	return r.db.Save(&song).Error
}

func (r *songRepository) DeleteSong(id uint) error {
	return r.db.Delete(&models.Song{}, id).Error
}
