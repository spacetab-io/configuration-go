package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WithConfigPath(t *testing.T) {
	type tc struct {
		name string
		in   string
		exp  string
		err  error
	}

	tcs := []tc{
		{
			name: "existed config path",
			in:   "./config_examples/configuration",
			exp:  "./config_examples/configuration",
			err:  nil,
		},
		{
			name: "default config path",
			in:   "",
			exp:  "./configuration",
			err:  ErrEmptyConfigPath,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.ErrorIs(t, WithConfigPath(tc.in)(&merger{}), tc.err)
		})
	}
}
