package service

import (
	"errors"

	"movies_service/model"
	"movies_service/repository"

	"gorm.io/gorm"
)

type MovieService interface {
	CreateMovie(movie *model.Movie) error
	GetMovies() ([]model.Movie, error)
	GetMovie(id uint) (*model.Movie, error)
	UpdateMovie(id uint, data *model.Movie) error
	DeleteMovie(id uint) error
}

type movieServiceImpl struct {
	movieRepo repository.MovieRepository
}

func NewMovieService(movieRepo repository.MovieRepository) MovieService {
	return &movieServiceImpl{movieRepo: movieRepo}
}

func (s *movieServiceImpl) CreateMovie(movie *model.Movie) error {
	return s.movieRepo.Create(movie)
}

func (s *movieServiceImpl) GetMovies() ([]model.Movie, error) {
	return s.movieRepo.GetAll()
}

func (s *movieServiceImpl) GetMovie(id uint) (*model.Movie, error) {
	movie, err := s.movieRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return movie, nil
}

func (s *movieServiceImpl) UpdateMovie(id uint, data *model.Movie) error {
	data.ID = id
	err := s.movieRepo.Update(data)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (s *movieServiceImpl) DeleteMovie(id uint) error {
	if err := s.movieRepo.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return err
	}
	return nil
}
