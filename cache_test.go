package cache

import (
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestCacheSet(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name:  "Test int",
			key:   "int",
			value: 42,
		},
		{
			name:  "Test string",
			key:   "string",
			value: "foo",
		},
		{
			name:  "Test pointer",
			key:   "pointer",
			value: &struct{}{},
		},
		{
			name:  "Test struct",
			key:   "struct",
			value: struct{}{},
		},
		{
			name:  "Test nil",
			key:   "nil",
			value: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Set(tt.key, tt.value, time.Hour)
			got, found := c.Get(tt.key)
			if !found {
				t.Errorf("cache.Get() found = %v, want %v", found, true)
			}
			if !reflect.DeepEqual(got, tt.value) {
				t.Errorf("cache.Get() got = %v, want %v", got, tt.value)
			}
		})
	}
}

func TestCacheConcurrentSet(t *testing.T) {
	c := New()
	var wg sync.WaitGroup
	const n = 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(strconv.Itoa(i), i, time.Hour)
		}(i)
	}
	wg.Wait()

	for i := 0; i < n; i++ {
		got, found := c.Get(strconv.Itoa(i))
		if !found {
			t.Errorf("cache.Get() found = %v, want %v", found, true)
		}
		if !reflect.DeepEqual(got, i) {
			t.Errorf("cache.Get() got = %v, want %v", got, i)
		}
	}
}

func TestCacheConcurrentGet(t *testing.T) {
	c := New()
	const n = 1000
	for i := 0; i < n; i++ {
		c.Set(strconv.Itoa(i), i, time.Hour)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			got, found := c.Get(strconv.Itoa(i))
			if !found {
				t.Errorf("cache.Get() found = %v, want %v", found, true)
			}
			if !reflect.DeepEqual(got, i) {
				t.Errorf("cache.Get() got = %v, want %v", got, i)
			}
		}(i)
	}
	wg.Wait()
}

func TestCacheConcurrentDelete(t *testing.T) {
	c := New()
	const n = 1000
	for i := 0; i < n; i++ {
		c.Set(strconv.Itoa(i), i, time.Hour)
	}

	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := strconv.Itoa(i)
			c.Delete(key)
			_, found := c.Get(key)
			if found {
				t.Errorf("cache.Get() found = %v, want %v", found, false)
			}
		}(i)
	}
	wg.Wait()
}

func TestCacheConcurrency(t *testing.T) {
	c := New()
	var wg sync.WaitGroup
	const n = 1000
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			key := strconv.Itoa(i)
			c.Set(key, i, time.Hour)
			got, found := c.Get(key)
			if !found {
				t.Errorf("cache.Get() found = %v, want %v", found, true)
			}
			if !reflect.DeepEqual(got, i) {
				t.Errorf("cache.Get() got = %v, want %v", got, i)
			}

			c.Delete(key)
			_, found = c.Get(key)
			if found {
				t.Errorf("cache.Get() found = %v, want %v", found, false)
			}
		}(i)
	}
	wg.Wait()
}

func TestCacheTTL(t *testing.T) {
	c := New()

	c.Set("foo", "foo", 100*time.Millisecond)
	c.Set("bar", "bar", 200*time.Millisecond)

	time.Sleep(100 * time.Millisecond)

	_, found := c.Get("foo")
	if found {
		t.Errorf("cache.Get(\"foo\") found = %v, want %v", found, false)
	}

	got, found := c.Get("bar")
	if !found {
		t.Errorf("cache.Get(\"bar\") found = %v, want %v", found, true)
	}
	if !reflect.DeepEqual(got, "bar") {
		t.Errorf("cache.Get(\"bar\") got = %v, want %v", got, "bar")
	}

	time.Sleep(100 * time.Millisecond)

	_, found = c.Get("bar")
	if found {
		t.Errorf("cache.Get(\"bar\") found = %v, want %v", found, false)
	}
}
