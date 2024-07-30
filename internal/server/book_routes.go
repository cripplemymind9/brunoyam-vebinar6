package server

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) GetBooksHandler(ctx *gin.Context) {
	claims, err := ValidateToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{Message: err.Error()})
		return
	}

	uid, err := s.store.GetUserId(*claims)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{Message: "Unauthorized"})
		return
	}

	books, err := s.store.GetBooksByUserId(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, books)
}

func (s *Server) InsertBooksHandler(ctx *gin.Context) {
	claims, err := ValidateToken(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{Message: err.Error()})
		return
	}

	uid, err := s.store.GetUserId(*claims)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{Message: "Unauthorized"})
		return
	}
	
	var books []models.Book
	if err := ctx.ShouldBindBodyWithJSON(&books); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	for _, book := range books {
		if err := s.validate.Struct(book); err != nil {
			ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
			return
		}
	}
	
	if err := s.store.InsertBooks(books, uid); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Message: "Books were added"})
}