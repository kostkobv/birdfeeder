package utils_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"utils"
)

func TestInitEncoder(t *testing.T) {
	encoder := utils.InitEncoder()

	assert.NotEmpty(t, encoder)
}

func TestUdhenc_Encode(t *testing.T) {
	encoder := utils.InitEncoder()

	t.Run("GSM 7-bit encode", func(t *testing.T) {
		t.Run("encode regular symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"c8329bfd768192a736c89c769759a0b09b0ccabfeb3f"}}
			m, _ := encoder.Encode("Hello. I'm fine, and you?")
			assert.Equal(t, e, *m)
		})

		t.Run("encode 2 space char symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"9bd706b8416d521bdec6b7e96dca1b0a"}}
			m, _ := encoder.Encode(`\|{}[]~€^`)
			assert.Equal(t, e, *m)
		})

		t.Run("encode mixed symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Plain, []string{"406499cd7e8740dd74ffed0279c565b90b34a98036a879be2d7eb3e79b14"}}
			m, _ := encoder.Encode("¡Hello! Ñiño. Über. ΓΩ {symbols}")
			assert.Equal(t, e, *m)
		})
	})

	t.Run("UC-2 encode", func(t *testing.T) {
		t.Run("encode UC-2 symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Unicode, []string{"041f044004380432043504420021d83dde00"}}
			m, _ := encoder.Encode("Привет!😀")
			assert.Equal(t, e, *m)
		})

		t.Run("encode only UC-2 symbols", func(t *testing.T) {
			e := utils.Encoded{utils.Unicode, []string{"041f04400438043204350442d83dde000020041a0430043a002004340435043b0430"}}
			m, _ := encoder.Encode("Привет😀 Как дела")
			assert.Equal(t, e, *m)
		})
	})
}
