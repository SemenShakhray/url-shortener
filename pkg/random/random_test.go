package random_test

import (
	"testing"

	"github.com/SemenShakhray/url-shortener/pkg/random"
	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	cases := []struct {
		name string
		size int
	}{
		{
			name: "size=1",
			size: 1,
		},
		{
			name: "size=4",
			size: 4,
		},
		{
			name: "size=10",
			size: 10,
		},
		{
			name: "size=20",
			size: 20,
		},
	}

	for _, item := range cases {
		t.Run(item.name, func(t *testing.T) {
			alies1, err := random.NewRandomString(item.size)
			alies2, err := random.NewRandomString(item.size)
			assert.NoError(t, err)

			assert.Equal(t, item.size, len(alies1))
			assert.Equal(t, item.size, len(alies2))

			assert.NotEqual(t, alies1, alies2)
		})
	}
}
