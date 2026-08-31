package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Agentic_Bank_Server/token"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func addAuthorization(username string, duration time.Duration, tokenMaker token.Maker, t *testing.T,
	request *http.Request, authorizationType string) {

	fmt.Print("Shivam inside authorization")
	token, err := tokenMaker.CreateToken(username, duration)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationType, token)
	fmt.Printf("authorization header %s %s", authorizationHeader, authorizationHeaderKey)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestMiddleWare(t *testing.T) {
	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization("user", time.Minute, tokenMaker, t, request, authorizationTypeBearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {

				require.Equal(t, http.StatusOK, recorder.Code)
				fmt.Print(recorder.Body)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {

			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "InvalidAuthorizationFormat",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization("user", time.Minute, tokenMaker, t, request, "")
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization("user", -time.Minute, tokenMaker, t, request, authorizationTypeBearer)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusUnauthorized)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t, nil)
			authPath := "/auth"

			server.router.GET(
				authPath,
				authMiddleware(server.tokenMaker),
				func(c *gin.Context) {
					c.JSON(http.StatusOK, gin.H{})

				})

			request, err := http.NewRequest(http.MethodGet, authPath, nil)
			require.NoError(t, err)

			recorder := httptest.NewRecorder()

			tc.setupAuth(t, request, server.tokenMaker)

			server.router.ServeHTTP(recorder, request)
			fmt.Println(recorder)
			tc.checkResponse(t, recorder)

		})
	}
}
