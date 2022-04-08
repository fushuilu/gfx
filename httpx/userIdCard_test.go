package httpx

import (
	"testing"

	"github.com/fushuilu/golibrary/libx"
	"github.com/stretchr/testify/assert"
)

func TestJwtUserIdCard(t *testing.T) {
	card := IdCard{UserId: 99, Uid: "99abc"}
	ji := NewJwtUserIdCard("998877", cmn.NewTestXRedis())
	token, err := ji.CreateJwtToken(60, card)
	assert.Nil(t, err)

	idCard, err := ji.GetIdCard(token)
	assert.Nil(t, err)
	assert.Equal(t, idCard.UserId, card.UserId)
	assert.Equal(t, idCard.Uid, card.Uid)
}
