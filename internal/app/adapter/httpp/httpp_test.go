package httpp

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ghostorange/cmd/client/tui/clconfig"
	"ghostorange/internal/app/model"
	"ghostorange/internal/app/server"
	"ghostorange/internal/app/srvconfig"
	"ghostorange/internal/app/storage"
	mockstorage "ghostorange/internal/app/storage/mock"
)

func TestProvider(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	strg := mockstorage.NewMockStorage(ctrl)

	srv := testSrv(strg)
	go srv.Run()
	time.Sleep(time.Second)
	defer srv.Shutdown(context.Background())

	cfg := clconfig.New(clconfig.WithAddress("localhost:8080"))

	prov, err := New(cfg, nil)
	require.NoError(t, err)

	strg.EXPECT().
		AddUser(gomock.Any(), gomock.Any(), gomock.Any()).
		Return("user_id", nil)

		strg.EXPECT().
		UserExists(gomock.Any(), gomock.Any()).
		Return(false, nil)

	err = prov.Register(model.Credentials{Login: "test", Password: "test"})
	require.NoError(t, err)

	t.Run("Get Credentials", func(t *testing.T) {
		tt := []model.ItemCredentials{
			{ID: "id",
				User: "user_id",
				Credentials: model.Credentials{
					Login:    "login",
					Password: "password",
				},
				Name:    "case 1",
				Comment: "lucky green",
			},
		}

		strg.EXPECT().
			GetData(gomock.Any(), model.KeyCredentials).
			Return(tt, nil)

		res, err := prov.GetData(model.KeyCredentials)
		require.NoError(t, err)

		v, ok := res.([]model.ItemCredentials)
		assert.True(t, ok)

		assert.Equal(t, tt, v)
	})

	t.Run("Add Credentials", func(t *testing.T) {
		tt := model.ItemCredentials{
			ID: "id",
				User: "user_id",
				Credentials: model.Credentials{
					Login:    "login",
					Password: "password",
				},
				Name:    "case 1",
				Comment: "lucky green",
			}

		strg.EXPECT().
			AddData(gomock.Any(), model.KeyCredentials, gomock.Any(), tt).
			Return(nil)

		err := prov.AddData(model.KeyCredentials, tt)
		require.NoError(t, err)
	})

	t.Run("Count", func(t *testing.T) {
		
		tt := 100

		strg.EXPECT().
			Count(gomock.Any(), model.KeyCredentials, gomock.Any()).
			Return(tt, nil)

		res, err := prov.Count(model.KeyCredentials)
		require.NoError(t, err)

		assert.Equal(t, strconv.Itoa(tt), res)
	})

}

func testSrv(strg storage.Storage) *server.Server {
	vars := map[string]string{
		"SERVER_ADDRESS":   "localhost:8080",
		"SESSION_LIFETIME": "100000000",
	}

	cfg := srvconfig.New(srvconfig.WithEnvVars(vars))

	return server.New(cfg, strg)
}
