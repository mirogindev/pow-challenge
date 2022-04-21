package hashcash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHashCash(t *testing.T) {

	t.Run("Check validation with valid message", func(t *testing.T) {
		payload := HashCashData{
			Ver:      1,
			Bits:     3,
			Date:     "220421192952",
			Resource: "127.0.0.1:58520",
			Rand:     "MTg4ODMy",
			Counter:  "OTkyMzM1",
		}
		sha1, valid := IsChallengeValid(&payload)
		assert.True(t, valid)
		assert.Equal(t, sha1, "00054d961d431e978500d4f75728632d16c3d6d7")
	})

	t.Run("Check validation with 4 zeros message", func(t *testing.T) {
		payload := HashCashData{
			Ver:      1,
			Bits:     4,
			Date:     "220421192952",
			Resource: "127.0.0.1:58520",
			Rand:     "MTg4ODMy",
			Counter:  "OTkyMzM1",
		}
		sha1, valid := IsChallengeValid(&payload)
		assert.False(t, valid)
		assert.Equal(t, sha1, "600886ecc85bd9d01d3d0bca28dafe014be5d91f")
	})

	t.Run("Check zeros with valid hash", func(t *testing.T) {
		hash := "00054d961d431e978500d4f75728632d16c3d6d7"
		ok := checkZeros(hash, Zero, 3)
		assert.True(t, ok)
	})

	t.Run("Check zeros with invalid hash", func(t *testing.T) {
		hash := "0054d961d431e978500d4f75728632d16c3d6d7"
		ok := checkZeros(hash, Zero, 3)
		assert.False(t, ok)
	})

}
