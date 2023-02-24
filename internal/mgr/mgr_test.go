package mgr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_DoubleInsert(t *testing.T) {
	mgr := New(time.Second)
	err := mgr.Push("test")
	require.NoError(t, err)
	err = mgr.Push("test")
	require.Error(t, err)
}

func Test_Check(t *testing.T) {
	mgr := New(time.Second)
	err := mgr.Push("test")
	require.NoError(t, err)
	err = mgr.Check("test")
	require.NoError(t, err)
}

func Test_Use(t *testing.T) {
	mgr := New(time.Second)
	err := mgr.Push("test")
	require.NoError(t, err)
	err = mgr.Use("test")
	require.NoError(t, err)
	err = mgr.Use("test")
	require.Error(t, err)
}

func Test_Expire(t *testing.T) {
	mgr := New(time.Second)
	err := mgr.Push("test")
	require.NoError(t, err)
	time.Sleep(time.Second)
	err = mgr.Use("test")
	require.Error(t, err)
}
