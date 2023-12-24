package cache

import (
	"dnevnik-rg.ru/internal/models"
	"sync"
)

type Cache struct {
	Pupils  Pupils
	Coaches Coaches
	Admins  Admins
}

type Pupils struct {
	Pupils map[int]models.Pupil
	mx     *sync.Mutex
}

type Coaches struct {
	Pupils map[int]models.Coach
	mx     *sync.Mutex
}

type Admins struct {
	Pupils map[int]models.Admin
	mx     *sync.Mutex
}

func NewCache() *Cache {
	return &Cache{}
}
