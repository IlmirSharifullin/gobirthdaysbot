package mapped

type KVStorage map[int64]map[string]any

func New() *KVStorage {
	return &KVStorage{}
}

func (s *KVStorage) Update(ID int64, key string, value any) {
	if _, ok := (*s)[ID]; !ok {
		(*s)[ID] = map[string]any{key: value}
	} else {
		(*s)[ID][key] = value
	}
}

func (s *KVStorage) Get(ID int64, key string) any {
	if _, ok := (*s)[ID]; !ok {
		return nil
	}
	return (*s)[ID][key]
}

func (s *KVStorage) Clear(ID int64) {
	(*s)[ID] = map[string]any{}
}
