package pow

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func Test_Simple(t *testing.T) {
	data := uuid.New()
	ch := make(chan struct{})
	hash, nonce, err := GetHash(ch, []byte(data.String()))
	require.NoError(t, err)
	fmt.Printf("%x, %d\n", hash.Bytes(), nonce)
	_, result, err := CheckHash([]byte(data.String()), nonce)
	require.NoError(t, err)
	require.True(t, result)
}

func Test_Interrupt(t *testing.T) {
	data := uuid.New()
	ch := make(chan struct{})
	go func() {
		ch <- struct{}{}
	}()
	time.Sleep(time.Millisecond * 100)
	_, _, err := GetHash(ch, []byte(data.String()))
	require.Error(t, err)
}
