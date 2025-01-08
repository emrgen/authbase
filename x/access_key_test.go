package x

import (
	"github.com/google/uuid"
	"testing"
)

func TestGenerateAccessKey(t *testing.T) {
	accessKey := NewAccessKey()
	if accessKey.ID == uuid.Nil {
		t.Error("ID is nil")
	}
	if accessKey.Value == "" {
		t.Error("Value is empty")
	}
}
