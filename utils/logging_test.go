package utils

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	assert := assert.New(t)

	n := "test"
	logger := NewLogger(n)
	assert.Equal(logger.Name, n)

	logger.Info("Testing info")

	var buf bytes.Buffer
	logger.SetWriter(&buf)

	logger.Info("Testing info")
	assert.Contains(buf.String(), "Testing info")

	logger.Warn("Testing warn")
	assert.Contains(buf.String(), "Testing warn")

	logger.Debug("Testing debug")
	assert.Contains(buf.String(), "Testing debug")

	logger.Error("Testing error")
	assert.Contains(buf.String(), "Testing error")

	logger.Infof("%s", "Testing info")
	assert.Contains(buf.String(), "Testing info")

	logger.Warnf("%s", "Testing warn")
	assert.Contains(buf.String(), "Testing warn")

	logger.Debugf("%s", "Testing debug")
	assert.Contains(buf.String(), "Testing debug")

	logger.Errorf("%s", "Testing error")
	assert.Contains(buf.String(), "Testing error")
}
