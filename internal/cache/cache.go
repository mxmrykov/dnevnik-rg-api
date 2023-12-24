package cache

import "dnevnik-rg.ru/internal/models"

type Cache struct {
	Pupils  map[int]models.Pupil
	Coaches map[int]models.Coach
}

func NewCache() *Cache {
	return &Cache{}
}
