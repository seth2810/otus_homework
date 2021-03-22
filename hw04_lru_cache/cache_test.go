package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)
		s := []string{"a", "b", "c", "d"}

		for i, n := range s {
			c.Set(Key(n), i)
		}

		val, ok := c.Get(Key(s[0]))
		require.Nil(t, val)
		require.False(t, ok)
	})

	t.Run("purge logic complex", func(t *testing.T) {
		c := NewCache(3)
		s := []string{"a", "b", "c", "d"}

		c.Set(Key(s[0]), 0)
		c.Set(Key(s[1]), 1)
		c.Set(Key(s[2]), 2)

		c.Set(Key(s[1]), 1)
		c.Get(Key(s[2]))
		c.Get(Key(s[0]))
		c.Set(Key(s[3]), 3)

		val, ok := c.Get(Key(s[0]))
		require.True(t, ok)
		require.NotNil(t, val)

		val, ok = c.Get(Key(s[1]))
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get(Key(s[2]))
		require.True(t, ok)
		require.NotNil(t, val)
	})

	t.Run("purge logic random", func(t *testing.T) {
		c := NewCache(3)
		s := []string{"a", "b", "c", "d", "e", "f"}

		for i := 0; i < 3; i++ {
			c.Set(Key(s[i]), i)
		}

		for i := 5; i > 3; i-- {
			c.Get(Key(s[i-3]))
			c.Set(Key(s[i]), i)
		}

		val, ok := c.Get(Key(s[5]))
		require.True(t, ok)
		require.NotNil(t, val)

		val, ok = c.Get(Key(s[4]))
		require.True(t, ok)
		require.NotNil(t, val)

		val, ok = c.Get(Key(s[1]))
		require.True(t, ok)
		require.NotNil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
