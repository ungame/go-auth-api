package security

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

func TestNewToken(t *testing.T) {
	id := bson.NewObjectId()

	token, err := NewToken(id.Hex())
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken(t *testing.T) {
	id := bson.NewObjectId()

	token, err := NewToken(id.Hex())
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	payload, err := ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, id.Hex(), payload.Id)
	assert.Equal(t, id.Hex(), payload.Issuer)

	now := time.Now()
	expiresAt := time.Unix(payload.ExpiresAt, 0)
	assert.True(t, now.Before(expiresAt))

	issuedAt := time.Unix(payload.IssuedAt, 0)
	assert.Equal(t, now.Year(), issuedAt.Year())
	assert.Equal(t, now.Month(), issuedAt.Month())
	assert.Equal(t, now.Day(), issuedAt.Day())
	assert.Equal(t, now.Hour(), issuedAt.Hour())
	assert.Equal(t, now.Minute(), issuedAt.Minute())
}
