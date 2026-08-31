package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/Agentic_Bank_Server/db/mock"
	db "github.com/Agentic_Bank_Server/db/sqlc"
	"github.com/Agentic_Bank_Server/utils"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func newTestServer(t *testing.T, store *mockdb.MockStore) *Server {
	config := utils.Config{
		SymmetricKey:   utils.RandomString(32),
		AccessDuration: time.Minute,
	}

	// start test server and make request
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestGetAccountApi(t *testing.T) {
	user := randomUser(t)
	account := randomAccount(user.Username)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	store := mockdb.NewMockStore(ctrl)

	// build stubs
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	// start test server and make request
	server := newTestServer(t, store)

	url := fmt.Sprintf("/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	addAuthorization("user", time.Minute, server.tokenMaker, t, request, authorizationTypeBearer)
	recorder := httptest.NewRecorder()

	server.router.ServeHTTP(recorder, request)

	// check response
	require.Equal(t, http.StatusOK, recorder.Code)
	requireBodyMatchAccount(t, recorder.Body, account)
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       utils.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)

}

func randomUser(t *testing.T) *db.User {
	return &db.User{
		Username:          utils.RandomOwner(),
		HashedPassword:    "hashedpassword",
		FullName:          utils.RandomString(32),
		Email:             utils.RandomEmail(),
		PasswordChangedAt: time.Now(),
		CreatedAt:         time.Now(),
	}
}
