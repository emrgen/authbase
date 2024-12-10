package x

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_SendMail(t *testing.T) {
	from := os.Getenv("EMAIL_FROM")
	err := SendMail(from, "minorblocker@gmail.com", "subject", "body")
	assert.NoError(t, err)
}
