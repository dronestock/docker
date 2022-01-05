package main

type env struct {
	key   string
	value string
}

func newEnv(key string, value string) *env {
	return &env{
		key:   key,
		value: value,
	}
}
