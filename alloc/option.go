package alloc

const (
	B  = 1
	KB = 1024
	MB = 1024 * 1024
)

type Malloc func(size int) []byte

type Option func(o *Options)

func WithSlabSize(size int) Option {
	return func(o *Options) {
		o.SlabSize = size
	}
}

func WithGrowthFactor(factor float64) Option {
	return func(o *Options) {
		o.GrowthFactor = factor
	}
}

func WithMalloc(malloc Malloc) Option {
	return func(o *Options) {
		o.Malloc = malloc
	}
}

func WithAlign(align bool) Option {
	return func(o *Options) {
		o.Align = align
	}
}

func newDefaultOptions() *Options {
	return &Options{
		GrowthFactor: 2,
		SlabSize:     128 * KB,
		Malloc: func(size int) []byte {
			return make([]byte, size)
		},
		Align: false,
	}
}

type Options struct {
	// Size of a single slab in bytes.
	// Defaults to 1MB.
	SlabSize int
	// Step defines Slab classes growing proportion.
	// Must be > 1.0.
	// Defaults to 1.25.
	GrowthFactor float64
	Malloc       Malloc

	Align bool
}
