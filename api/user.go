package api

import (
	"fmt"
	"net/http"

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

	c.JSON(http.StatusOK, user)
	return

}
