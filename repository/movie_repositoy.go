package repository

import (
	"movies_service/model"

	"gorm.io/gorm"
)

type MovieRepository interface {
	Create(movie *model.Movie) error
	GetAll() ([]model.Movie, error)
	GetByID(id uint) (*model.Movie, error)
	Update(movie *model.Movie) error
	Delete(id uint) error
}

type movieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) MovieRepository {
	return &movieRepository{db: db}
}

func (r *movieRepository) Create(movie *model.Movie) error {
	return r.db.Create(movie).Error
}

func (r *movieRepository) GetAll() ([]model.Movie, error) {
	var movies []model.Movie
	err := r.db.Find(&movies).Error
	return movies, err
}

func (r *movieRepository) GetByID(id uint) (*model.Movie, error) {
	var movie model.Movie
	err := r.db.First(&movie, id).Error
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *movieRepository) Update(movie *model.Movie) error {
	var existing model.Movie
	if err := r.db.First(&existing, movie.ID).Error; err != nil {
		return err
	}
	return r.db.Save(movie).Error
}

func (r *movieRepository) Delete(id uint) error {
	res := r.db.Delete(&model.Movie{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
