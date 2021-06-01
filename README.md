# go-cache

Basic golang in-memory cache based on map with no garbage collection.

Since there is no garbage collection for expired items, it is only suitable for small fixed key space.

Example usage:
```go
package main

import (
	"fmt"
	"time"

	"github.com/need-being/go-cache"
)

func main() {
	c := cache.New()

	c.Set("foo", "bar", 100*time.Millisecond)
	if item, ok := c.Get("foo"); ok {
		fmt.Println("found", item)
	}

	time.Sleep(100 * time.Millisecond)
	if _, ok := c.Get("foo"); !ok {
		fmt.Println("foo missed")
	}
}
```
