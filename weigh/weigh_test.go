package weigh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNeatSize(t *testing.T) {

	cases := []struct {
		expected string
		size     int64
	}{
		{
			expected: "123 bytes",
			size:     int64(123),
		},
		{
			expected: "12.06 KiB",
			size:     int64(12345),
		},
		{
			expected: "11.77 MiB",
			size:     int64(12345678),
		},
		{
			expected: "11.50 GiB",
			size:     int64(12345678901),
		},
		{
			expected: "11.23 TiB",
			size:     int64(12345678012345),
		},
		{
			expected: "10.97 PiB",
			size:     int64(12345678012345678),
		},
	}

	for _, tc := range cases {
		require.Equal(t, tc.expected, neatSize(tc.size))
	}

}
