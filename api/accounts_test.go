package api

import (
	mockdb "backend_masterclass/db/mock"
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {

	//Creating a random account
	account := randomAccount()

	//The controller watches out if the calls that are supposed to be made are made
	controller := gomock.NewController(t)
	defer controller.Finish()

	//Create a fake store
	store := mockdb.NewMockStore(controller)

	//This stub just makes sure we're calling the right method
	//when the endpoint "/accounts/:id" is called
	store.EXPECT().
		GetAccount(gomock.Any(), gomock.Eq(account.ID)).
		Times(1).
		Return(account, nil)

	server := NewServer(store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	server.router.ServeHTTP(recorder, request)

	//Here we are checking response
	require.Equal(t, http.StatusCreated, recorder.Code)

	//We want to match the body recieved from the the test http call with the random one we generated at
	//the top, we have to write a function to convert the buffer in the `recorder` to an Accounts struct
	requireBodyMatchAccount(t, recorder.Body, account)
}

func randomAccount() db.Accounts {
	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Accounts) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gottenAccount db.Accounts
	err = json.Unmarshal(data, &gottenAccount)

	fmt.Println(gottenAccount)

	require.NoError(t, err)
	require.Equal(t, gottenAccount.ID, account.ID)
	require.Equal(t, gottenAccount.Owner, account.Owner)
	require.Equal(t, gottenAccount.Balance, account.Balance)
	require.Equal(t, gottenAccount.Currency, account.Currency)
}
