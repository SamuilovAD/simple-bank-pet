package gapi

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"testing"
	"time"

	mockdb "github.com/SamuilovAD/simple-bank-pet/db/mock"
	db "github.com/SamuilovAD/simple-bank-pet/db/sqlc"
	"github.com/SamuilovAD/simple-bank-pet/pb"
	"github.com/SamuilovAD/simple-bank-pet/token"
	"github.com/SamuilovAD/simple-bank-pet/util"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t)

	newName := util.RandomOwner()
	newEmail := util.RandomEmail()

	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildStubs    func(store *mockdb.MockStore)
		buildContext  func(t *testing.T, tokenMaker token.Maker) context.Context
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				FullName: &newName,
				Email:    &newEmail,
			},
			buildContext: func(t *testing.T, tokenMaker token.Maker) context.Context {
				return newContextWithBearerToken(
					t,
					tokenMaker,
					user.Username,
					time.Minute,
				)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.UpdateUserParams{
					Username: user.Username,
					FullName: pgtype.Text{
						String: newName,
						Valid:  true,
					},
					Email: pgtype.Text{
						String: newEmail,
						Valid:  true,
					},
				}
				updatedUser := db.User{
					Username:          user.Username,
					HashedPassword:    user.HashedPassword,
					FullName:          newName,
					Email:             newEmail,
					PasswordChangedAt: user.PasswordChangedAt,
					CreatedAt:         user.CreatedAt,
				}
				store.EXPECT().
					UpdateUser(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(updatedUser, nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				updatedUser := res.GetUser()
				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newName, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)

			tc.buildStubs(store)
			server := newTestServer(t, store, nil)

			res, err := server.UpdateUser(tc.buildContext(t, server.tokenMaker), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
