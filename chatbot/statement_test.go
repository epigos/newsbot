package chatbot

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStatement(t *testing.T) {
	assert := assert.New(t)
	n := "Test"

	st := NewStatement(n, user.ID)
	assert.Equal(st.Text, n)

	res := "Hi"
	st.AddTextResponse(res)
	assert.Equal(res, fmt.Sprintf("%s", st.Responses[0]))

	st.SetScore(1)
	assert.Equal(st.Score, float32(1))

	st.SetPayload("test")
	assert.Equal(st.Payload, "test")

	s := st.SerializeResponse()
	assert.Contains(s, "Hi")
}
