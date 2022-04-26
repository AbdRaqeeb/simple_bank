package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	mockdb "github.com/AbdRaqeeb/simple_bank/db/mock"
	db "github.com/AbdRaqeeb/simple_bank/db/sqlc"
	"github.com/AbdRaqeeb/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTransferAPI(t *testing.T) {
	accountOne := randomAccount()
	accountTwo := randomAccount()
	accountThree := randomAccount()

	currencyOne := util.CAD
	currencyTwo := util.NAR

	accountOne.Currency = currencyOne
	accountTwo.Currency = currencyOne
	accountThree.Currency = currencyTwo

	amount := int64(10)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"from_account_id": accountOne.ID,
				"to_account_id":   accountTwo.ID,
				"amount":          amount,
				"currency":        currencyOne,
			},
			buildStubs: func(store *mockdb.MockStore) {
				// expect get account to be called twice
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountOne.ID)).Times(1).Return(accountOne, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountTwo.ID)).Times(1).Return(accountTwo, nil)
				args := db.TransferTxParams{
					FromAccountID: accountOne.ID,
					ToAccountID:   accountTwo.ID,
					Amount:        amount,
				}
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(args)).Times(1)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusOK)
			},
		},
		{
			name: "Invalid Currency",
			body: gin.H{
				"from_account_id": accountOne.ID,
				"to_account_id":   accountTwo.ID,
				"amount":          amount,
				"currency":        "fake",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "Invalid Amount",
			body: gin.H{
				"from_account_id": accountOne.ID,
				"to_account_id":   accountTwo.ID,
				"amount":          -amount,
				"currency":        currencyOne,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "Not Found - FromAccount",
			body: gin.H{
				"from_account_id": accountOne.ID,
				"to_account_id":   accountTwo.ID,
				"amount":          amount,
				"currency":        currencyOne,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountOne.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountTwo.ID)).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusNotFound)
			},
		},
		{
			name: "Not Found - ToAccount",
			body: gin.H{
				"from_account_id": accountOne.ID,
				"to_account_id":   accountTwo.ID,
				"amount":          amount,
				"currency":        currencyOne,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountOne.ID)).Times(1).Return(accountOne, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountTwo.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusNotFound)
			},
		},
		{
			name: "Invalid Currency Mismatch - ToAccount",
			body: gin.H{
				"from_account_id": accountOne.ID,
				"to_account_id":   accountThree.ID,
				"amount":          amount,
				"currency":        currencyOne,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountOne.ID)).Times(1).Return(accountOne, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountThree.ID)).Times(1).Return(accountThree, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "Invalid Currency Mismatch - FromAccount",
			body: gin.H{
				"from_account_id": accountThree.ID,
				"to_account_id":   accountTwo.ID,
				"amount":          amount,
				"currency":        currencyOne,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountThree.ID)).Times(1).Return(accountThree, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountTwo.ID)).Times(0)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
		{
			name: "TransTXError",
			body: gin.H{
				"from_account_id": accountOne.ID,
				"to_account_id":   accountTwo.ID,
				"amount":          amount,
				"currency":        currencyOne,
			},
			buildStubs: func(store *mockdb.MockStore) {
				args := db.TransferTxParams{
					FromAccountID: accountOne.ID,
					ToAccountID:   accountTwo.ID,
					Amount:        amount,
				}

				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountOne.ID)).Times(1).Return(accountOne, nil)
				store.EXPECT().GetAccount(gomock.Any(), gomock.Eq(accountTwo.ID)).Times(1).Return(accountTwo, nil)
				store.EXPECT().TransferTx(gomock.Any(), gomock.Eq(args)).Times(1).Return(db.TransferTxResult{}, sql.ErrTxDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Code, http.StatusInternalServerError)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)

			server := NewServer(store)

			tc.buildStubs(store)

			url := "/transfers"

			recorder := httptest.NewRecorder()

			// Marshall body to JSON
			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}
