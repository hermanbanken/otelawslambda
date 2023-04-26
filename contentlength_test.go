package otelawslambda

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContentLength(t *testing.T) {
	assert.Equal(t, 3000, contentLength(strings.Repeat("foo", 1000), false))

	testWithNr := func(l int) func(t *testing.T) {
		return func(t *testing.T) {
			var withPadding, withoutPadding string
			withoutPadding = base64.RawStdEncoding.EncodeToString([]byte(strings.Repeat("A", l)))
			withPadding = base64.StdEncoding.EncodeToString([]byte(strings.Repeat("A", l)))
			t.Log(withPadding, len(withPadding), len(strings.TrimPrefix(withPadding, "A")))

			d, e := base64.RawStdEncoding.DecodeString(withoutPadding)
			assert.NoError(t, e)
			assert.Equal(t, len(d), contentLength(withoutPadding, true), "contentLength was %d but expected length of %q was %d", contentLength(withoutPadding, true), withoutPadding, len(d))

			d2, e := base64.StdEncoding.DecodeString(withPadding)
			assert.NoError(t, e)
			assert.Equal(t, len(d2), contentLength(withPadding, true), "contentLength was %d but expected length of %q was %d", contentLength(withPadding, true), withPadding, len(d2))
		}
	}

	for l := 2; l <= 16; l++ {
		t.Run(fmt.Sprintf("%d", l), testWithNr(l))
	}
}

func FuzzContentLength(f *testing.F) {
	for l := 2; l <= 16; l++ {
		f.Add(strings.Repeat("A", l))
		f.Add(base64.RawStdEncoding.EncodeToString([]byte(strings.Repeat("A", l))))
		f.Add(base64.StdEncoding.EncodeToString([]byte(strings.Repeat("A", l))))
	}
	f.Add("")
	f.Add("\r\r")
	f.Add("0\r0")
	f.Fuzz(func(t *testing.T, a string) {
		result, e := base64.RawStdEncoding.DecodeString(a)
		if e != nil {
			var e2 error
			result, e2 = base64.StdEncoding.DecodeString(a)
			if e2 != nil {
				t.Skip("not valid base64")
			}
		}
		if base64.StdEncoding.EncodeToString(result) != a && base64.RawStdEncoding.EncodeToString(result) != a {
			t.Skip("not identical before & after encoding")
		}

		if contentLength(a, true) != len(result) {
			t.Errorf("contentLength was %d but expected length of %q was %d", contentLength(a, true), a, len(result))
		}
	})

}
