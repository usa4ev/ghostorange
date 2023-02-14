package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeDecodeItems(t *testing.T) {
	cases := []ItemCredentials{
		{Credentials: Credentials{
			Login:    "testlogin",
			Password: "testpassword",
		},
			Comment: "confidential",
			Name:    "My login",
		},
		{Credentials: Credentials{
			Login:    "tanya",
			Password: "dragon",
		},
			Comment: "lol",
			Name:    "Tanya's login",
		},
	}

	b, err := EncodeItemsJSON(cases)
	require.NoError(t, err)

	val, err := DecodeItemsJSON(KeyCredentials, b)
	require.NoError(t, err)

	assert.EqualValues(t, cases, val)
}

