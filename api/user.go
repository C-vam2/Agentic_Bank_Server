package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	db "github.com/Agentic_Bank_Server/db/sqlc"
	"github.com/Agentic_Bank_Server/utils"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	Email             string    `json:"email"`
	FullName          string    `json:"full_name"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (server *Server) createUser(c *gin.Context) {
	var req createUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponsne(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return
	}

	user, err := server.store.CreateUser(c, db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	})

	if err != nil {

		pqErr := pq.As(err)
		fmt.Print(pqErr.Code.Name())

		switch pqErr.Code.Name() {
		case "unique_violation":
			c.JSON(http.StatusForbidden, errResponsne(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return
	}

	c.JSON(http.StatusOK, newUserResponse(user))
	return
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(c *gin.Context) {
	var req loginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponsne(err))
		return
	}

	user, err := server.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponsne(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return
	}

	if err := utils.CheckPassword(req.Password, user.HashedPassword); err != nil {
		c.JSON(http.StatusUnauthorized, errResponsne(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(
		req.Username,
		server.config.AccessDuration,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return
	}

	res := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserResponse(user),
	}

	c.JSON(http.StatusOK, res)
}
