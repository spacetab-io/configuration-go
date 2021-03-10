package config

import (
	"bytes"
	logs "log"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_checkConfigPath(t *testing.T) {
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
			err:  syscall.ENOENT,
		},
	}

	t.Parallel()

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out, err := checkConfigPath(tc.in)

			if tc.err != nil && err != nil {
				assert.ErrorIs(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tc.exp, out)
		})
	}
}

func Test_iSay(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logs.SetOutput(&buf)

	defer func() {
		logs.SetOutput(os.Stderr)
	}()

	log("some text")

	assert.Equal(t, time.Now().Format("2006/01/02 15:04:05")+" [config] some text\n", buf.String())
}
