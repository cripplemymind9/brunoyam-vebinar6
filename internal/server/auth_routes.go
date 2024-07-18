package server

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) LoginHandler(ctx *gin.Context) {
	var input models.LoginUser

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(input); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	id, err := s.store.Login(input)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	token, err := s.CreateToken(id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Message: token})
}

func (s *Server) ProfileHandler (ctx *gin.Context) {
	claims, err := ValidateToken(ctx)
	if err != nil  {
		ctx.JSON(http.StatusUnauthorized, models.Response{Message: err.Error()})
		return
	}

	user, err := s.store.Profile(*claims)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{Message: "Unauthorized"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}