# ArenaMap

`ArenaMap` is a memory-efficient hash map implementation in Go that utilizes an arena allocator to minimize garbage collection overhead. This implementation is particularly useful for high-performance applications where memory allocation and deallocation need to be tightly controlled.

## Features

- **Memory Efficiency**: Utilizes an arena allocator (`arena.Arena`) for memory management, reducing the overhead of individual memory allocations and deallocations.
- **Dynamic Resizing**: Automatically resizes the hash map when the number of stored elements exceeds a certain threshold, ensuring efficient memory usage.
- **Key-Value Storage**: Stores key-value pairs in buckets using a hash table implementation for quick retrieval and insertion.

## Usage

### Initialization

To create a new `ArenaHashMap`, use the `NewArenaHashMap` function, passing an `arena.Arena` instance as a parameter.

```go
mem := arena.NewArena()
hm := NewArenaHashMap(mem)
```

### Insertion

To insert a key-value pair into the hash map, use the `Put` method.

```go
hm.Put("key", "value")
```

### Retrieval

To retrieve a value associated with a key from the hash map, use the `Get` method.

```go
value, found := hm.Get("key")
if found {
    fmt.Println("Value:", value)
} else {
    fmt.Println("Key not found")
}
```

### Deletion

To delete a key-value pair from the hash map, use the `Delete` method.

```go
hm.Delete("key")
```

### Memory Management

To release the memory allocated by the hash map, use the `Free` method.

```go
hm.Free()
```


## Example
```go
package main

import (
    "fmt"

    "github.com/ongniud/arena"
    "github.com/ongniud/arena-map"
)

func main() {
    mem := arena.NewArena()
    defer mem.Close()

    hm := amap.NewArenaHashMap[int, string](mem)

    hm.Put(1, "one")
    hm.Put(2, "two")
    hm.Put(3, "three")

    if val, ok := hm.Get(2); ok {
        fmt.Println("Value for key 2:", val)
    }

    hm.Delete(2)
    if _, ok := hm.Get(2);!ok {
        fmt.Println("Key 2 has been deleted.")
    }
    hm.Free()
}
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

---

Feel free to contribute, report issues, or suggest improvements by creating a pull request or issue on GitHub.
