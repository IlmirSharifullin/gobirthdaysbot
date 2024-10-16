package kvstorage

type KVStorage interface {
	Update(ID int64, key string, value any)
	Get(ID int64, key string) any
	Clear(ID int64)
}
