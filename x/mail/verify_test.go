package mail

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_VerifyMail(t *testing.T) {
	err := VerifyEmail("minorblocker@gmail.com", "http://localhost:4001/verify/1234")
	assert.NoError(t, err)
}
