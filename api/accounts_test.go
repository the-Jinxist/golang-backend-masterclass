package api

import (
	mockdb "backend_masterclass/db/mock"
	db "backend_masterclass/db/sqlc"
	"backend_masterclass/token"
	"backend_masterclass/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {

	//Creating a random account
	user, _ := createRandomUser(t)
	account := randomAccount(user.Username)

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

	server := NewTestServer(t, store)
	recorder := httptest.NewRecorder()

	url := fmt.Sprintf("/account/%d", account.ID)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)

	addAuthorization(t, request, server.tokenMaker, authorizationTypeBearer, user.Username, time.Minute)

	server.router.ServeHTTP(recorder, request)

	//Here we are checking response
	require.Equal(t, http.StatusCreated, recorder.Code)

	//We want to match the body recieved from the the test http call with the random one we generated at
	//the top, we have to write a function to convert the buffer in the `recorder` to an Accounts struct
	requireBodyMatchAccount(t, recorder.Body, account)
}

//This method uses an array of test cases so we can cover every possible scenario with the api.
func TestGetAccountAPIWithMultipleTestCases(t *testing.T) {
	user, _ := createRandomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusCreated, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},

		//Here, we start adding more cases for different scenarios

		//This scenario is for when there is no data in the database that matches this account
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Accounts{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)

			},
		},

		//In this unauthorized scenario
		{
			name:      "UnAuthorized",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)

			},
		},

		//This scenario is for when there is an internal error with our logic
		{
			name:      "InternalServerError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Accounts{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)

			},
		},

		//This scenario is for when there's a bad request with our logic, i.e we send a negative id
		{
			name:      "BadRequest",
			accountID: -1,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)

			},
		},
	}

	for i := range testCases {
		testCase := testCases[i]

		t.Run(testCase.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			store := mockdb.NewMockStore(controller)
			testCase.buildStubs(store)

			server := NewTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", testCase.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)

			testCase.setupAuth(t, request, server.tokenMaker)

			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			testCase.checkResponse(t, recorder)

		})
	}
}

func randomAccount(owner string) db.Accounts {

	return db.Accounts{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
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
