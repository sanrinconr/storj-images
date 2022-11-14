package upload_test

import (
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func defaultTestTimer() func() time.Time {
	return func() time.Time {
		return time.Date(2020, time.February, 12, 0, 0, 0, 0, time.UTC)
	}
}

func goldenFile(t *testing.T, fileName string) []byte {
	const location = "testdata/%s"

	f, err := os.Open(fmt.Sprintf(location, fileName))
	assert.Nil(t, err)
	body, err := io.ReadAll(f)
	assert.Nil(t, err)

	return body
}
