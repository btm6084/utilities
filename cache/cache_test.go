package cache

import (
	"testing"

	"github.com/btm6084/utilities/metrics"
	"github.com/stretchr/testify/assert"
)

func TestGetSet(t *testing.T) {
	m := &metrics.NoOp{}

	t.Run("Bytes In, Bytes Out Binary", func(t *testing.T) {
		data := []byte{'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07', '\x08', '\x09', '\x10'}
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual []byte
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, data, actual)
	})

	t.Run("Bytes In, String Out Binary", func(t *testing.T) {
		data := []byte{'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07', '\x08', '\x09', '\x10'}
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual string
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, "\x00\x01\x02\x03\x04\x05\x06\a\b\t\x10", actual)
	})

	t.Run("String In, Bytes Out Binary", func(t *testing.T) {
		data := string([]byte{'\x00', '\x01', '\x02', '\x03', '\x04', '\x05', '\x06', '\x07', '\x08', '\x09', '\x10'})

		// We expect this because the input gets encoded. On output, it stays encoded. Use a string to extract it if you want the same thing back.
		expected := []byte{0x5c, 0x75, 0x30, 0x30, 0x30, 0x30, 0x5c, 0x75, 0x30, 0x30, 0x30, 0x31, 0x5c, 0x75, 0x30, 0x30, 0x30, 0x32, 0x5c, 0x75, 0x30, 0x30, 0x30, 0x33, 0x5c, 0x75, 0x30, 0x30, 0x30, 0x34, 0x5c, 0x75, 0x30, 0x30, 0x30, 0x35, 0x5c, 0x75, 0x30, 0x30, 0x30, 0x36, 0x5c, 0x75, 0x30, 0x30, 0x30, 0x37, 0x5c, 0x75, 0x30, 0x30, 0x30, 0x38, 0x5c, 0x74, 0x5c, 0x75, 0x30, 0x30, 0x31, 0x30}
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual []byte
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Bytes In, String Out String", func(t *testing.T) {
		data := []byte{'H', 'E', 'L', 'L', 'O', '!'}
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual string
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, "HELLO!", actual)
	})

	t.Run("Bytes In, Int Out", func(t *testing.T) {
		data := []byte{'1', '2', '5', '4', '2', '7'}
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual int
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, 125427, actual)
	})

	t.Run("Int In, Bytes Out", func(t *testing.T) {
		data := 125427
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual []byte
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, []byte{'1', '2', '5', '4', '2', '7'}, actual)
	})

	t.Run("String In, Int Out", func(t *testing.T) {
		data := "125427"
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual int
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, 125427, actual)
	})

	t.Run("Int In, String Out", func(t *testing.T) {
		data := 125427
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual string
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, "125427", actual)
	})

	t.Run("String In, String Out", func(t *testing.T) {
		data := "Hello There. General Kenobi!"
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual string
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, "Hello There. General Kenobi!", actual)
	})

	t.Run("Int In, Int Out", func(t *testing.T) {
		data := 125427
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual int
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, 125427, actual)
	})

	t.Run("Struct In, Struct Out", func(t *testing.T) {
		type Substruct struct {
			Things map[string]string `json:"things"`
		}
		type TestStruct struct {
			Bool   bool        `json:"bool"`
			Int    int         `json:"int"`
			Slice  []byte      `json:"slice"`
			String string      `json:"string"`
			Stuff  []Substruct `json:"stuff"`
		}

		data := TestStruct{
			Bool:   true,
			Int:    12786,
			String: `Such a nice string`,
			Stuff: []Substruct{
				{
					Things: map[string]string{
						"Oh":      "Didn't See you there",
						"So Nice": "to see you!",
					},
				},
				{
					Things: map[string]string{
						"This is": "a second slice entry.",
					},
				},
			},
		}
		err := Set(m, t.Name(), data)
		assert.Nil(t, err)

		var actual TestStruct
		err = Get(m, t.Name(), &actual)
		assert.Nil(t, err)
		assert.Equal(t, data, actual)
	})

}
