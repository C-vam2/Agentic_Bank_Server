package api

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Agentic_Bank_Server/token"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		//get authorization header
		authorizationHeader := c.GetHeader(authorizationHeaderKey)

		//check if the autohrization header is valid or not
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is empty")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponsne(err))
			return
		}

		fields := strings.Fields(authorizationHeader)

		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponsne(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])

		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported header type %s", authorizationType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, errResponsne(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])

		if err != nil {
			c.JSON(http.StatusUnauthorized, errResponsne(err))
			return
		}

		c.Set(authorizationPayloadKey, payload)
		c.Next()

	}
}
