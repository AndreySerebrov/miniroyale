package payload

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestXxx(t *testing.T) {
	p, err := New("../../assets/Words-of-Wisdom.txt")
	require.NoError(t, err)
	quote, err := p.GetRandomQuote()
	require.NoError(t, err)
	require.True(t, len(quote) > 0)
	fmt.Println(quote)
}
