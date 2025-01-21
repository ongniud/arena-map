This repo implements an arena-based hashmap designed to significantly reduce the impact of garbage collection (GC).

In certain scenarios, maps are used to store item attributes and process stage features for later operations such as joining, filtering, or conversion. These maps may handle large amounts of data but do not require high concurrency or frequent updates. Once the process is completed, the map can be released.

In Go, frequent memory allocations and deallocations can negatively impact performance due to the overhead of managing numerous memory blocks and pointers. An arena-based map is ideal for these situations, as it minimizes memory management overhead.

## Guide

### What is memory arena?
A memory arena is a memory management technique where a large block of memory is allocated at once, and parts of it are used for allocations.

In garbage-collected languages like Go, memory arenas offer:
- Memory Reuse: By reusing memory within the arena, it reduces the need for frequent allocations, which can lead to performance improvements.
- Reduced GC Pause: Fewer allocations mean less work for the garbage collector, leading to shorter GC pauses.
- Improved Memory Locality: Allocating related objects close together enhances cache performance, speeding up memory operations.

However, memory arenas also have trade-offs, such as potentially higher memory usage due to unused space. Careful consideration and profiling are needed to determine if they are the right solution for your application.

### What is slab allocation?

Slab allocation is a simple memory management method where memory is divided into fixed-size chunks called slabs.

Each arena manages slabs, which are large blocks of memory divided into chunks of the same size. When a new allocation is needed, the first available chunk is used. If no chunks are available, a new slab is allocated.

For more details, see http://en.wikipedia.org/wiki/Slab_allocation


### Suitable Scenarios:

- GC-Intensive Workloads:  Perfect for applications where frequent memory allocations result in significant GC overhead, helping to reduce performance impacts.
- Frequent Memory Allocations: Ideal for operations with high allocation frequency, as ArenaMap reuses pre-allocated memory, reducing allocation overhead.


## Limitations:
- Supports Only Base Types: Only basic types (including string) are supported currently. Complex types will be supported soon.
- No Memory Alignment Support: Memory alignment is not currently supported due to implementation complexity, however, solutions are being explored and tested.
- Higher Memory Cost: ArenaMap consumes more memory than Go's standard map.
- Not Concurrency-Safe: ArenaMap is not thread-safe. Please use your own locking.

## Example

## Requirement
- go 1.18+
- go module project

## TODO
- Support for memory alignment.
- Support for complex object allocation.
- Improve hashmap realism and efficiency.
- Improve documentation.
- Increase test coverage.
- Improve error handling.

## Contributing
Contributions from the community are welcome! This project accepts contributions via GitHub pull requests:

- Fork the repository.
- Create a feature branch (git checkout -b my-feature)
- Commit your changes (git commit -am 'Add some feature')
- Push to the branch (git push origin my-feature)
- Create a pull request.

If you encounter any issues, please let us know by opening a GitHub issue. We appreciate detailed and accurate reports that help us identify and reproduce the issue.

## Contact
If you have any questions, feedback or suggestions, please feel free to contact me. 
I'm always open to feedback and would love to hear from you!

### License
This project is licensed under the MIT License.
