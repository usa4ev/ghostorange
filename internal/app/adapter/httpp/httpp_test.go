package httpp

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/usa4ev/ghostorange/internal/app/model"
	"github.com/usa4ev/ghostorange/internal/app/server"
	"github.com/usa4ev/ghostorange/internal/app/srvconfig"
	"github.com/usa4ev/ghostorange/internal/app/storage"
	mockstorage "github.com/usa4ev/ghostorange/internal/app/storage/mock"
	"github.com/usa4ev/ghostorange/internal/app/tui/clconfig"
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

	t.Run("Get Card", func(t *testing.T) {
		tt := model.ItemCard{
			ID:                 "id",
			Number:             "1001",
			Exp:                time.Now().Add(time.Hour * 24000),
			CardholderName:     "mr. Cardholder",
			CardholderSurename: "Smith",
			CVVHash:            "secure",
			Name:               "case 1",
			Comment:            "lucky green",
		}

		strg.EXPECT().
			GetCardInfo(gomock.Any(), tt.ID, gomock.Any()).
			Return(tt, nil)

		item, err := prov.GetCard(tt.ID, tt.CVVHash)
		require.NoError(t, err)

		assert.WithinDuration(t, tt.Exp, item.Exp, 0)
		tt.Exp = time.Time{}
		item.Exp = time.Time{}
		assert.Equal(t, tt, item)
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
