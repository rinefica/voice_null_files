package crypto

import (
	"testing"

	"github.com/rinefica/voice_null_files/internal/lib/sl"
	"github.com/stretchr/testify/require"
)

func TestCryptoService(t *testing.T) {
	log := sl.SetupLogger("local")
	sut := NewCryptoService(log, []byte("somecryptopasswrod"))
	require.NotNil(t, sut)

	data := "data message"
	crypto, err := sut.Encrypt([]byte(data))

	println(string(crypto))
	cipherText, err := sut.Decrypt(crypto)
	require.NoError(t, err)
	require.Equal(t, data, string(cipherText))
}
