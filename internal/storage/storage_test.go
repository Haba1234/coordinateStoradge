package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Parallel()
	t.Run("add point", func(t *testing.T) {
		t.Parallel()
		s := NewStorage()
		z := 312345
		s.AddPoint(uint64(z))
		s.RLock()
		val, ok := s.ReadPoint(uint64(z))
		s.RUnlock()
		require.Equal(t, true, val)
		require.Equal(t, true, ok)
		require.Equal(t, 1, s.Len())
	})

	t.Run("read point", func(t *testing.T) {
		t.Parallel()
		s := NewStorage()
		z := 312345
		s.RLock()
		val, ok := s.ReadPoint(uint64(z))
		s.RUnlock()
		require.Equal(t, false, val)
		require.Equal(t, false, ok)
	})
}
