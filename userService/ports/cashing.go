package ports

type Caching interface {
	Set(key string, value interface{})
	Get(key string) interface{}
	Evict(key string)
	// Other caching methods...
}
