package service

import (
	"errors"
	"song-library/internal/models"
	"song-library/internal/repository"
	"strconv"
	"strings"
)

type SongService interface {
	GetAllSongs(group string, title string, page int, limit int) ([]models.Song, error)
	GetSongByID(id uint) (models.Song, error)
	GetSongLyrics(songID uint, page string, limit string) ([]string, error)
	CreateSong(song models.Song) error
	UpdateSong(song models.Song) error
	DeleteSong(id uint) error
}

type songService struct {
	repo repository.SongRepository
}

func NewSongService(repo repository.SongRepository) SongService {
	return &songService{repo: repo}
}

func (s *songService) GetAllSongs(group string, title string, page int, limit int) ([]models.Song, error) {
	return s.repo.GetAllSongs(group, title, page, limit)
}

func (s *songService) GetSongByID(id uint) (models.Song, error) {
	return s.repo.GetSongByID(id)
}

func (s *songService) GetSongLyrics(songID uint, page string, limit string) ([]string, error) {
	song, err := s.GetSongByID(songID)
	if err != nil {
		return nil, err
	}

	lyrics := strings.Split(song.Text, "\n")
	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		return nil, errors.New("invalid page number")
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		return nil, errors.New("invalid limit number")
	}

	start := (pageNum - 1) * limitNum
	end := start + limitNum

	if start >= len(lyrics) {
		return []string{}, nil
	}
	if end > len(lyrics) {
		end = len(lyrics)
	}

	return lyrics[start:end], nil
}

func (s *songService) CreateSong(song models.Song) error {
	return s.repo.CreateSong(song)
}

func (s *songService) UpdateSong(song models.Song) error {
	return s.repo.UpdateSong(song)
}

func (s *songService) DeleteSong(id uint) error {
	return s.repo.DeleteSong(id)
}
