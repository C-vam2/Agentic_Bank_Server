package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Agentic_Bank_Server/db/sqlc"
	"github.com/Agentic_Bank_Server/token"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(c *gin.Context) {
	var req createAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponsne(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	account, err := server.store.CreateAccount(c, db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	})

	if err != nil {

		pqErr := pq.As(err)
		fmt.Print(pqErr.Code.Name())

		switch pqErr.Code.Name() {
		case "foreign_key_violation", "unique_violation":
			c.JSON(http.StatusForbidden, errResponsne(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return
	}

	c.JSON(http.StatusOK, account)
	return

}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(c *gin.Context) {
	var req getAccountRequest

	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponsne(err))
		return
	}

	account, err := server.store.GetAccount(c, req.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponsne(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return
	}

	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")

		c.JSON(http.StatusUnauthorized, errResponsne(err))
		return
	}

	c.JSON(http.StatusOK, account)
	return

}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccount(c *gin.Context) {
	var req listAccountRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponsne(err))
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := server.store.ListAccounts(c, arg)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponsne(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return
	}

	c.JSON(http.StatusOK, accounts)
	return

}
