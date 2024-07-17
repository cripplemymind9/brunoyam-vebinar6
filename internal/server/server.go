package server

import (
	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"

	"fmt"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Storage interface {
	GetAllUsers() ([]models.User, error)
	GetUser(int) (models.User, error)
	InsertUser(models.User) error
	UpdateUser(int, models.User) error
	DeleteUser(int) error
}

type Server struct {
	addr 		string
	store 		Storage
	validate 	*validator.Validate
}

func NewServer (addr string, store Storage) *Server {
	return &Server{
		addr: 		addr,
		store: 		store,
		validate: 	validator.New(),
	}
}

func (s *Server) Run() error {
	router := gin.Default()

	usersRoutes := router.Group("/users") 
	{
		usersRoutes.GET("/", s.GetAllUsersHandler)
		usersRoutes.GET("/:uid", s.GetUserHandler)
		usersRoutes.POST("/", s.InsertUserHandler)
		usersRoutes.PUT("/:uid", s.UpdateUserHandler)
		usersRoutes.DELETE("/:uid", s.DeleteTaskHandler)
	}
	
	return router.Run(s.addr)
}

func (s *Server) GetAllUsersHandler(ctx *gin.Context) {
	users, err := s.store.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (s *Server) InsertUserHandler(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(user); err != nil  {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.InsertUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Message: "User was saved!"})
}

func (s *Server) GetUserHandler(ctx *gin.Context) {
	param := ctx.Param("uid")
	uid, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	user, err := s.store.GetUser(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (s *Server) UpdateUserHandler(ctx *gin.Context) {
	param := ctx.Param("uid")
	uid, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	var user models.User
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.validate.Struct(user); err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.UpdateUser(uid, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	} 

	ctx.JSON(http.StatusOK, models.Response{Message: fmt.Sprintf("User №%v was updated!", uid)})
}

func (s *Server) DeleteTaskHandler(ctx *gin.Context) {
	param := ctx.Param("uid")
	uid, err := strconv.Atoi(param)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{Message: err.Error()})
		return
	}

	if err := s.store.DeleteUser(uid); err != nil {
		ctx.JSON(http.StatusInternalServerError, models.Response{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{Message: fmt.Sprintf("User №%v was deleted!", uid)})
}