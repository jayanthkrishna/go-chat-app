package main

import "fmt"

type Store struct {
	data  map[int]string
	cache Cacher
}

func NewStore(c Cacher) *Store {

	data := map[int]string{
		1: "Elon musk is the owner of twitter",
		2: "Foo is not bar",
		3: "Must watch show",
	}
	return &Store{
		data:  data,
		cache: c,
	}
}

func (s *Store) getFromCache(key int) (string, bool) {
	val, ok := s.cache.Get(key)

	if !ok {

		fmt.Println("returning key from cache")
		return val, ok
	}
	return val, ok
}

func (s *Store) Get(key int) (string, error) {

	val, ok := s.cache.Get(key)

	if ok {
		if err := s.cache.Remove(key); err != nil {
			fmt.Println(err)
		}

		fmt.Println("returning the value from the cache")
		return val, nil

	}
	val, err := s.data[key]

	if !err {
		return "", fmt.Errorf("key not found : %d", key)
	}

	if err := s.cache.Set(key, val); err != nil {
		return "", err
	}
	fmt.Println("returning key from internal storage")
	return val, nil
}
