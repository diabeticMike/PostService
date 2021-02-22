package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfigJSON(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		filepath := "../../config.json"

		config, err := readConfigJSON(filepath)

		assert.NoError(t, err)
		assert.NotEmpty(t, config)
	})

	t.Run("wrong filePath", func(t *testing.T) {
		filepath := ""

		config, err := readConfigJSON(filepath)

		assert.Error(t, err)
		assert.Empty(t, config)
	})
}

func TestNew(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		filepath := "../../config.json"

		config, err := New(filepath)

		assert.NoError(t, err)
		assert.NotEmpty(t, config)
	})

	t.Run("wrong filePath", func(t *testing.T) {
		filepath := ""

		config, err := New(filepath)

		assert.Error(t, err)
		assert.Empty(t, config)
	})
}
