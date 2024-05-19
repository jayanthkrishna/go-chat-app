package main

type Cacher interface {
	Get(int) (string, bool)
	Remove(int) error
	Set(int, string) error
}

type NopCache struct{}

func (c NopCache) Get(int) (string, bool) {
	return "", false
}

func (c NopCache) Set(int, string) error {
	return nil
}

func (c NopCache) Remove(int) error {
	return nil
}
