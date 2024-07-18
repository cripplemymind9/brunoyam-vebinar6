package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/cripplemymind9/brunoyam-vebinar6/internal/domain/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Storage interface {
	GetAllUsers() ([]models.User, error)
	GetUser(int) (models.User, error)
	InsertUser(models.User) error
	UpdateUser(int, models.User) error
	DeleteUser(int) error

	//Auth methods
	Login(models.LoginUser) (int, error)
	Profile(models.Claims) (models.User, error)
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

	router.POST("/users", s.InsertUserHandler)

	//Auth routes
	router.POST("/login", s.LoginHandler)
	router.GET("/profile", s.ProfileHandler)

	//Protected routes
	protectedUsers := router.Group("/users")
	protectedUsers.Use(AuthMiddleWare())
	{
		protectedUsers.GET("/", s.GetAllUsersHandler)
		protectedUsers.GET("/:uid", s.GetUserHandler)
		protectedUsers.PUT("/:uid", s.UpdateUserHandler)
		protectedUsers.DELETE("/:uid", s.DeleteTaskHandler)
	}

	return router.Run(s.addr)
}

func (s *Server) CreateToken(id int) (string, error) {
	user, err := s.store.GetUser(id)
	if err != nil  {
		return "", err
	}

	claims := &models.Claims{
		UID: user.UID,
		Login: user.Login,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte("secret_key"))
}

func ValidateToken(ctx *gin.Context) (*models.Claims, error) {
	tokenString := ctx.GetHeader("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	
	return claims, nil
}

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := ValidateToken(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, models.Response{Message: "Unauthorized"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}