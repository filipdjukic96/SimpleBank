package gapi

import (
	mockdb "bank/db/mock"
	db "bank/db/sqlc"
	"bank/pb"
	"bank/util"
	"bank/worker"
	mockwk "bank/worker/mock"
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type eqCreateUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (expected eqCreateUserTxParamsMatcher) Matches(x interface{}) bool {
	actualArgument, ok := x.(db.CreateUserTxParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.password, actualArgument.HashedPassword)
	if err != nil {
		return false
	}

	expected.arg.HashedPassword = actualArgument.HashedPassword
	if !reflect.DeepEqual(expected.arg.CreateUserParams, actualArgument.CreateUserParams) {
		return false
	}

	// call the AfterCreate function from the actual argument
	err = actualArgument.AfterCreate(expected.user)

	return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreateUserTxParamsMatcher{arg, password, user}
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	// returns user, password variables, no need to specify them
	return
}

func TestCreateUserAPI(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		requestBody   *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			requestBody: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						Username: user.Username,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}

				store.EXPECT().CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).Times(1).Return(db.CreateUserTxResult{User: user}, nil)

				taskPayload := &worker.PayloadSendVerifyEmail{
					Username: user.Username,
				}
				taskDistributor.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), taskPayload, gomock.Any()).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)

				createdUser := res.GetUser()
				require.Equal(t, user.Username, createdUser.Username)
				require.Equal(t, user.FullName, createdUser.FullName)
				require.Equal(t, user.Email, createdUser.Email)
			},
		},
		{
			name: "InternalError",
			requestBody: &pb.CreateUserRequest{
				Username: user.Username,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, taskDistributor *mockwk.MockTaskDistributor) {
				store.EXPECT().CreateUserTx(gomock.Any(), gomock.Any()).Times(1).Return(db.CreateUserTxResult{}, sql.ErrConnDone)
				taskDistributor.EXPECT().DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.Error(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), codes.Internal)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			taskCtrl := gomock.NewController(t)
			defer storeCtrl.Finish() // make sure all expected methods are called
			defer taskCtrl.Finish()

			store := mockdb.NewMockStore(storeCtrl)
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)

			// build stubs for mock store
			tc.buildStubs(store, taskDistributor)

			// start test server and send http request
			server := NewTestServer(t, store, taskDistributor)

			// call server gRPC API handler function directly
			res, err := server.CreateUser(context.Background(), tc.requestBody)

			//check response
			tc.checkResponse(t, res, err)
		})
	}
}
