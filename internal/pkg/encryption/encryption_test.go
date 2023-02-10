package encryption

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"ghostorange/internal/app/model"
)

func TestEncryptDecrypt(t *testing.T) {
	tt := model.Credentials{}

	// Encrypt
	buf := bytes.NewBuffer(nil)
	enc := json.NewEncoder(buf)
	require.NoError(t, enc.Encode(tt))

	encrypted, err := Encrypt(buf.Bytes())
	require.NoError(t, err)

	// Decrypt
	b, err := Decrypt(encrypted)
	require.NoError(t, err)

	buf = bytes.NewBuffer(b)
	dec := json.NewDecoder(buf)

	var res model.Credentials

	require.NoError(t, dec.Decode(&res))

	assert.Equal(t, tt, res)
}
