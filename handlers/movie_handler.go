package handlers

import (
	"net/http"
	"strconv"

	"movies_service/model"
	"movies_service/service"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	movieService service.MovieService
}

func NewMovieHandler(movieService service.MovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

// CreateMovie godoc
// @Summary Create a movie
// @Description Add a new movie to the collection
// @Tags Movies
// @Accept json
// @Produce json
// @Param movie body model.Movie true "Movie data"
// @Success 201 {object} model.Movie
// @Failure 400 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /movies [post]
// @Security BearerAuth
func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movie model.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie data"})
		return
	}
	if err := h.movieService.CreateMovie(&movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create movie"})
		return
	}
	c.JSON(http.StatusCreated, movie)
}

// GetMovies godoc
// @Summary List movies
// @Description Get a list of all movies
// @Tags Movies
// @Accept json
// @Produce json
// @Success 200 {array} model.Movie
// @Failure 401 {object} model.ErrorResponse
// @Router /movies [get]
// @Security BearerAuth
func (h *MovieHandler) GetMovies(c *gin.Context) {
	movies, err := h.movieService.GetMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch movies"})
		return
	}
	c.JSON(http.StatusOK, movies)
}

// GetMovie godoc
// @Summary Get movie
// @Description Get details of a movie by ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} model.Movie
// @Failure 404 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /movies/{id} [get]
// @Security BearerAuth
func (h *MovieHandler) GetMovie(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}
	movie, err := h.movieService.GetMovie(uint(id))
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch movie"})
		}
		return
	}
	c.JSON(http.StatusOK, movie)
}

// UpdateMovie godoc
// @Summary Update movie
// @Description Update an existing movie by ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Param movie body model.Movie true "Movie data"
// @Success 200 {object} model.Movie
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /movies/{id} [put]
// @Security BearerAuth
func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}
	var movieUpdates model.Movie
	if err := c.ShouldBindJSON(&movieUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie data"})
		return
	}
	err = h.movieService.UpdateMovie(uint(id), &movieUpdates)
	if err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update movie"})
		}
		return
	}
	movieUpdates.ID = uint(id)
	c.JSON(http.StatusOK, movieUpdates)
}

// DeleteMovie godoc
// @Summary Delete movie
// @Description Delete a movie by ID
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} model.ErrorResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /movies/{id} [delete]
// @Security BearerAuth
func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}
	if err := h.movieService.DeleteMovie(uint(id)); err != nil {
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete movie"})
		}
		return
	}
	c.Status(http.StatusNoContent)
}
