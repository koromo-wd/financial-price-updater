package updater

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestNewGoogleSheet(t *testing.T) {
	assert.EqualValues(
		t,
		&GoogleSheet{
			Option:     option.WithCredentialsFile("/test"),
			SheetID:    "test",
			WriteRange: "A:B",
		},
		NewGoogleSheet("/test", "test", "A:B"),
	)
}
