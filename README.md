This repo implements an arena-based hashmap designed to significantly reduce the impact of garbage collection (GC).

In some scenarios, maps are required within a process context to maintain items' attributes and features of different process stages for later joining, filtering, or conversion.
In these cases, the map may store a large amount of data but does not require high concurrency or frequent updates.
Once the process is complete, the map can be released.

Since Go is a garbage-collected language, frequent memory allocations and releases, maintain huge count memory blocks and pointers, can negatively impact service performance.
An arena-based map is well-suited for this situation, as it minimizes the overhead of memory management.


## Guide

### What is memory arena?
A memory arena is a method of memory management where a large block of memory is allocated at once and portions of it are used to satisfy allocation requests from the program.

In the context of a garbage-collected language such as Go, the use of memory arenas can offer several advantages:

- Performance Improvement: By allocating memory in large blocks, memory arenas reduce the overhead associated with frequent calls to the system's memory allocator. This can lead to performance improvements, especially in applications that perform many small allocations.
- Enhanced Cache Locality: Memory arenas can also improve cache locality by allocating closely related objects within the same block of memory. This arrangement increases the likelihood that when one object is accessed, other related objects are already in the cache, thus reducing cache misses and enhancing overall application performance.

However, while memory arenas offer these advantages, they are not a silver bullet and come with trade-offs, such as potentially increased memory usage due to unused space within the allocated blocks. Careful consideration and profiling are necessary to determine whether using a memory arena is beneficial for a particular application.

### What is slab allocation?

The memory is managed via a simple slab allocator algorithm

Each arena tracks one or more slabClass structs. Each slabClass manages a different "chunk size", where chunk sizes are computed using a simple "growth factor" (e.g., the "power of 2 growth" in the above example). Each slabClass also tracks zero or more slabs, where every slab tracked by a slabClass will all have the same chunk size. A slab manages a (usually large) continguous array of memory bytes (1MB from the above example), and the slab's memory is subdivided into many chunks of the same chunk size. All the chunks in a new slab are placed on a free-list that's part of the slabClass.

When Alloc() is invoked, the first "large enough" slabClass is found, and a chunk from the free-list is taken to service the allocation. If there are no more free chunks available in a slabClass, then a new slab (e.g., 1MB) is allocated, chunk'ified, and the request is processed as before.

See: http://en.wikipedia.org/wiki/Slab_allocation


### Suitable Scenarios:

- GC-Intensive Workloads: Designed for applications where frequent memory allocations cause significant GC overhead, helping to reduce performance impacts.
- Frequent Memory Allocations: Ideal for operations with high allocation frequency, as ArenaMap reuses pre-allocated memory, reducing allocation overhead.

### Unsuitable Scenarios:
- High Concurrency: Not suitable for concurrent access since it is not thread-safe. For concurrent operations, use Go’s built-in Map or other thread-safe structures.
- Dynamic Object Lifecycle Management: If your application requires frequent object creation and destruction with fine-grained lifecycle control, this library may not be ideal.
- Low Memory Usage: Not necessary for applications with small memory usage where GC overhead is minimal.
- Severe Memory Fragmentation: Not suitable for scenarios with a lot of small objects or highly variable object sizes, as it may lead to memory fragmentation.


## Limitations:
- Supports Only Base Types: Only basic types (including string) are supported currently. Complex types will be supported soon.
- No Memory Alignment Support: Memory alignment is not currently supported due to implementation complexity, but solutions are being explored and are under verification.
- Higher Memory Cost: ArenaMap consumes more memory than Go's standard map.
- Not Concurrency-Safe: ArenaMap is not thread-safe. Please use your own locking.

## Example

## Requirement
- go1.18+
- go module project

## TODO
- Memory alignment support.
- Complex object allocation.
- More effective hashmap realism.
- More documents.
- More testing coverage.
- More clear error reminders.

## Contributing
Contributions from the community are welcome! This project accepts contributions via GitHub pull requests:

- Fork it
- Create your feature branch (git checkout -b my-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-feature)
- Create new Pull Request

If you experience any issues, please let us know via GitHub issues.
I appreciate detailed and accurate reports that help us to identity and replicate the issue.

## Contact
If you have any questions, feedback or suggestions, please feel free to contact me. 
I'm always open to feedback and would love to hear from you!

### License
This project is licensed under the terms of the MIT license.
