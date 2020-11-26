package cache

type Loader interface {
	Load(key string) ([]byte, error)
}

type LoaderFunc func(key string) ([]byte, error)

func (f LoaderFunc) Load(key string) ([]byte, error) {
	return f(key)
}
