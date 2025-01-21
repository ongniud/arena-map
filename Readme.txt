



1. 类似 tslab 那样，实现一个自定义指针，格式类似 prt:=sc_id|slb_id|ck_id ，通过 ref() 获取 pointer 指针；
这样释放、调试、增/减引用时比较方便，不需要再根据 pointer 和 size 去查（size比较麻烦），带来的问题时，使用者需要理解和保存这个指针；

type Ptr[T any] uint32
func (ptr Ptr[T]) ref() ref {
	return ref{slabID: int(ptr >> 16), slotID: int(ptr & 0xffff)}
}
func (ptr Ptr[T]) IsNil() bool { return ptr == 0 }

或者

type EObject[T any] uintptr
func (e EObject[T]) Value() T {
	return *(*T)(unsafe.Pointer(e))
}
func (e EObject[T]) Set(t T) {
	p := (*T)(unsafe.Pointer(e))
	*p = t
}

2. 实现 share 引用，通过 ref 控制是否释放，这在有些场景中可以减少分配，节约内存；

3. 是否需要清空内存？延伸一点，是否需要给 malloc 支持各种 options ，比如 align、reset ...

4. 是否需要缩容？

5. 因为 unsafe.Pointer 引用的内存还是会被 GC 识别，如果在 Arena 释放时，内部仍有内存被引用，那么整块 arena 可能不会被回收，可能有内存泄漏；
一种解决思路是在整个 arena 被释放时，检查内部是否所有 slab 都已释放完毕，也即 free list 均满的，否则给出 warning ，虽不能解决问题但也没有忽视；
这种问题无法解决，因为还有引用关系在，直接释放会导致程序异常，但是 arena 带来的问题是可能放大了这个问题的影响面积，导致大片内存被占用。

可以参考下 golang arena 的实现，看看其是如何解决类似问题的。


6. 内存地址对齐，参考 go-mem 中 arrow 的实现；

7. 对于指针类型的 T ，还需要返回 *T 吗？这是二级指针，用起来有点复杂。
[重要]如果能较好的支持指针，复杂类型的支持都比较简单，可以先通过 arena 分配具体数据，再保存指针；把识别数据类型的心智负担交给用户；

背景：
很多时候，在一个请求上下文中，使用 map 存储 item 维度的信息，用于后续聚合处理；
这个 map 没有高并发，不会频繁删除，在处理结束后可以释放；

通过这个库你可以学习：
1. slab 的分配原理
2. 基于 arena 实现 map
3. ...

说明：
1. 只支持基础类型（包括string)；不支持复杂类型；
2. 不支持内存对齐；支持思路：
    (1) 每个 slab class 中，分配 chunk 时遍历找到地址符合要求的，这要求对 free chunk list 进行改造，支持遍历删除等；但是可能遍历完也找不到合适的；
    (2) 每个 slab class 中，分配不同 align 的子 chunk free 列表，如果找不到，在新建 slab 时使用 mmap 支持从指定 offset 分配堆内存；
3. 

