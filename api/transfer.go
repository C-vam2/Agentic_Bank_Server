package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	db "github.com/Agentic_Bank_Server/db/sqlc"
	"github.com/Agentic_Bank_Server/token"
	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) validAccount(c *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponsne(err))
			return account, false
		}

		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch %s vs %s", accountID, account.Currency, currency)
		c.JSON(http.StatusBadRequest, errResponsne(err))
		return account, false
	}

	return account, true
}

func (server *Server) createTransfer(c *gin.Context) {
	var req transferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errResponsne(err))
		return
	}

	fromAccount, valid := server.validAccount(c, req.FromAccountID, req.Currency)
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)

	if !valid {
		return
	}

	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		c.JSON(http.StatusUnauthorized, errResponsne(err))
		return
	}

	if _, valid := server.validAccount(c, req.ToAccountID, req.Currency); !valid {
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	result, err := server.store.TransferTx(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponsne(err))
		return
	}

	c.JSON(http.StatusOK, result)
}
